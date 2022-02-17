package resourcecache

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	lmlog "github.com/logicmonitor/k8s-argus/pkg/log"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	util "github.com/logicmonitor/k8s-argus/pkg/utilities"
	"github.com/logicmonitor/lm-sdk-go/client/lm"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func (rc *ResourceCache) populateCacheStore(lctx *lmctx.LMContext, cmList *corev1.ConfigMapList, selectedDumpID int64, tmpCache *Store) error {
	log := lmlog.Logger(lctx)
	conf, err := config.GetConfig(lctx)
	if err != nil {
		return err
	}
	for _, cm := range cmList.Items {
		dumpID, err2 := strconv.ParseInt(cm.Labels["dumpID"], 10, 64)
		if err2 != nil || (selectedDumpID != -1 && dumpID != selectedDumpID) ||
			conf.ClusterName != cm.Annotations["clusterName"] ||
			cm.Labels["version"] != Version {
			log.Warnf("Failed to load cache chunk with dumpID: %v, selectedDumpID %v, clusterName: %v, version: %v", cm.Labels["dumpID"], selectedDumpID, cm.Annotations["clusterName"], cm.Labels["version"])
			continue
		}
		m := make(map[types.ResourceName][]types.ResourceMeta)
		er := json.Unmarshal([]byte(cm.Data["cache"]), &m)
		if er != nil {
			log.Errorf("Failed to parse stored configmap [%s]: %s", cm.Name, er)

			return er
		}
		tmpCache.AddAll(m)
	}

	return nil
}

func (rc *ResourceCache) selectDumpID(lctx *lmctx.LMContext, cmList *corev1.ConfigMapList) int64 {
	log := lmlog.Logger(lctx)
	// Select latest dump if multiple dump exists
	selectedDumpID := int64(-1)
	for _, cm := range cmList.Items {
		dumpID, err2 := strconv.ParseInt(cm.Labels["dumpID"], 10, 64)
		if err2 == nil && dumpID > selectedDumpID {
			selectedDumpID = dumpID
		}
	}
	if selectedDumpID == -1 {
		log.Warn("No dumpID found in any of the listed configMap, loading all configmaps into cache (may go wrong if multiple dumps are present)")
	}

	return selectedDumpID
}

func (rc *ResourceCache) listAllCacheConfigMaps() (*corev1.ConfigMapList, error) {
	ns, err := config.GetWatchConfig("namespace")
	if err != nil {
		ns = "argus"
	}
	cmList, err := config.GetClientSet().CoreV1().ConfigMaps(ns).List(metav1.ListOptions{
		LabelSelector: labels.Set{"argus": "cache"}.String(),
	})
	if err != nil {
		return cmList, fmt.Errorf("failed to list all configmaps with selector \"argus=cache\": %w", err)
	}
	return cmList, nil
}

