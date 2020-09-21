package cronjob

import (
	"strings"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/devicegroup"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/logicmonitor/k8s-argus/pkg/watch/deployment"
	"github.com/logicmonitor/lm-sdk-go/client"
	"github.com/logicmonitor/lm-sdk-go/models"
	"github.com/sirupsen/logrus"
	"gopkg.in/robfig/cron.v2"
	"k8s.io/client-go/kubernetes"
)

// UpdateTelemetryCron a cron job to update K8s & Helm properties in cluster device group
func UpdateTelemetryCron(base *types.Base) {
	lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"res": "update-telemetry"}))
	log := lmlog.Logger(lctx)
	c := cron.New()
	// scheduling is done in the machine's local time zone at midnight
	// _, err := cron.AddFunc("@midnight", func() {
	_, err := c.AddFunc("@every 0h3m0s", func() {
		updateTelemetry(lctx, base)
	})
	if err != nil {
		log.Errorf("Failed to create cron job. Error: %v", err)
	}
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
		historyKey := k + constants.HistorySuffix
		historyVal, historyValExists := existingPropertiesMap[historyKey]
		value, propertyExists := existingPropertiesMap[k]

		if !propertyExists {
			entityProperty := models.EntityProperty{Name: k, Value: v, Type: constants.DeviceGroupCustomType}
			devicegroup.AddDeviceGroupProperty(lctx, groupID, &entityProperty, client)
			newValue := ""
			if historyValExists {
				newValue = getUpdatedHistoryValue(historyVal, v)
				updateProperty(lctx, historyKey, newValue, groupID, client)
			}
		} else if value == "" || value != v {
			updateProperty(lctx, k, v, groupID, client)
			newValue := ""
			if !historyValExists {
				newValue = getNewHistoryValue(value, v)
				entityProperty := models.EntityProperty{Name: historyKey, Value: newValue, Type: constants.DeviceGroupCustomType}
				devicegroup.AddDeviceGroupProperty(lctx, groupID, &entityProperty, client)
			} else {
				newValue = getUpdatedHistoryValue(historyVal, v)
				updateProperty(lctx, historyKey, newValue, groupID, client)
			}
		}
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

func getNewHistoryValue(value, v string) string {
	values := []string{}
	if value != "" {
		values = append(values, value)
	}
	values = append(values, v)
	return strings.Join(values, constants.PropertySeparator)
}

func getUpdatedHistoryValue(historyVal, v string) string {
	values := []string{}
	if historyVal != "" {
		values = strings.Split(historyVal, constants.PropertySeparator)
	}
	length := len(values)
	// to retain last 10 records, trim last 9 record then append 10th record if not same as last
	if length > 9 {
		values = values[length-9 : length]
		length = len(values) // calculated length again after removing elements from slice
	}
	if length == 0 || (length > 0 && values[length-1] != v) {
		values = append(values, v)
	}
	return strings.Join(values, constants.PropertySeparator)
}
