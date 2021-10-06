package utilities

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/logicmonitor/k8s-argus/pkg/client/csc"
	"github.com/logicmonitor/k8s-argus/pkg/config"
	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/logicmonitor/k8s-argus/pkg/enums"
	"github.com/logicmonitor/k8s-argus/pkg/lmctx"
	"github.com/logicmonitor/k8s-argus/pkg/types"
	"github.com/logicmonitor/lm-sdk-go/client/lm"
	"github.com/logicmonitor/lm-sdk-go/models"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
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
// nolint: cyclop
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
	case constants.K8sAPIVersionBatchV1:
		return clientset.BatchV1().RESTClient()
	case constants.K8sAPIVersionBatchV1Beta1:
		return clientset.BatchV1beta1().RESTClient()
	case constants.K8sAPIVersionExtensionsV1Beta1:
		return clientset.ExtensionsV1beta1().RESTClient()
	case constants.K8sAPIVersionNetworkingV1:
		return clientset.NetworkingV1().RESTClient()
	case constants.K8sAPIVersionNetworkingV1Beta1:
		return clientset.NetworkingV1beta1().RESTClient()
	default:

		return clientset.CoreV1().RESTClient()
	}
}

// GetHTTPStatusCodeFromLMSDKError get code
func GetHTTPStatusCodeFromLMSDKError(err error) int {
	if err == nil {
		return -2
	}
	if errors.Is(err, context.DeadlineExceeded) {
		// 408 client timeout error
		return http.StatusRequestTimeout
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

// GetResourceMetaFromResource get meta
func GetResourceMetaFromResource(resource *models.Device) (types.ResourceMeta, error) {
	labels := make(map[string]string)
	for _, prop := range resource.CustomProperties {
		if strings.HasPrefix(*prop.Name, constants.LabelCustomPropertyPrefix) {
			labels[strings.TrimPrefix(*prop.Name, constants.LabelCustomPropertyPrefix)] = *prop.Value
		}
	}
	categoriesStr := GetResourcePropertyValue(resource, constants.K8sSystemCategoriesPropertyKey)
	categories := strings.Split(categoriesStr, ",")

	return types.ResourceMeta{
		Container:     GetResourcePropertyValue(resource, constants.K8sResourceNamespacePropertyKey),
		LMID:          resource.ID,
		DisplayName:   *resource.DisplayName,
		Name:          *resource.Name,
		Labels:        labels,
		SysCategories: categories,
		UID:           k8stypes.UID(GetResourcePropertyValue(resource, constants.K8sResourceUIDPropertyKey)),
	}, nil
}

// GetResourceMetaFromDeviceGroup get meta
func GetResourceMetaFromDeviceGroup(resourceGroup *models.DeviceGroup) (types.ResourceMeta, error) {
	labels := make(map[string]string)
	for _, prop := range resourceGroup.CustomProperties {
		labels[*prop.Name] = *prop.Value
	}
	categoriesStr := GetResourceGroupPropertyValue(resourceGroup, constants.K8sSystemCategoriesPropertyKey)
	categories := strings.Split(categoriesStr, ",")

	return types.ResourceMeta{
		Container:     fmt.Sprintf("%d", resourceGroup.ParentID),
		LMID:          resourceGroup.ID,
		Name:          *resourceGroup.Name,
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
	return constants.ClusterResourceGroupPrefix + clusterName
}

// GetDisplayName returns desired name
func GetDisplayName(rt enums.ResourceType, meta *metav1.PartialObjectMetadata, conf *config.Config) string {
	// Use full name always to enforce unique name, as freedom to displayname impacts argus performance landing into several lm calls per event
	return fmt.Sprintf("%s-%s", rt.LMName(meta), conf.ClusterName)
}

// GetClusterGroupID avoid call to Santaba and
// returns from cache, cluster root grp is stagnant once created
// nolint: cyclop
func GetClusterGroupID(lctx *lmctx.LMContext, client *types.LMRequester) (int32, error) {
	if clusterGroupID > 0 {
		return clusterGroupID, nil
	}
	conf, err := config.GetConfig(lctx)
	if err != nil {
		return -2, fmt.Errorf("failed to get config: %w", err)
	}
	clusterGroupName := ClusterGroupName(conf.ClusterName)
	clctx := lctx.LMContextWith(map[string]interface{}{constants.PartitionKey: clusterGroupName})
	params := lm.NewGetDeviceGroupByIDParams()
	params.SetID(conf.ClusterGroupID)

	command := client.GetResourceGroupByIDCommand(clctx, params)
	g, err := client.SendReceive(clctx, command)
	if err != nil && GetHTTPStatusCodeFromLMSDKError(err) != http.StatusNotFound {
		return -3, fmt.Errorf("error while fetching cluster resource group %w", err)
	}
	if err != nil && GetHTTPStatusCodeFromLMSDKError(err) == http.StatusNotFound {
		params = lm.NewGetDeviceGroupByIDParams()
		params.SetID(1)

		command = client.GetResourceGroupByIDCommand(clctx, params)
		g, err = client.SendReceive(clctx, command)
	}
	if err != nil {
		return -4, fmt.Errorf("error while fetching root (1) resource group %w", err)
	}

	clusterGrpID := int32(-5)

	for _, sg := range g.(*lm.GetDeviceGroupByIDOK).Payload.SubGroups {
		if sg.Name == clusterGroupName {
			clusterGrpID = sg.ID

			break
		}
	}
	if clusterGrpID <= 0 {
		return clusterGrpID, fmt.Errorf("no child resource group present of name [%s] under resource group [%d]", clusterGroupName, conf.ClusterGroupID)
	}
	clusterGroupID = clusterGrpID

	return clusterGrpID, nil
}

// BuildResource build
func BuildResource(lctx *lmctx.LMContext, c *config.Config, d *models.Device, options ...types.ResourceOption) (*models.Device, error) {
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
			DeviceType:      constants.K8sResourceType,
		}

		for _, option := range options {
			option(d)
		}

		collectorID, err := getCollectorID()
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

func getCollectorID() (int32, error) {
	// when argus is running out cluster on local env, grpc connection with csc of cluster cannot be opened hence
	// returns static id
	if IsLocal() {
		id, err := strconv.ParseInt(os.Getenv("COLLECTOR_ID"), 10, 32)
		if err != nil {
			return 0, fmt.Errorf("could not parse collector id from ENV: COLLECTOR_ID: %w", err)
		}

		return int32(id), nil
	}
	return csc.GetCollectorID()
}

// DoesResourceExistInCacheUtil exists
func DoesResourceExistInCacheUtil(lctx *lmctx.LMContext, rt enums.ResourceType, resourceCache types.ResourceCache, resource *models.Device, softRefresh bool) (types.ResourceMeta, bool) {
	resourceName, err := GetResourceNameFromResource(rt, resource)
	if err != nil {
		return types.ResourceMeta{}, false // nolint: exhaustivestruct
	}

	return resourceCache.Exists(lctx, resourceName, GetResourcePropertyValue(resource, constants.K8sResourceNamespacePropertyKey), softRefresh)
}

// GetResourceNameFromResource get
func GetResourceNameFromResource(rt enums.ResourceType, resource *models.Device) (types.ResourceName, error) {
	return types.ResourceName{
		Name:     GetResourcePropertyValue(resource, constants.K8sResourceNamePropertyKey),
		Resource: rt,
	}, nil
}

func EvaluateExclusion(labels map[string]string) bool {
	for k, v := range labels {
		if k == "logicmonitor/monitoring" && v == "disable" {
			return false
		}
	}

	return true
}

func IsArgusPod(lctx *lmctx.LMContext, rt enums.ResourceType, resource *models.Device) bool {
	if rt == enums.Pods {
		for _, kv := range resource.CustomProperties {
			if *kv.Name == constants.LabelCustomPropertyPrefix+"app" &&
				(*kv.Value == "argus" || *kv.Value == "collectorset-controller") {
				return true
			}
		}
	}
	return false
}

func IsArgusPodCacheMeta(lctx *lmctx.LMContext, rt enums.ResourceType, meta types.ResourceMeta) bool {
	if rt == enums.Pods {
		for k, v := range meta.Labels {
			if k == "app" &&
				(v == "argus" || v == "collectorset-controller") {
				return true
			}
		}
	}
	return false
}
