package devicegroup

import (
	"fmt"
	"net/url"

	"github.com/logicmonitor/k8s-argus/pkg/utilities"
	lm "github.com/logicmonitor/lm-sdk-go"
	log "github.com/sirupsen/logrus"
)

const hasCategoryOpen = "hasCategory("
const hasCategoryClose = ")"

// Options are the options for creating a device group.
type Options struct {
	Name            string
	ParentID        int32
	DisableAlerting bool
	AppliesTo       AppliesToBuilder
	Client          *lm.DefaultApi
}

// AppliesToBuilder is an interface for building and appliesTo string.
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
	a.value += " == " + val
	return a
}

func (a *appliesToBuilder) HasCategory(category string) AppliesToBuilder {
	a.value += hasCategoryOpen + category + hasCategoryClose
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
	clusterDeviceGroup, err := find(opts.ParentID, opts.Name, opts.Client)
	if err != nil {
		return 0, err
	}

	if clusterDeviceGroup == nil {
		cdg, err := create(opts.Name, opts.AppliesTo.String(), opts.DisableAlerting, opts.ParentID, opts.Client)
		if err != nil {
			return 0, err
		}

		clusterDeviceGroup = cdg
	}

	return clusterDeviceGroup.Id, nil
}

func find(parentID int32, name string, client *lm.DefaultApi) (*lm.RestDeviceGroup, error) {
	filter := fmt.Sprintf("name:%s", url.QueryEscape(name))
	restResponse, apiResponse, err := client.GetDeviceGroupList("name,id,parentId", -1, 0, filter)
	if _err := utilities.CheckAllErrors(restResponse, apiResponse, err); _err != nil {
		log.Errorf("Failed to find device group %q: %v", name, _err)
	}

	log.Debugf("%#v", restResponse)

	var deviceGroup *lm.RestDeviceGroup
	for _, d := range restResponse.Data.Items {
		if d.ParentId == parentID {
			log.Infof("Found device group %q with id %d", name, parentID)
			deviceGroup = &d

			return deviceGroup, nil
		}
	}

	return nil, fmt.Errorf("Unexpected error when searching for device group %q", name)
}

func create(name, appliesTo string, disableAlerting bool, parentID int32, client *lm.DefaultApi) (*lm.RestDeviceGroup, error) {
	restResponse, apiResponse, err := client.AddDeviceGroup(lm.RestDeviceGroup{
		Name:            name,
		Description:     "A dynamic device group for Kubernetes.",
		ParentId:        parentID,
		AppliesTo:       appliesTo,
		DisableAlerting: disableAlerting,
	})
	if e := utilities.CheckAllErrors(restResponse, apiResponse, err); e != nil {
		return nil, fmt.Errorf("Failed to add device group: %v", e)
	}

	deviceGroup := &restResponse.Data
	log.Infof("Created device group with id %d", deviceGroup.Id)

	return deviceGroup, nil
}