// Save saves cache to configmaps
func (rc *ResourceCache) Save(lctx *lmctx.LMContext) error {
	log := lmlog.Logger(lctx)
	rc.flushMU.Lock()
	defer rc.flushMU.Unlock()
	if rc.store.Size() == 0 {
		log.Tracef("store is empty so not storing it")

		return fmt.Errorf("store is empty hence not storing to cm")
	}

	chunks, err := getChunks(rc.store.getMap())
	if err != nil {
		log.Errorf("Failed to marshal cache map to json string or failed to split into chunks %v", err)

		return err
	}
	ns, err := config.GetWatchConfig("namespace")
	if err != nil {
		ns = "argus"
	}
	conf, err := config.GetConfig(lctx)
	if err != nil {
		return err
	}
	dumpID := time.Now().Unix()
	dumpIDStr := fmt.Sprintf("%v", dumpID)
	for idx, chunk := range chunks {
		m := map[string]string{"cache": fmt.Sprintf("%s", chunk)}
		cm := &corev1.ConfigMap{ // nolint: exhaustivestruct
			TypeMeta: metav1.TypeMeta{}, // nolint: exhaustivestruct
			ObjectMeta: metav1.ObjectMeta{ // nolint: exhaustivestruct
				Name:      fmt.Sprintf("cache-%v-%v", dumpID, idx),
				Namespace: ns,
				Labels: map[string]string{
					"argus":       "cache",
					"dumpID":      dumpIDStr,
					"chunkNumber": fmt.Sprintf("%v", idx),
					"version":     Version,
				},
				Annotations: map[string]string{
					"contentSize": fmt.Sprintf("%v", rc.store.Size()),
					"clusterName": conf.ClusterName,
				},
			},
			Data:       m,
			BinaryData: nil,
		}
		_, err1 := config.GetClientSet().CoreV1().ConfigMaps(ns).Create(cm)
		if err1 != nil {
			_, err2 := config.GetClientSet().CoreV1().ConfigMaps(ns).Update(cm)
			if err2 != nil {
				err3 := fmt.Errorf("failed to store cache chunk %v to cm: %w %v", idx, err1, err2)
				log.Errorf("%s", err3)

				return err3
			}
		}
	}
	// Delete previous cache configmaps
	err2 := config.GetClientSet().CoreV1().ConfigMaps(ns).DeleteCollection(
		&metav1.DeleteOptions{}, // nolint: exhaustivestruct
		metav1.ListOptions{ // nolint: exhaustivestruct
			LabelSelector: fmt.Sprintf("argus==cache,dumpID!=%s", dumpIDStr),
		})
	if err2 != nil {
		return fmt.Errorf("failed to delete previous cache configmaps: %w", err2)
	}

	return nil
}

func getChunks(m map[types.ResourceName][]types.ResourceMeta) ([][]byte, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	var result [][]byte
	if len(bytes) > CMMaxBytes {
		n := 1
		m1 := make(map[types.ResourceName][]types.ResourceMeta)
		m2 := make(map[types.ResourceName][]types.ResourceMeta)
		for k, v := range m {
			if n%2 == 0 {
				m1[k] = v
			} else {
				m2[k] = v
			}
			n++
		}
		ch1, err2 := getChunks(m1)
		if err2 != nil {
			return result, err2
		}

		result = append(result, ch1...)
		ch2, err3 := getChunks(m2)
		if err3 != nil {
			return result, err3
		}

		result = append(result, ch2...)

		return result, nil
	}
	result = append(result, bytes)

	return result, nil
}

func (rc *ResourceCache) updateIncrementalCache(lctx *lmctx.LMContext) {
	log := lmlog.Logger(lctx)
	log.Infof("Update incremental cache that might not have dumped into configmap cache..")
	conf, err := config.GetConfig(lctx)
	if err != nil {
		log.Warnf("Failed to get config")
		return
	}
	for {
		clctx := lctx.LMContextWith(map[string]interface{}{constants.PartitionKey: conf.ClusterName})
		params := lm.NewGetImmediateDeviceListByDeviceGroupIDParams()
		params.SetID(*conf.ResourceContainerGroupID)
		// keeping 10 seconds less to cover corner cases - wherein cache was just dumped and objects updated
		// https://www.logicmonitor.com/support/rest-api-developers-guide/v1/resources/get-resources
		filter := fmt.Sprintf("updatedOn>:\"%d\"", rc.dumpID-10) // nolint: gomnd
		params.SetFilter(&filter)
		command := rc.GetImmediateResourceListByResourceGroupIDCommand(clctx, params)
		resp, err := rc.SendReceive(clctx, command)
		if err != nil && util.GetHTTPStatusCodeFromLMSDKError(err) == http.StatusTooManyRequests {
			log.Warnf("Rate limits reached, retrying in %f seconds...", rateLimitBackoffTime.Seconds())
			time.Sleep(rateLimitBackoffTime)

			continue
		}
		count := 0
		if resp != nil && resp.(*lm.GetImmediateDeviceListByDeviceGroupIDOK) != nil {
			for _, resource := range resp.(*lm.GetImmediateDeviceListByDeviceGroupIDOK).Payload.Items {
				if rc.storeDevice(lctx, resource, conf.ClusterName, rc.store) {
					count++
				}
			}
		}
		log.Infof("Updated incremental cache update with %d entries", count)
		break
	}
}
