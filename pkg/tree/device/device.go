package device

import (
	"fmt"

	"github.com/logicmonitor/k8s-argus/pkg/utilities"
	lm "github.com/logicmonitor/lm-sdk-go"
)

// FindByDisplayName searches for a device by it's display name. It will only
// return a device if and only if one device was found, and return nil
// otherwise.
func FindByDisplayName(name string, client *lm.DefaultApi) (*lm.RestDevice, error) {
	filter := fmt.Sprintf("displayName:%s", name)
	restResponse, apiResponse, err := client.GetDeviceList("", -1, 0, filter)
	if _err := utilities.CheckAllErrors(restResponse, apiResponse, err); _err != nil {
		return nil, _err
	}
	if restResponse.Data.Total == 1 {
		return &restResponse.Data.Items[0], nil
	}

	return nil, nil
}

// Add adds a device to a LogicMonitor account.
func Add(d *lm.RestDevice, client *lm.DefaultApi) error {
	restResponse, apiResponse, err := client.AddDevice(*d, false)
	return utilities.CheckAllErrors(restResponse, apiResponse, err)
}

// UpdateAndReplace updatess a device using the 'replace' OpType.
func UpdateAndReplace(d *lm.RestDevice, id int32, client *lm.DefaultApi) error {
	restResponse, apiResponse, err := client.UpdateDevice(*d, id, "replace")
	return utilities.CheckAllErrors(restResponse, apiResponse, err)
}

// UpdateAndReplaceField updatess a device using the 'replace' OpType for a
// specific field of a device.
func UpdateAndReplaceField(d *lm.RestDevice, id int32, field string, client *lm.DefaultApi) error {
	restResponse, apiResponse, err := client.PatchDeviceById(*d, id, "replace", field)
	return utilities.CheckAllErrors(restResponse, apiResponse, err)
}

// Delete deletes a device by device ID.
func Delete(d *lm.RestDevice, client *lm.DefaultApi) error {
	restResponse, apiResponse, err := client.DeleteDevice(d.Id)
	return utilities.CheckAllErrors(restResponse, apiResponse, err)
}

// BuildSystemCategoriesFromLabels formats a system.categories string.
func BuildSystemCategoriesFromLabels(categories string, labels map[string]string) string {
	for k, v := range labels {
		categories += "," + k + "=" + v

	}

	return categories
}
