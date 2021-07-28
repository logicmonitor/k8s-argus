package cronjob

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/logicmonitor/k8s-argus/pkg/aerrors"
	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/resourcegroup"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/logicmonitor/lm-sdk-go/models"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const maxHistoryRecordsDefault = 10

// StartTelemetryCron a cron job to update K8s & Helm properties in cluster resource group
func StartTelemetryCron(resourceCache types.ResourceCache, requester *types.LMRequester) error {
	tu := telemetryUpdater{
		ResourceCache: resourceCache,
		LMRequester:   requester,
		seq:           0,
	}
	conf, err := config.GetConfig()
	if err != nil {
		return err
	}
	_, err = RegisterFunc(*conf.TelemetryCronString, func() {
		tu.seq++
		lctx := lmlog.NewLMContextWith(logrus.WithFields(logrus.Fields{"res": "update-telemetry", "seq": tu.seq}))
		log := lmlog.Logger(lctx)
		log.Infof("Starting telemetry runner execution")
		err := tu.Run(lctx)
		if err != nil {
			log.Errorf("Telemetry update runner failed with error: %s", err)
			return
		}
		log.Infof("Telemetry runner execution completed")
	})
	return err
}

type telemetryUpdater struct {
	types.ResourceCache
	*types.LMRequester
	seq int64
}

func (tu *telemetryUpdater) Run(lctx *lmctx.LMContext) error {
	return tu.run(lctx)
}

// nolint: cyclop
func (tu *telemetryUpdater) run(lctx *lmctx.LMContext) error {
	conf, err := config.GetConfig()
	if err != nil {
		return err
	}
	cacheKey := types.ResourceName{
		Name:     util.ClusterGroupName(conf.ClusterName),
		Resource: enums.Namespaces,
	}
	cacheContainer := fmt.Sprintf("%d", conf.ClusterGroupID)
	meta, ok := tu.ResourceCache.Exists(lctx, cacheKey, cacheContainer, false)
	if !ok {
		err := fmt.Errorf("cluster's root resource group not found")
		return err
	}

	clctx := lctx.LMContextWith(map[string]interface{}{constants.PartitionKey: conf.ClusterName})

	errSlice := make([]error, 0)
	err = tu.audit(clctx, constants.ArgusAppVersion, &meta, func() (string, error) {
		return constants.Version, nil
	})
	if err != nil {
		errSlice = append(errSlice, err)
	}

	err = tu.audit(clctx, constants.KubernetesVersionKey, &meta, getKubernetesVersion)
	if err != nil {
		errSlice = append(errSlice, err)
	}

	m, err := getHelmChartDetailsFromConfigMap()
	if err != nil {
		errSlice = append(errSlice, err)
	} else {
		err = tu.audit(clctx, constants.ArgusHelmChartAuditKey, &meta, func() (string, error) {
			v, ok := m[constants.ArgusHelmChartAuditKey]
			if !ok {
				return "", fmt.Errorf("%s not found", constants.ArgusHelmChartAuditKey)
			}
			return v, nil
		})
		if err != nil {
			errSlice = append(errSlice, err)
		}

		err = tu.audit(clctx, constants.CSCHelmChartAuditKey, &meta, func() (string, error) {
			v, ok := m[constants.CSCHelmChartAuditKey]
			if !ok {
				return "", fmt.Errorf("%s not found", constants.CSCHelmChartAuditKey)
			}
			return v, nil
		})
		if err != nil {
			errSlice = append(errSlice, err)
		}

		err = tu.audit(clctx, constants.ArgusHelmRevisionAuditKey, &meta, func() (string, error) {
			v, ok := m[constants.ArgusHelmRevisionAuditKey]
			if !ok {
				return "", fmt.Errorf("%s not found", constants.ArgusHelmRevisionAuditKey)
			}
			return v, nil
		})
		if err != nil {
			errSlice = append(errSlice, err)
		}

		err = tu.audit(clctx, constants.CSCHelmRevisionAuditKey, &meta, func() (string, error) {
			v, ok := m[constants.CSCHelmRevisionAuditKey]
			if !ok {
				return "", fmt.Errorf("%s not found", constants.CSCHelmRevisionAuditKey)
			}
			return v, nil
		})
		if err != nil {
			errSlice = append(errSlice, err)
		}
	}

	if len(errSlice) > 0 {
		return aerrors.GetMultiError("Failed telemetry update", errSlice...)
	}

	tu.ResourceCache.Set(lctx, cacheKey, meta)
	return nil
}

