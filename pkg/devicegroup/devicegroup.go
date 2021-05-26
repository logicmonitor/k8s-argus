package devicegroup

import (
	"fmt"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/devicecache"
	cache2 "github.com/logicmonitor/k8s-argus/pkg/devicecache/cache"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/logicmonitor/lm-sdk-go/client"
	"github.com/logicmonitor/lm-sdk-go/client/lm"
	"github.com/logicmonitor/lm-sdk-go/models"
	"github.com/sirupsen/logrus"
)

const (
	hasCategoryOpen = "hasCategory("
	existsOpen      = "exists("
	closingBracket  = ")"
)

// Options are the options for creating a device group.
type Options struct {
	AppliesTo                         AppliesToBuilder
	AppliesToDeletedGroup             AppliesToBuilder
	AppliesToConflict                 AppliesToBuilder
	Client                            *client.LMSdkGo
	Name                              string
	ParentID                          int32
	DisableAlerting                   bool
	DeleteDevices                     bool
	FullDisplayNameIncludeClusterName bool
	CustomProperties                  PropertyBuilder
}

// AppliesToBuilder is an interface for building an appliesTo string.
type AppliesToBuilder interface {
	HasCategory(string) AppliesToBuilder
	Auto(string) AppliesToBuilder
	And() AppliesToBuilder
	Custom(string) AppliesToBuilder
	Or() AppliesToBuilder
	Equals(string) AppliesToBuilder
	Exists(string) AppliesToBuilder
	String() string
}

// AppliesToBuilderImpl impl
type AppliesToBuilderImpl struct {
	value string
}

// NewAppliesToBuilder is the builder for appliesTo.
func NewAppliesToBuilder() AppliesToBuilder {
	return &AppliesToBuilderImpl{value: ""} // nolint: exhaustivestruct
}

// And and
func (a *AppliesToBuilderImpl) And() AppliesToBuilder {
	a.value += " && "

	return a
}

// Or or
func (a *AppliesToBuilderImpl) Or() AppliesToBuilder {
	a.value += " || "

	return a
}

// Equals equals
func (a *AppliesToBuilderImpl) Equals(val string) AppliesToBuilder {
	a.value += " == " + fmt.Sprintf(`"%s"`, val)

	return a
}

// HasCategory has
func (a *AppliesToBuilderImpl) HasCategory(category string) AppliesToBuilder {
	a.value += hasCategoryOpen + fmt.Sprintf(`"%s"`, category) + closingBracket

	return a
}

// Exists exists
func (a *AppliesToBuilderImpl) Exists(property string) AppliesToBuilder {
	a.value += existsOpen + fmt.Sprintf(`"%s"`, property) + closingBracket

	return a
}

// Auto auto
func (a *AppliesToBuilderImpl) Auto(property string) AppliesToBuilder {
	a.value += "auto." + property

	return a
}

// Custom custom
func (a *AppliesToBuilderImpl) Custom(property string) AppliesToBuilder {
	a.value += property

	return a
}

// String string
func (a *AppliesToBuilderImpl) String() string {
	return a.value
}

// PropertyBuilder is an interface for building properties
type PropertyBuilder interface {
	Add(string, string, bool) PropertyBuilder
	AddProperties([]config.PropOpts) PropertyBuilder
	Build([]*models.NameAndValue) []*models.NameAndValue
}

type propertyBuilder struct {
	properties []config.PropOpts
}

// NewPropertyBuilder is the builder for properties
func NewPropertyBuilder() PropertyBuilder {
	return &propertyBuilder{} // nolint: exhaustivestruct
}

func (p *propertyBuilder) Add(key string, value string, override bool) PropertyBuilder {
	opts := config.PropOpts{
		Name:     key,
		Value:    value,
		Override: &override,
	}
	p.properties = append(p.properties, opts)

	return p
}

func (p *propertyBuilder) AddProperties(properties []config.PropOpts) PropertyBuilder {
	p.properties = append(p.properties, properties...)

	return p
}

func (p *propertyBuilder) Build(existingProps []*models.NameAndValue) []*models.NameAndValue {
	props := make([]*models.NameAndValue, 0)
	exProps := make(map[string]string)

	for _, elm := range existingProps {
		exProps[*elm.Name] = *elm.Value
	}

	for _, prop := range p.properties {
		key := prop.Name
		value := prop.Value
		override := true
		if prop.Override != nil {
			override = *prop.Override
		}

		if val, ok := exProps[key]; !override && ok {
			value = val
		}

		props = append(props, &models.NameAndValue{Name: &key, Value: &value})
	}

	return props
}

