package cronjob

import (
	"strings"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/devicegroup"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/logicmonitor/k8s-argus/pkg/watch/namespace"
	"github.com/logicmonitor/lm-sdk-go/client"
	"github.com/logicmonitor/lm-sdk-go/models"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const maxHistoryRecordsDefault = 10

// UpdateTelemetryCron a cron job to update K8s & Helm properties in cluster device group
func UpdateTelemetryCron(base *types.Base) {
	lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"res": "update-telemetry"}))
	RegisterFunc(lctx, "@midnight", func() { updateTelemetry(lctx, base) })
}

func updateTelemetry(lctx *lmctx.LMContext, base *types.Base) {
	log := lmlog.Logger(lctx)
	parentID := base.Config.ClusterGroupID
	groupName := util.ClusterGroupName(base.Config.ClusterName)
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
		// update history property
		historyKey := k + constants.HistorySuffix
		updatedHistoryVal := getUpdatedHistoryValue(existingPropertiesMap[historyKey], v)
		updateProperty(lctx, historyKey, updatedHistoryVal, groupID, client)

		// update latest property
		updateProperty(lctx, k, v, groupID, client)
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
	// add Argus app version
	customProperties[constants.ArgusAppVersion] = constants.Version
	customProperties = getKubernetesVersion(lctx, customProperties, kubeClient)
	customProperties = getHelmChartDetailsFromConfigMap(lctx, customProperties, kubeClient)

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

// getHelmChartDetailsFromConfigMap fetches configmap from kubernetes cluster and read annotations
func getHelmChartDetailsFromConfigMap(lctx *lmctx.LMContext, customProperties map[string]string, kubeClient kubernetes.Interface) map[string]string {
	log := lmlog.Logger(lctx)

	// get list of namespace for fetching deployments
	namespaceList := namespace.GetNamespaceList(lctx, kubeClient)

	regex := constants.Chart + " in (" + constants.Argus + ", " + constants.CollectorsetController + ")"
	opts := metav1.ListOptions{ // nolint: exhaustivestruct
		LabelSelector: regex,
	}
	for i := range namespaceList {
		configMapList, err := kubeClient.CoreV1().ConfigMaps(namespaceList[i]).List(opts)
		if err != nil || configMapList == nil {
			log.Errorf("Failed to get the configMap from k8s. Error: %v", err)

			continue
		}
		for i := range configMapList.Items {
			annotations := configMapList.Items[i].GetAnnotations()
			labelVal := configMapList.Items[i].GetLabels()[constants.Chart]
			for key, value := range annotations {
				if key == constants.HelmChart || key == constants.HelmRevision {
					name := labelVal + "." + key
					customProperties[name] = value
				}
			}
		}
	}

	return customProperties
}

// update the property and if it does not exists then add it
func updateProperty(lctx *lmctx.LMContext, key string, value string, groupID int32, client *client.LMSdkGo) {
	entityProperty := models.EntityProperty{Name: key, Value: value, Type: constants.DeviceGroupCustomType} // nolint: exhaustivestruct
	isUpdated := devicegroup.UpdateDeviceGroupPropertyByName(lctx, groupID, &entityProperty, client)
	if !isUpdated {
		devicegroup.AddDeviceGroupProperty(lctx, groupID, &entityProperty, client)
	}
}

func getUpdatedHistoryValue(historyVal, newValue string) string {
	values := make([]string, 0)
	if historyVal != "" {
		values = strings.Split(historyVal, constants.PropertySeparator)
	}
	length := len(values)
	// append record if not same as last
	if length == 0 || values[length-1] != newValue {
		values = append(values, newValue)
		length = len(values) // calculate length again after adding record
	}
	// retain last maxHistoryRecordsDefault records
	if length >= maxHistoryRecordsDefault {
		values = values[length-maxHistoryRecordsDefault : length]
	}

	return strings.Join(values, constants.PropertySeparator)
}