func (tu *telemetryUpdater) audit(lctx *lmctx.LMContext, auditKey string, meta *types.ResourceMeta, getter func() (string, error)) error {
	historyKey := auditKey + constants.HistorySuffix
	latestVal, err := getter()
	if err != nil {
		return fmt.Errorf("failed to get %s value: %w", auditKey, err)
	}
	if latestVal == "" {
		return fmt.Errorf("value retrieved for %s is empty [%s]", auditKey, latestVal)
	}

	// First update history
	prevVal, ok := meta.Labels[historyKey]
	updatedVal := getUpdatedHistoryValue(prevVal, latestVal)
	errSlice := make([]error, 0)
	if prevVal != updatedVal {
		err := tu.createOrUpdateProperty(lctx, historyKey, updatedVal, meta.LMID, ok)
		if err != nil {
			errSlice = append(errSlice, err)
		} else {
			meta.Labels[historyKey] = updatedVal
		}
	}

	// Update audit prop
	prevVal, ok = meta.Labels[auditKey]
	if prevVal != latestVal {
		err := tu.createOrUpdateProperty(lctx, auditKey, latestVal, meta.LMID, ok)
		if err != nil {
			errSlice = append(errSlice, err)
		} else {
			meta.Labels[auditKey] = latestVal
		}
	}
	if len(errSlice) > 0 {
		return aerrors.GetMultiError(fmt.Sprintf("Failed audit update for %s", auditKey), errSlice...)
	}

	return nil
}

// getKubernetesVersion Fetches Kubernetes version
func getKubernetesVersion() (string, error) {
	serverVersion, err := config.GetClientSet().Discovery().ServerVersion()
	if err != nil {
		return "", err
	}
	return serverVersion.String(), nil
}

// getHelmChartDetailsFromConfigMap fetches configmap from kubernetes cluster and read annotations
func getHelmChartDetailsFromConfigMap() (map[string]string, error) {
	regex := constants.Chart + " in (" + constants.Argus + ", " + constants.CollectorsetController + ")"
	opts := metav1.ListOptions{ // nolint: exhaustivestruct
		LabelSelector: regex,
	}
	configMapList, err := config.GetClientSet().CoreV1().ConfigMaps(metav1.NamespaceAll).List(opts)
	if err != nil || configMapList == nil {
		return nil, fmt.Errorf("failed to get the configMap from k8s. Error: %w", err)
	}
	customProperties := map[string]string{}
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

	return customProperties, nil
}

// update the property and if it does not exists then add it
func (tu *telemetryUpdater) createOrUpdateProperty(lctx *lmctx.LMContext, key string, value string, groupID int32, exists bool) error {
	entityProperty := models.EntityProperty{Name: key, Value: value, Type: constants.ResourceGroupCustomType} // nolint: exhaustivestruct
	if exists {
		_, err := resourcegroup.UpdateResourceGroupPropertyByName(lctx, groupID, &entityProperty, tu.LMRequester)
		if err != nil && util.GetHTTPStatusCodeFromLMSDKError(err) == http.StatusNotFound {
			_, err := resourcegroup.AddResourceGroupProperty(lctx, groupID, &entityProperty, tu.LMRequester)
			return err
		}
		return err
	}
	_, err := resourcegroup.AddResourceGroupProperty(lctx, groupID, &entityProperty, tu.LMRequester)
	// do not update property here on conflict, if call is for add then it might delete old history.
	return err
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