// Create creates a device group.
func Create(lctx *lmctx.LMContext, opts *Options, cache *devicecache.ResourceCache) (int32, error) {
	clusterDeviceGroup, err := Find(opts.ParentID, opts.Name, opts.Client)
	if err != nil {
		return 0, err
	}

	if clusterDeviceGroup == nil {
		cdg, err := create(opts.Name, opts.AppliesTo.String(), opts.DisableAlerting, opts.ParentID, opts.CustomProperties.Build(nil), opts.Client)
		if err != nil {
			return 0, err
		}

		clusterDeviceGroup = cdg
	} else {
		updateCustomProperties(lctx, clusterDeviceGroup, opts)
	}

	if !opts.DeleteDevices && opts.AppliesToDeletedGroup != nil {
		deletedDeviceGroup, err := Find(clusterDeviceGroup.ID, constants.DeletedDeviceGroup, opts.Client)
		if err != nil {
			return 0, err
		}
		if deletedDeviceGroup == nil {
			_, err := create(constants.DeletedDeviceGroup, opts.AppliesToDeletedGroup.String(), true, clusterDeviceGroup.ID, NewPropertyBuilder().Build(nil), opts.Client)
			if err != nil {
				return 0, err
			}
		}
	}

	err = createConflictDynamicGroup(opts, clusterDeviceGroup.ID)
	if err != nil {
		return 0, err
	}

	cache.Set(cache2.ResourceName{
		Name:     *clusterDeviceGroup.Name,
		Resource: enums.Namespaces,
	}, cache2.ResourceMeta{ // nolint: exhaustivestruct
		Container:   "Nodes",
		LMID:        clusterDeviceGroup.ParentID,
		DisplayName: *clusterDeviceGroup.Name,
		Labels:      nil,
	})

	return clusterDeviceGroup.ID, nil
}

