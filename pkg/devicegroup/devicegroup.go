package devicegroup

import (
	"fmt"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/lm-sdk-go/client"
	"github.com/logicmonitor/lm-sdk-go/client/lm"
	"github.com/logicmonitor/lm-sdk-go/models"
	log "github.com/sirupsen/logrus"
)

const hasCategoryOpen = "hasCategory("
const hasCategoryClose = ")"

// Options are the options for creating a device group.
type Options struct {
	AppliesTo             AppliesToBuilder
	AppliesToDeletedGroup AppliesToBuilder
	Client                *client.LMSdkGo
	Name                  string
	ParentID              int32
	DisableAlerting       bool
	DeleteDevices         bool
}

// AppliesToBuilder is an interface for building an appliesTo string.
type AppliesToBuilder interface {
	HasCategory(string) AppliesToBuilder
	Auto(string) AppliesToBuilder
	And() AppliesToBuilder
	Or() AppliesToBuilder
	Equals(string) AppliesToBuilder
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
	a.value += hasCategoryOpen + fmt.Sprintf(`"%s"`, category) + hasCategoryClose
	return a
}

func (a *appliesToBuilder) Auto(property string) AppliesToBuilder {
	a.value += "auto." + property
	return a
}

func (a *appliesToBuilder) String() string {
	return a.value
}

// Create creates a device group.
func Create(opts *Options) (int32, error) {
	clusterDeviceGroup, err := Find(opts.ParentID, opts.Name, opts.Client)
	if err != nil {
		return 0, err
	}

	if clusterDeviceGroup == nil {
		log.Infof("Could not find device group %q", opts.Name)
		cdg, err := create(opts.Name, opts.AppliesTo.String(), opts.DisableAlerting, opts.ParentID, opts.Client)
		if err != nil {
			return 0, err
		}

		clusterDeviceGroup = cdg
	}

	if !opts.DeleteDevices && opts.AppliesToDeletedGroup != nil {
		deletedDeviceGroup, err := Find(clusterDeviceGroup.ID, constants.DeletedDeviceGroup, opts.Client)
		if err != nil {
			return 0, err
		}
		if deletedDeviceGroup == nil {
			_, err := create(constants.DeletedDeviceGroup, opts.AppliesToDeletedGroup.String(), true, clusterDeviceGroup.ID, opts.Client)
			if err != nil {
				return 0, err
			}
		}

	}

	return clusterDeviceGroup.ID, nil
}

// Find searches for a device group identified by the parent ID and name.
func Find(parentID int32, name string, client *client.LMSdkGo) (*models.DeviceGroup, error) {
	params := lm.NewGetDeviceGroupListParams()
	fields := "name,id,parentId,subGroups"
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

// FindDeviceGroupByID searches for a device group by ID.
func FindDeviceGroupByID(groupID int32, client *client.LMSdkGo) (*models.DeviceGroup, error) {
	params := lm.NewGetDeviceGroupByIDParams()
	params.SetID(groupID)
	fields := "name,id,parentId,subGroups"
	params.SetFields(&fields)
	restResponse, err := client.LM.GetDeviceGroupByID(params)
	if err != nil {
		return nil, fmt.Errorf("failed to get device group (id=%v): %v", groupID, err)
	}

	var deviceGroup *models.DeviceGroup
	if restResponse != nil && restResponse.Payload != nil {
		deviceGroup = restResponse.Payload
	}

	return deviceGroup, nil
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
		log.Infof("Deleting subGroup:\"%s\" from deviceGroup:\"%v\"", subGroup.Name, *deviceGroup.Name)
		_, err := client.LM.DeleteDeviceGroupByID(params)
		return err
	}

	return nil
}

func create(name, appliesTo string, disableAlerting bool, parentID int32, client *client.LMSdkGo) (*models.DeviceGroup, error) {
	params := lm.NewAddDeviceGroupParams()
	params.SetBody(&models.DeviceGroup{
		Name:            &name,
		Description:     "A dynamic device group for Kubernetes.",
		ParentID:        parentID,
		AppliesTo:       appliesTo,
		DisableAlerting: disableAlerting,
	})

	restResponse, err := client.LM.AddDeviceGroup(params)
	if err != nil {
		return nil, fmt.Errorf("failed to add device group %q: %v", name, err)
	}

	deviceGroup := restResponse.Payload
	log.Infof("Created device group %q with id %d", name, deviceGroup.ID)

	return deviceGroup, nil
}
