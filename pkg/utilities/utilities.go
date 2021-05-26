package utilities

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/devicecache/cache"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/logicmonitor/lm-sdk-go/client"
	"github.com/logicmonitor/lm-sdk-go/client/lm"
	"github.com/logicmonitor/lm-sdk-go/models"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var clusterGroupID = int32(-1)

// GetLabelByPrefix takes a list of labels returns the first label matching the specified prefix
func GetLabelByPrefix(prefix string, labels map[string]string) (string, string) {
	for k, v := range labels {
		if match, err := regexp.MatchString("^"+prefix, k); match {
			if err != nil {
				continue
			}

			return k, v
		}
	}

	return "", ""
}

// GetLabelsByPrefix takes a list of labels returns the first label matching the specified prefix
func GetLabelsByPrefix(prefix string, labels map[string]string) map[string]string {
	m := make(map[string]string)

	for k, v := range labels {
		if match, err := regexp.MatchString("^"+prefix, k); match {
			if err != nil {
				continue
			}
			m[k] = v
		}
	}

	return m
}

// GetShortUUID returns short ids. introduced this util function to start for traceability of events and its logs
func GetShortUUID() uint32 {
	return uuid.New().ID()
}

// GetK8sRESTClient get the K8s RESTClient by apiVersion, use the default V1 version if there is no match
func GetK8sRESTClient(clientset *kubernetes.Clientset, apiVersion string) rest.Interface {
	switch apiVersion {
	case constants.K8sAPIVersionV1:

		return clientset.CoreV1().RESTClient()
	case constants.K8sAPIVersionAppsV1beta2:

		return clientset.AppsV1beta2().RESTClient()
	case constants.K8sAPIVersionAppsV1:

		return clientset.AppsV1().RESTClient()
	case constants.K8sAutoscalingV1:

		return clientset.AutoscalingV1().RESTClient()
	default:

		return clientset.CoreV1().RESTClient()
	}
}

// GetHTTPStatusCodeFromLMSDKError get code
func GetHTTPStatusCodeFromLMSDKError(err error) int {
	if err == nil {
		return -2
	}
	errRegex := regexp.MustCompile(`(?P<api>\[.*\])\[(?P<code>\d+)\].*`)
	matches := errRegex.FindStringSubmatch(err.Error())
	if len(matches) < 3 { // nolint: gomnd
		return -1
	}

	code, err := strconv.Atoi(matches[2])
	if err != nil {
		return -1
	}

	return code
}

// GetCurrentFunctionName get current
func GetCurrentFunctionName() string {
	// Skip GetCurrentFunctionName

	return GetNthCallerFunctionName(2) // nolint: gomnd
}

// GetCallerFunctionName get caller
func GetCallerFunctionName() string {
	// Skip GetCallerFunctionName and the function to get the caller of

	return GetNthCallerFunctionName(3) // nolint: gomnd
}

// GetNthCallerFunctionName get nth caller
func GetNthCallerFunctionName(n int) string {
	// Skip GetCallerFunctionName and the function to get the caller of

	return strings.TrimPrefix(getFrame(n).Function, "github.com/logicmonitor/k8s-argus/")
}

// Referenced here: https://play.golang.org/p/cv-SpkvexuM
func getFrame(skipFrames int) runtime.Frame {
	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
	targetFrameIndex := skipFrames + 2 // nolint: gomnd

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
	programCounters := make([]uintptr, targetFrameIndex+2) // nolint: gomnd
	frame := runtime.Frame{Function: "unknown"}            // nolint: exhaustivestruct

	n := runtime.Callers(0, programCounters)
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])

		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}

	return frame
}

// GetResourceMetaFromDevice get meta
func GetResourceMetaFromDevice(device *models.Device) (cache.ResourceMeta, error) {
	labels := make(map[string]string)
	for _, prop := range device.CustomProperties {
		if strings.HasPrefix(*prop.Name, constants.LabelCustomPropertyPrefix) {
			labels[strings.TrimPrefix(*prop.Name, constants.LabelCustomPropertyPrefix)] = *prop.Value
		}
	}
	categoriesStr := GetPropertyValue(device, constants.K8sSystemCategoriesPropertyKey)
	categories := strings.Split(categoriesStr, ",")

	return cache.ResourceMeta{
		Container:     GetPropertyValue(device, constants.K8sDeviceNamespacePropertyKey),
		LMID:          device.ID,
		DisplayName:   *device.DisplayName,
		Labels:        labels,
		SysCategories: categories,
	}, nil
}