func createConflictDynamicGroup(opts *Options, clusterGrpID int32) error {
	if opts.AppliesToConflict != nil && !opts.FullDisplayNameIncludeClusterName {
		conflictingDeviceGroup, err := Find(clusterGrpID, constants.ConflictDeviceGroup, opts.Client)
		if err != nil {
			return err
		}
		if conflictingDeviceGroup == nil {
			_, err := create(constants.ConflictDeviceGroup, opts.AppliesToConflict.String(), true, clusterGrpID, nil, opts.Client)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Find searches for a device group identified by the parent ID and name.
func Find(parentID int32, name string, client *client.LMSdkGo) (*models.DeviceGroup, error) {
	params := lm.NewGetDeviceGroupListParams()
	fields := "name,id,parentId,subGroups,customProperties"
	params.SetFields(&fields)
	filter := fmt.Sprintf("parentId:\"%v\",name:\"%s\"", parentID, name)
	params.SetFilter(&filter)
	restResponse, err := client.LM.GetDeviceGroupList(params)
	if err != nil {
		return nil, fmt.Errorf("failed to get device group list when searching for %q: %w", name, err)
	}

	var deviceGroup *models.DeviceGroup

	if len(restResponse.Payload.Items) == 0 {
		logrus.Errorf("Could not find device group %q with parentId %v", name, parentID)

		return deviceGroup, nil
	}
	deviceGroup = restResponse.Payload.Items[0]

	return deviceGroup, nil
}

// FindDeviceGroupsByName searches for a device group by name.
func FindDeviceGroupsByName(name string, client *client.LMSdkGo) ([]*models.DeviceGroup, error) {
	params := lm.NewGetDeviceGroupListParams()
	fields := "name,id,parentId,subGroups"
	params.SetFields(&fields)
	filter := fmt.Sprintf("name:\"%s\"", name)
	params.SetFilter(&filter)
	restResponse, err := client.LM.GetDeviceGroupList(params)
	if err != nil {
		return nil, err
	}

	var deviceGroups []*models.DeviceGroup
	if restResponse != nil && restResponse.Payload != nil {
		deviceGroups = restResponse.Payload.Items
	}

	return deviceGroups, nil
}

// Exists returns true if the specified device group exists in the account
func Exists(parentID int32, name string, client *client.LMSdkGo) bool {
	deviceGroup, err := Find(parentID, name, client)
	if err != nil {
		logrus.Warnf("Failed looking up device group for node role %q: %v", name, err)
	}

	if deviceGroup != nil {
		return true
	}

	return false
}

// ExistsByID returns true if we could get the group by id
func ExistsByID(groupID int32, client *client.LMSdkGo) bool {
	params := lm.NewGetDeviceGroupByIDParams()
	params.SetID(groupID)
	fields := "name,id"
	params.SetFields(&fields)
	restResponse, err := client.LM.GetDeviceGroupByID(params)
	if err != nil {
		logrus.Warnf("Failed to get device group (id=%v): %v", groupID, err)

		return false
	}

	logrus.Debugf("Exists device group by id response:%#v", restResponse)

	if restResponse.Payload != nil && restResponse.Payload.ID == groupID {
		return true
	}

	return false
}

// GetByID returns true if we could get the group by id
func GetByID(groupID int32, client *client.LMSdkGo) (*models.DeviceGroup, error) {
	params := lm.NewGetDeviceGroupByIDParams()
	params.SetID(groupID)
	restResponse, err := client.LM.GetDeviceGroupByID(params)
	if err != nil {
		return nil, err
	}

	return restResponse.Payload, nil
}

func GetClusterGroupProperty(lctx *lmctx.LMContext, name string, client *client.LMSdkGo) string {
	log := lmlog.Logger(lctx)
	clusterGroupID := util.GetClusterGroupID(lctx, client)
	log.Infof("Cluster group ID: %v", clusterGroupID)
	deviceGroup, err := GetByID(clusterGroupID, client)
	if err != nil {
		log.Errorf("Failed to fetch cluster group: %s", err)
	}
	return GetPropertyValue(deviceGroup, name)
}

// DeleteSubGroup deletes a subgroup from a device group with the specified
// name.
func DeleteSubGroup(deviceGroup *models.DeviceGroup, name string, client *client.LMSdkGo) error {
	for _, subGroup := range deviceGroup.SubGroups {
		if subGroup.Name != name {
			continue
		}
		params := lm.NewDeleteDeviceGroupByIDParams()
		params.ID = subGroup.ID
		deleteChildren := true
		params.SetDeleteChildren(&deleteChildren)
		deleteHard := true
		params.SetDeleteHard(&deleteHard)
		_, err := client.LM.DeleteDeviceGroupByID(params)

		return err
	}

	return nil
}

// DeleteGroup deletes a device group with the specified deviceGroupID.
func DeleteGroup(deviceGroup *models.DeviceGroup, client *client.LMSdkGo) error {
	params := lm.NewDeleteDeviceGroupByIDParams()
	params.ID = deviceGroup.ID
	deleteChildren := true
	params.SetDeleteChildren(&deleteChildren)
	deleteHard := true
	params.SetDeleteHard(&deleteHard)
	logrus.Infof("Deleting deviceGroup:\"%s\" ID:\"%d\" ParentID:\"%d\"", *deviceGroup.Name, deviceGroup.ID, deviceGroup.ParentID)
	_, err := client.LM.DeleteDeviceGroupByID(params)

	return err
}

func create(name, appliesTo string, disableAlerting bool, parentID int32, customProperties []*models.NameAndValue, client *client.LMSdkGo) (*models.DeviceGroup, error) {
	params := lm.NewAddDeviceGroupParams()
	params.SetBody(&models.DeviceGroup{ // nolint: exhaustivestruct
		Name:             &name,
		Description:      "A dynamic device group for Kubernetes.",
		ParentID:         parentID,
		AppliesTo:        appliesTo,
		DisableAlerting:  disableAlerting,
		CustomProperties: customProperties,
	})

	restResponse, err := client.LM.AddDeviceGroup(params)
	if err != nil {
		return nil, fmt.Errorf("failed to add device group %q: %w", name, err)
	}

	deviceGroup := restResponse.Payload
	logrus.Infof("Created device group %q with id %d", name, deviceGroup.ID)

	return deviceGroup, nil
}

// GetDeviceGroupPropertyList Fetches device group property list
func GetDeviceGroupPropertyList(lctx *lmctx.LMContext, groupID int32, client *client.LMSdkGo) []*models.EntityProperty {
	log := lmlog.Logger(lctx)
	params := lm.NewGetDeviceGroupPropertyListParams()
	params.SetGid(groupID)
	restResponse, err := client.LM.GetDeviceGroupPropertyList(params)
	if err != nil || restResponse == nil {
		log.Errorf("Failed to fetch device group (id - '%v') property list. Error: %v", groupID, err)

		return []*models.EntityProperty{}
	}

	return restResponse.Payload.Items
}

// UpdateDeviceGroupPropertyByName Updates device group property by name
func UpdateDeviceGroupPropertyByName(lctx *lmctx.LMContext, groupID int32, entityProperty *models.EntityProperty, client *client.LMSdkGo) bool {
	log := lmlog.Logger(lctx)
	params := lm.NewUpdateDeviceGroupPropertyByNameParams()
	params.SetBody(entityProperty)
	params.SetGid(groupID)
	params.SetName(entityProperty.Name)
	restResponse, err := client.LM.UpdateDeviceGroupPropertyByName(params)
	if err != nil || restResponse == nil {
		log.Errorf("Failed to update device group property '%v'. Error: %v", entityProperty.Name, err)

		return false
	}
	log.Debugf("Successfully updated device group property '%v'", entityProperty.Name)

	return true
}

// DeleteDeviceGroupPropertyByName Updates device group property by name
func DeleteDeviceGroupPropertyByName(lctx *lmctx.LMContext, groupID int32, entityProperty *models.EntityProperty, client *client.LMSdkGo) bool {
	log := lmlog.Logger(lctx)
	params := lm.NewDeleteDeviceGroupPropertyByNameParams()
	params.SetGid(groupID)
	params.SetName(entityProperty.Name)
	restResponse, err := client.LM.DeleteDeviceGroupPropertyByName(params)
	if err != nil || restResponse == nil {
		log.Errorf("Failed to delete device group property '%v'. Error: %v", entityProperty.Name, err)

		return false
	}
	log.Debugf("Successfully deleted device group property '%v'", entityProperty.Name)

	return true
}

// AddDeviceGroupProperty Adds new property in device group
func AddDeviceGroupProperty(lctx *lmctx.LMContext, groupID int32, entityProperty *models.EntityProperty, client *client.LMSdkGo) bool {
	log := lmlog.Logger(lctx)
	params := lm.NewAddDeviceGroupPropertyParams()
	params.SetBody(entityProperty)
	params.SetGid(groupID)
	restResponse, err := client.LM.AddDeviceGroupProperty(params)
	if err != nil || restResponse == nil {
		log.Errorf("Failed to add device group property '%v'. Error: %v", entityProperty.Name, err)

		return false
	}
	log.Debugf("Successfully added device group property '%v'", entityProperty.Name)

	return true
}

// updateCustomProperties adds/replaces the custom properties for device group
func updateCustomProperties(lctx *lmctx.LMContext, clusterDeviceGroup *models.DeviceGroup, opts *Options) {
	existingGroupProperties := make(map[string]string)
	customProperties := opts.CustomProperties.Build(clusterDeviceGroup.CustomProperties)
	for _, v := range clusterDeviceGroup.CustomProperties {
		existingGroupProperties[*v.Name] = *v.Value
	}
	for _, v := range customProperties {
		value, exists := existingGroupProperties[*v.Name]
		entityProperty := models.EntityProperty{Name: *v.Name, Value: *v.Value, Type: constants.DeviceGroupCustomType} // nolint: exhaustivestruct
		if !exists {
			// Add property if property does not exists
			AddDeviceGroupProperty(lctx, clusterDeviceGroup.ID, &entityProperty, opts.Client)
		} else if value == "" || value != *v.Value {
			// Update value if it is empty or not equal to value in configmap
			UpdateDeviceGroupPropertyByName(lctx, clusterDeviceGroup.ID, &entityProperty, opts.Client)
		}
	}
}

// GetPropertyValue get device property value by property name
func GetPropertyValue(deviceGroup *models.DeviceGroup, propertyName string) string {
	if deviceGroup == nil {
		return ""
	}
	if len(deviceGroup.CustomProperties) > 0 {
		for _, cp := range deviceGroup.CustomProperties {
			if *cp.Name == propertyName {
				return *cp.Value
			}
		}
	}

	return ""
}
