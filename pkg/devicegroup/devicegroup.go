package devicegroup

import (
	"fmt"
	"strings"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/lm-sdk-go/client"
	"github.com/logicmonitor/lm-sdk-go/client/lm"
	"github.com/logicmonitor/lm-sdk-go/models"
	log "github.com/sirupsen/logrus"
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

type appliesToBuilder struct {
	value string
}

// NewAppliesToBuilder is the builder for appliesTo.
func NewAppliesToBuilder() AppliesToBuilder {
	return &appliesToBuilder{}
}

func (a *appliesToBuilder) And() AppliesToBuilder {
	a.value += " && "
	return a
}

func (a *appliesToBuilder) Or() AppliesToBuilder {
	a.value += " || "
	return a
}
func (a *appliesToBuilder) Equals(val string) AppliesToBuilder {
	a.value += " == " + fmt.Sprintf(`"%s"`, val)
	return a
}

func (a *appliesToBuilder) HasCategory(category string) AppliesToBuilder {
	a.value += hasCategoryOpen + fmt.Sprintf(`"%s"`, category) + closingBracket
	return a
}

func (a *appliesToBuilder) Exists(property string) AppliesToBuilder {
	a.value += existsOpen + fmt.Sprintf(`"%s"`, property) + closingBracket
	return a
}

func (a *appliesToBuilder) Auto(property string) AppliesToBuilder {
	a.value += "auto." + property
	return a
}

func (a *appliesToBuilder) Custom(property string) AppliesToBuilder {
	a.value += property
	return a
}

func (a *appliesToBuilder) String() string {
	return a.value
}

// PropertyBuilder is an interface for building properties
type PropertyBuilder interface {
	Add(string, string) PropertyBuilder
	AddDeleteAfter(string) PropertyBuilder
	GetDeleteAfter() string
	Build() []*models.NameAndValue
}

type propertyBuilder struct {
	deleteAfter string
	properties  map[string]string
}

// NewPropertyBuilder is the builder for properties
func NewPropertyBuilder() PropertyBuilder {
	return &propertyBuilder{deleteAfter: "", properties: make(map[string]string)}
}

func (p *propertyBuilder) Add(key string, value string) PropertyBuilder {
	p.properties[key] = value
	return p
}

func (p *propertyBuilder) AddDeleteAfter(deleteAfter string) PropertyBuilder {
	p.deleteAfter = deleteAfter
	return p
}

func (p *propertyBuilder) GetDeleteAfter() string {
	return p.deleteAfter
}

func (p *propertyBuilder) Build() []*models.NameAndValue {
	props := []*models.NameAndValue{}
	for k, v := range p.properties {
		props = append(props, &models.NameAndValue{Name: &k, Value: &v})
	}
	return props
}

// Create creates a device group.
func Create(opts *Options) (int32, error) {
	lctx := lmlog.NewLMContextWith(log.WithFields(log.Fields{"res": "create-device-group", "device_group_name": opts.Name}))
	log := lmlog.Logger(lctx)
	clusterDeviceGroup, err := Find(opts.ParentID, opts.Name, opts.Client)
	if err != nil {
		return 0, err
	}

	if clusterDeviceGroup == nil {
		log.Infof("Could not find device group %q", opts.Name)
		cdg, err := create(opts.Name, opts.AppliesTo.String(), opts.DisableAlerting, opts.ParentID, opts.CustomProperties.Build(), opts.Client)
		if err != nil {
			return 0, err
		}

		clusterDeviceGroup = cdg
	} else {
		checkDeleteAfterDurationProperty(lctx, clusterDeviceGroup, opts)
	}

	if !opts.DeleteDevices && opts.AppliesToDeletedGroup != nil {
		deletedDeviceGroup, err := Find(clusterDeviceGroup.ID, constants.DeletedDeviceGroup, opts.Client)
		if err != nil {
			return 0, err
		}
		if deletedDeviceGroup == nil {
			_, err := create(constants.DeletedDeviceGroup, opts.AppliesToDeletedGroup.String(), true, clusterDeviceGroup.ID, NewPropertyBuilder().Build(), opts.Client)
			if err != nil {
				return 0, err
			}
		}

	}

	err = createConflictDynamicGroup(opts, clusterDeviceGroup.ID)
	if err != nil {
		return 0, err
	}

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
	filter := fmt.Sprintf("name:\"%s\"", name)
	params.SetFilter(&filter)
	restResponse, err := client.LM.GetDeviceGroupList(params)
	if err != nil {
		return nil, fmt.Errorf("failed to get device group list when searching for %q: %v", name, err)
	}

	log.Debugf("%#v", restResponse)

	var deviceGroup *models.DeviceGroup
	for _, d := range restResponse.Payload.Items {
		if d.ParentID == parentID {
			log.Infof("Found device group %q with id %d", name, d.ID)
			deviceGroup = d
			break
		}
	}

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
		log.Warnf("Failed looking up device group for node role %q: %v", name, err)
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
		log.Warnf("Failed to get device group (id=%v): %v", groupID, err)
		return false
	}

	log.Debugf("%#v", restResponse)

	if restResponse.Payload != nil && restResponse.Payload.ID == groupID {
		return true
	}

	return false
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
	log.Infof("Deleting deviceGroup:\"%s\" ID:\"%d\" ParentID:\"%d\"", *deviceGroup.Name, deviceGroup.ID, deviceGroup.ParentID)
	_, err := client.LM.DeleteDeviceGroupByID(params)
	return err
}

func create(name, appliesTo string, disableAlerting bool, parentID int32, customProperties []*models.NameAndValue, client *client.LMSdkGo) (*models.DeviceGroup, error) {
	params := lm.NewAddDeviceGroupParams()
	params.SetBody(&models.DeviceGroup{
		Name:             &name,
		Description:      "A dynamic device group for Kubernetes.",
		ParentID:         parentID,
		AppliesTo:        appliesTo,
		DisableAlerting:  disableAlerting,
		CustomProperties: customProperties,
	})

	restResponse, err := client.LM.AddDeviceGroup(params)
	if err != nil {
		return nil, fmt.Errorf("failed to add device group %q: %v", name, err)
	}

	deviceGroup := restResponse.Payload
	log.Infof("Created device group %q with id %d", name, deviceGroup.ID)

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

func checkDeleteAfterDurationProperty(lctx *lmctx.LMContext, clusterDeviceGroup *models.DeviceGroup, opts *Options) {
	if opts.CustomProperties.GetDeleteAfter() == "" {
		return
	}
	propertyExists := false
	propertyValueEmpty := false
	for _, prop := range clusterDeviceGroup.CustomProperties {
		if strings.EqualFold(*prop.Name, constants.K8sResourceDeleteAfterDurationPropertyKey) {
			propertyExists = true
			if *prop.Value == "" {
				propertyValueEmpty = true
			}
			break
		}
	}
	entityProperty := models.EntityProperty{Name: constants.K8sResourceDeleteAfterDurationPropertyKey, Value: opts.CustomProperties.GetDeleteAfter(), Type: constants.DeviceGroupCustomType}
	if !propertyExists {
		AddDeviceGroupProperty(lctx, clusterDeviceGroup.ID, &entityProperty, opts.Client)
	}
	if propertyValueEmpty {
		UpdateDeviceGroupPropertyByName(lctx, clusterDeviceGroup.ID, &entityProperty, opts.Client)
	}
}
