package cronjob

import (
	"strings"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/devicegroup"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/logicmonitor/k8s-argus/pkg/watch/deployment"
	"github.com/sirupsen/logrus"
	"github.com/vkumbhar94/lm-sdk-go/client"
	"github.com/vkumbhar94/lm-sdk-go/models"
	"k8s.io/client-go/kubernetes"
)

// UpdateTelemetryCron a cron job to update K8s & Helm properties in cluster device group
func UpdateTelemetryCron(base *types.Base) {
	lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"res": "update-telemetry"}))
	c := RegisterCron(lctx, "@midnight", func() { updateTelemetry(lctx, base) })
	c.Start()
}

func updateTelemetry(lctx *lmctx.LMContext, base *types.Base) {
	log := lmlog.Logger(lctx)
	parentID := base.Config.ClusterGroupID
	groupName := constants.ClusterDeviceGroupPrefix + base.Config.ClusterName
	deviceGroup, err := devicegroup.Find(parentID, groupName, base.LMClient)
	if err != nil || deviceGroup == nil {
		log.Errorf("Failed to fetch device group. Error: %v", err)
		return
	}
	updateDeviceGroupK8sAndHelmProperties(lctx, deviceGroup.ID, base.LMClient, base.K8sClient)
}

// updateDeviceGroupK8sAndHelmProperties will fetch existing properties and compare with actual values then update in cluster device group
func updateDeviceGroupK8sAndHelmProperties(lctx *lmctx.LMContext, groupID int32, client *client.LMSdkGo, kubeClient kubernetes.Interface) {
	existingPropertiesMap := getExistingDeviceGroupPropertiesMap(lctx, groupID, client)
	customPropertiesMap := getK8sAndHelmProperties(lctx, kubeClient)

	for k, v := range customPropertiesMap {
		updateProperty(lctx, k, v, groupID, client)

		// update history property
		historyKey := k + constants.HistorySuffix
		updatedHistoryVal := getUpdatedHistoryValue(existingPropertiesMap[historyKey], v)
		updateProperty(lctx, historyKey, updatedHistoryVal, groupID, client)
	}
}

func getExistingDeviceGroupPropertiesMap(lctx *lmctx.LMContext, groupID int32, client *client.LMSdkGo) map[string]string {
	entityProperties := devicegroup.GetDeviceGroupPropertyList(lctx, groupID, client)
	entityPropertiesMap := make(map[string]string)
	for _, property := range entityProperties {
		entityPropertiesMap[property.Name] = property.Value
	}
	return entityPropertiesMap
}

func getK8sAndHelmProperties(lctx *lmctx.LMContext, kubeClient kubernetes.Interface) map[string]string {
	customProperties := make(map[string]string)
	customProperties = getKubernetesVersion(lctx, customProperties, kubeClient)
	customProperties = deployment.GetHelmChartDetailsFromDeployments(lctx, customProperties, kubeClient)
	return customProperties
}

// getKubernetesVersion Fetches Kubernetes version
func getKubernetesVersion(lctx *lmctx.LMContext, customProperties map[string]string, kubeClient kubernetes.Interface) map[string]string {
	log := lmlog.Logger(lctx)
	serverVersion, err := kubeClient.Discovery().ServerVersion()
	if err != nil || serverVersion == nil {
		log.Errorf("Failed to get Kubernetes version. Error: %v", err)
		return customProperties
	}
	cpValue := serverVersion.String()
	customProperties[constants.KubernetesVersionKey] = cpValue
	return customProperties
}

// update the property and if it does not exists then add it
func updateProperty(lctx *lmctx.LMContext, key string, value string, groupID int32, client *client.LMSdkGo) {
	entityProperty := models.EntityProperty{Name: key, Value: value, Type: constants.DeviceGroupCustomType}
	isUpdated := devicegroup.UpdateDeviceGroupPropertyByName(lctx, groupID, &entityProperty, client)
	if !isUpdated {
		devicegroup.AddDeviceGroupProperty(lctx, groupID, &entityProperty, client)
	}
}

func getUpdatedHistoryValue(historyVal, newValue string) string {
	values := []string{}
	if historyVal != "" {
		values = strings.Split(historyVal, constants.PropertySeparator)
	}
	length := len(values)
	// append record if not same as last
	if length == 0 || (length > 0 && values[length-1] != newValue) {
		values = append(values, newValue)
		length = len(values) // calculate length again after adding record
	}
	// retain last 10 records
	if length >= 10 {
		values = values[length-10 : length]
	}
	return strings.Join(values, constants.PropertySeparator)
}
