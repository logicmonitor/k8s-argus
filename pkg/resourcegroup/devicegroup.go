package resourcegroup

import (
	"net/http"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/logicmonitor/lm-sdk-go/client/lm"
	"github.com/logicmonitor/lm-sdk-go/models"
)

// ExistsByID returns true if we could get the group by id
func ExistsByID(lctx *lmctx.LMContext, groupID int32, client *types.LMRequester) (bool, error) {
	conf, err := config.GetConfig(lctx)
	if err != nil {
		return false, err
	}
	clusterGroupName := util.ClusterGroupName(conf.ClusterName)
	clctx := lctx.LMContextWith(map[string]interface{}{constants.PartitionKey: clusterGroupName})
	params := lm.NewGetDeviceGroupByIDParams()
	params.SetID(groupID)
	fields := "name,id"
	params.SetFields(&fields)
	command := client.GetResourceGroupByIDCommand(clctx, params)
	restResponse, err := client.SendReceive(clctx, command)
	if err != nil && util.GetHTTPStatusCodeFromLMSDKError(err) == http.StatusNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return restResponse.(*lm.GetDeviceGroupByIDOK).Payload != nil && restResponse.(*lm.GetDeviceGroupByIDOK).Payload.ID == groupID, nil
}

// GetByID returns true if we could get the group by id
func GetByID(lctx *lmctx.LMContext, groupID int32, requester *types.LMRequester) (*models.DeviceGroup, error) {
	conf, err := config.GetConfig(lctx)
	if err != nil {
		return nil, err
	}
	clctx := lctx.LMContextWith(map[string]interface{}{constants.PartitionKey: conf.ClusterName})
	params := lm.NewGetDeviceGroupByIDParams()
	params.SetID(groupID)
	command := requester.GetResourceGroupByIDCommand(clctx, params)
	restResponse, err := requester.SendReceive(clctx, command)
	if err != nil {
		return nil, err
	}

	return restResponse.(*lm.GetDeviceGroupByIDOK).Payload, nil
}

func GetClusterGroupProperty(lctx *lmctx.LMContext, name string, client *types.LMRequester) string {
	log := lmlog.Logger(lctx)
	clusterGroupID, err := util.GetClusterGroupID(lctx, client)
	if err != nil {
		log.Error(err.Error())
		return ""
	}
	resourceGroup, err := GetByID(lctx, clusterGroupID, client)
	if err != nil {
		log.Errorf("Failed to fetch cluster group [%d]: %s", clusterGroupID, err)
		return ""
	}
	return GetPropertyValue(resourceGroup, name)
}

// UpdateResourceGroupPropertyByName Updates resource group property by name
func UpdateResourceGroupPropertyByName(lctx *lmctx.LMContext, groupID int32, entityProperty *models.EntityProperty, requester *types.LMRequester) (interface{}, error) {
	params := lm.NewUpdateDeviceGroupPropertyByNameParams()
	params.SetBody(entityProperty)
	params.SetGid(groupID)
	params.SetName(entityProperty.Name)
	command := requester.UpdateResourceGroupPropertyByNameCommand(lctx, params)
	return requester.SendReceive(lctx, command)
}

// DeleteResourceGroupPropertyByName Updates resource group property by name
func DeleteResourceGroupPropertyByName(lctx *lmctx.LMContext, groupID int32, entityProperty *models.EntityProperty, lmrequester *types.LMRequester) bool {
	log := lmlog.Logger(lctx)
	params := lm.NewDeleteDeviceGroupPropertyByNameParams()
	params.SetGid(groupID)
	params.SetName(entityProperty.Name)
	command := lmrequester.DeleteResourceGroupPropertyByNameCommand(lctx, params)
	_, err := lmrequester.SendReceive(lctx, command)
	statusCode := util.GetHTTPStatusCodeFromLMSDKError(err)
	if err != nil && statusCode != http.StatusNotFound {
		log.Errorf("Failed to delete resource group property '%v'. Error: %v", entityProperty.Name, err)

		return false
	}
	log.Debugf("Successfully deleted resource group property '%v'", entityProperty.Name)

	return true
}

// AddResourceGroupProperty Adds new property in resource group
func AddResourceGroupProperty(lctx *lmctx.LMContext, groupID int32, entityProperty *models.EntityProperty, requester *types.LMRequester) (interface{}, error) {
	params := lm.NewAddDeviceGroupPropertyParams()
	params.SetBody(entityProperty)
	params.SetGid(groupID)
	command := requester.AddResourceGroupPropertyCommand(lctx, params)
	return requester.SendReceive(lctx, command)
}

// GetPropertyValue get resource property value by property name
func GetPropertyValue(resourceGroup *models.DeviceGroup, propertyName string) string {
	if resourceGroup == nil {
		return ""
	}
	if len(resourceGroup.CustomProperties) > 0 {
		for _, cp := range resourceGroup.CustomProperties {
			if *cp.Name == propertyName {
				return *cp.Value
			}
		}
	}

	return ""
}