// IsLocal check is local
func IsLocal() bool {
	return os.Getenv(constants.IsLocal) == "true"
}

// ClusterGroupName cluster group name
func ClusterGroupName(clusterName string) string {
	return constants.ClusterDeviceGroupPrefix + clusterName
}

// GetDisplayNameNew returns desired name
// TODO: eventually get rid of these options, include namespace and include cluster name, always generate a unique name to reduce complexity
func GetDisplayNameNew(rt enums.ResourceType, meta *metav1.ObjectMeta, conf *config.Config) string {
	// Use full name always to enforce unique name, as freedom to displayname impacts argus performance landing into several lm calls per event
	return fmt.Sprintf("%s-%s", rt.LMName(meta), conf.ClusterName)
	// return getDisplayNameAsPerSettings(rt, meta, conf)
}

// getDisplayNameAsPerSettings just kept for backup
func getDisplayNameAsPerSettings(rt enums.ResourceType, meta *metav1.ObjectMeta, conf *config.Config) string {
	if conf.FullDisplayNameIncludeClusterName {
		return fmt.Sprintf("%s-%s", rt.LMName(meta), conf.ClusterName)
	}
	if conf.FullDisplayNameIncludeNamespace {
		return rt.LMName(meta)
	}
	s := enums.ShortResourceType(rt)

	return fmt.Sprintf("%s-%s", meta.Name, s.String())
}

// GetClusterGroupID avoid call to Santaba and
// returns from cache, cluster root grp is stagnant once created
func GetClusterGroupID(lctx *lmctx.LMContext, client *client.LMSdkGo) int32 {
	if clusterGroupID != -1 {
		return clusterGroupID
	}
	conf, err := config.GetConfig()
	if err != nil {
		logrus.Errorf("Failed to get config")
		return -1
	}
	params := lm.NewGetDeviceGroupByIDParams()
	params.SetID(conf.ClusterGroupID)

	g, err := client.LM.GetDeviceGroupByID(params)
	if err != nil {
		logrus.Errorf("Error while fetching cluster device group %v", err)

		return -1
	}

	clusterGroupName := ClusterGroupName(conf.ClusterName)
	clusterGrpID := int32(-1)

	for _, sg := range g.Payload.SubGroups {
		if sg.Name == clusterGroupName {
			clusterGrpID = sg.ID

			break
		}
	}
	clusterGroupID = clusterGrpID

	return clusterGrpID
}

// BuildDevice build
func BuildDevice(lctx *lmctx.LMContext, c *config.Config, d *models.Device, options ...types.DeviceOption) (*models.Device, error) {
	if d == nil {
		hostGroupIds := fmt.Sprintf("%d", *c.ResourceContainerGroupID)
		propertyName := constants.K8sClusterNamePropertyKey
		// use the copy value
		clusterName := c.ClusterName
		d = &models.Device{ // nolint: exhaustivestruct
			CustomProperties: []*models.NameAndValue{
				{
					Name:  &propertyName,
					Value: &clusterName,
				},
			},
			DisableAlerting: c.DisableAlerting,
			HostGroupIds:    &hostGroupIds,
			DeviceType:      constants.K8sDeviceType,
		}

		for _, option := range options {
			option(d)
		}

		collectorID, err := GetCollectorID()
		if err != nil {
			return d, &types.GetCollectorIDError{Err: err}
		}
		d.PreferredCollectorID = &collectorID
	} else {
		for _, option := range options {
			option(d)
		}
	}

	return d, nil
}

// DoesDeviceExistInCacheUtil exists
func DoesDeviceExistInCacheUtil(lctx *lmctx.LMContext, resource enums.ResourceType, resourceCache types.ResourceCache, device *models.Device) (cache.ResourceMeta, bool) {
	resourceName, err := GetResourceNameFromDevice(resource, device)
	if err != nil {
		return cache.ResourceMeta{}, false // nolint: exhaustivestruct
	}

	return resourceCache.Exists(lctx, resourceName, GetPropertyValue(device, constants.K8sDeviceNamespacePropertyKey))
}

// GetResourceNameFromDevice get
func GetResourceNameFromDevice(resource enums.ResourceType, device *models.Device) (cache.ResourceName, error) {
	return cache.ResourceName{
		Name:     GetPropertyValue(device, constants.K8sDeviceNamePropertyKey),
		Resource: resource,
	}, nil
}
