package controller

import (
	"fmt"
	"strings"

	crv1alpha1 "github.com/logicmonitor/k8s-collectorset-controller/pkg/apis/v1alpha1"
	"github.com/logicmonitor/k8s-collectorset-controller/pkg/constants"
	"github.com/logicmonitor/lm-sdk-go/client"
	"github.com/logicmonitor/lm-sdk-go/client/lm"
	"github.com/logicmonitor/lm-sdk-go/models"
	log "github.com/sirupsen/logrus"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	clientset "k8s.io/client-go/kubernetes"
)

// CreateOrUpdateCollectorSet creates a replicaset for each collector in
// a CollectorSet
func CreateOrUpdateCollectorSet(collectorset *crv1alpha1.CollectorSet, controller *Controller) ([]int32, error) {
	groupID := collectorset.Spec.GroupID
	if groupID == 0 || !checkCollectorGroupExistsByID(controller.LogicmonitorClient, groupID) {
		groupName := constants.ClusterCollectorGroupPrefix + collectorset.Spec.ClusterName
		log.Infof("Group name is %s", groupName)

		newGroupID, err := getCollectorGroupID(controller.LogicmonitorClient, groupName, collectorset)
		if err != nil {
			return nil, err
		}
		log.Infof("Adding collector group %q with ID %d", strings.Title(groupName), newGroupID)
		groupID = newGroupID
	}

	ids, err := getCollectorIDs(controller.LogicmonitorClient, groupID, collectorset)
	if err != nil {
		return nil, err
	}

	secretIsOptional := false
	collectorSize := strings.ToLower(collectorset.Spec.Size)
	log.Infof("Collector size is %s", collectorSize)

	statefulset := appsv1beta1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1beta1",
			Kind:       "StatefulSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      collectorset.Name,
			Namespace: collectorset.Namespace,
			Labels: map[string]string{
				"logicmonitor.com/collectorset": collectorset.Name,
			},
		},
		Spec: appsv1beta1.StatefulSetSpec{
			Replicas: collectorset.Spec.Replicas,
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: collectorset.Namespace,
					Labels: map[string]string{
						"logicmonitor.com/collectorset": collectorset.Name,
					},
				},
				Spec: apiv1.PodSpec{
					ServiceAccountName: constants.CollectorServiceAccountName,
					Affinity: &apiv1.Affinity{
						PodAntiAffinity: &apiv1.PodAntiAffinity{
							RequiredDuringSchedulingIgnoredDuringExecution: []apiv1.PodAffinityTerm{
								{
									LabelSelector: &metav1.LabelSelector{
										MatchLabels: map[string]string{
											"logicmonitor.com/collectorset": collectorset.Name,
										},
									},
									TopologyKey: "kubernetes.io/hostname",
								},
							},
						},
					},
					Containers: []apiv1.Container{
						{
							Name:            "collector",
							Image:           "logicmonitor/collector:latest",
							ImagePullPolicy: apiv1.PullAlways,
							Env: []apiv1.EnvVar{
								{
									Name: "account",
									ValueFrom: &apiv1.EnvVarSource{
										SecretKeyRef: &apiv1.SecretKeySelector{
											LocalObjectReference: apiv1.LocalObjectReference{
												Name: constants.CollectorsetControllerSecretName,
											},
											Key:      "account",
											Optional: &secretIsOptional,
										},
									},
								},
								{
									Name: "access_id",
									ValueFrom: &apiv1.EnvVarSource{
										SecretKeyRef: &apiv1.SecretKeySelector{
											LocalObjectReference: apiv1.LocalObjectReference{
												Name: constants.CollectorsetControllerSecretName,
											},
											Key:      "accessID",
											Optional: &secretIsOptional,
										},
									},
								},
								{
									Name: "access_key",
									ValueFrom: &apiv1.EnvVarSource{
										SecretKeyRef: &apiv1.SecretKeySelector{
											LocalObjectReference: apiv1.LocalObjectReference{
												Name: constants.CollectorsetControllerSecretName,
											},
											Key:      "accessKey",
											Optional: &secretIsOptional,
										},
									},
								},
								{
									Name:  "kubernetes",
									Value: "true",
								},
								{
									Name:  "collector_size",
									Value: collectorSize,
								},
								{
									Name:  "collector_version",
									Value: fmt.Sprint(collectorset.Spec.CollectorVersion), //the default value is 0, santaba will assign the latest version
								},
								{
									Name:  "use_ea",
									Value: fmt.Sprint(collectorset.Spec.UseEA), //the default value is false, santaba will assign the latest GD version
								},
								{
									Name:  "COLLECTOR_IDS",
									Value: strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ","), "[]"),
								},
							},
							Resources: getResourceRequirements(collectorSize),
						},
					},
				},
			},
			UpdateStrategy: appsv1beta1.StatefulSetUpdateStrategy{
				Type: appsv1beta1.RollingUpdateStatefulSetStrategyType,
			},
			PodManagementPolicy: appsv1beta1.ParallelPodManagement,
		},
	}

	setProxyConfiguration(collectorset, &statefulset)

	if _, _err := controller.Clientset.AppsV1beta1().StatefulSets(statefulset.ObjectMeta.Namespace).Create(&statefulset); _err != nil {
		if !apierrors.IsAlreadyExists(_err) {
			return nil, _err
		}
		if _, _err := controller.Clientset.AppsV1beta1().StatefulSets(statefulset.ObjectMeta.Namespace).Update(&statefulset); _err != nil {
			return nil, _err
		}
	}

	collectorset.Status.IDs = ids

	err = updateCollectors(controller.LogicmonitorClient, ids)
	if err != nil {
		log.Warnf("Failed to set collector backup agents: %v", err)
	}
	return collectorset.Status.IDs, nil
}

func setProxyConfiguration(collectorset *crv1alpha1.CollectorSet, statefulset *appsv1beta1.StatefulSet) {
	if collectorset.Spec.ProxyURL == "" {
		return
	}
	container := &statefulset.Spec.Template.Spec.Containers[0]
	container.Env = append(container.Env,
		apiv1.EnvVar{
			Name:  "proxy_url",
			Value: collectorset.Spec.ProxyURL,
		},
	)
	if collectorset.Spec.SecretName != "" {
		secretIsOptionalTrue := true
		container.Env = append(container.Env,
			apiv1.EnvVar{
				Name: "proxy_user",
				ValueFrom: &apiv1.EnvVarSource{
					SecretKeyRef: &apiv1.SecretKeySelector{
						LocalObjectReference: apiv1.LocalObjectReference{
							Name: collectorset.Spec.SecretName,
						},
						Key:      "proxyUser",
						Optional: &secretIsOptionalTrue,
					},
				},
			},
			apiv1.EnvVar{
				Name: "proxy_pass",
				ValueFrom: &apiv1.EnvVarSource{
					SecretKeyRef: &apiv1.SecretKeySelector{
						LocalObjectReference: apiv1.LocalObjectReference{
							Name: collectorset.Spec.SecretName,
						},
						Key:      "proxyPass",
						Optional: &secretIsOptionalTrue,
					},
				},
			},
		)
	}
}

func updateCollectors(client *client.LMSdkGo, ids []int32) error {
	// if there is only one collector, there will be no backup for it
	if len(ids) < 2 {
		return nil
	}

	for i := 0; i < len(ids); i++ {
		var backupAgentID int32
		if i == 0 {
			backupAgentID = ids[len(ids)-1]
		} else {
			backupAgentID = ids[i-1]
		}
		err := updateCollectorBackupAgent(client, ids[i], backupAgentID)
		if err != nil {
			log.Warnf("Failed to update the backup collector id: %v", err)
		}
	}

	return nil
}

// DeleteCollectorSet deletes the collectorset.
func DeleteCollectorSet(collectorset *crv1alpha1.CollectorSet, client clientset.Interface) error {
	data := []byte(`[{"op":"add","path":"/spec/replicas","value": 0}]`)
	if _, err := client.AppsV1beta1().StatefulSets(collectorset.Namespace).Patch(collectorset.Name, types.JSONPatchType, data); err != nil {
		return err
	}

	deleteOpts := metav1.DeleteOptions{}
	return client.AppsV1beta1().StatefulSets(collectorset.Namespace).Delete(collectorset.Name, &deleteOpts)
}

func checkCollectorGroupExistsByID(client *client.LMSdkGo, id int32) bool {
	params := lm.NewGetCollectorGroupByIDParams()
	params.SetID(id)
	fields := "id"
	params.SetFields(&fields)
	restResponse, err := client.LM.GetCollectorGroupByID(params)
	if err != nil || restResponse.Payload == nil {
		log.Warnf("Failed to get collector group with id %d", id)
		return false
	}
	return true
}

func getCollectorGroupID(client *client.LMSdkGo, name string, collectorset *crv1alpha1.CollectorSet) (int32, error) {
	params := lm.NewGetCollectorGroupListParams()
	filter := fmt.Sprintf("name:\"%s\"", name)
	params.SetFilter(&filter)
	restResponse, err := client.LM.GetCollectorGroupList(params)
	if err != nil {
		return -1, err
	}

	if restResponse.Payload == nil || restResponse.Payload.Total == 0 {
		log.Infof("Adding collector group with name %q", name)
		return addCollectorGroup(client, name, collectorset)
	}
	if restResponse.Payload.Total == 1 {
		return restResponse.Payload.Items[0].ID, err
	}
	return -1, fmt.Errorf("failed to get collector group ID")
}

func addCollectorGroup(client *client.LMSdkGo, name string, collectorset *crv1alpha1.CollectorSet) (int32, error) {

	kubernetesLabelApp := constants.CustomPropertyKubernetesLabelApp
	kubernetesLabelAppValue := constants.CustomPropertyKubernetesLabelAppValue
	autoClusterName := constants.CustomPropertyAutoClusterName
	AutoClusterNameValue := collectorset.Spec.ClusterName
	customProperties := []*models.NameAndValue{
		{Name: &kubernetesLabelApp, Value: &kubernetesLabelAppValue},
		{Name: &autoClusterName, Value: &AutoClusterNameValue},
	}

	body := &models.CollectorGroup{
		Name:             &name,
		CustomProperties: customProperties,
	}
	params := lm.NewAddCollectorGroupParams()
	params.SetBody(body)
	restResponse, err := client.LM.AddCollectorGroup(params)
	if err != nil {
		return -1, err
	}
	return restResponse.Payload.ID, nil
}

// $(statefulset name)-$(ordinal)
func getCollectorIDs(client *client.LMSdkGo, groupID int32, collectorset *crv1alpha1.CollectorSet) ([]int32, error) {
	var ids []int32
	for ordinal := int32(0); ordinal < *collectorset.Spec.Replicas; ordinal++ {
		name := fmt.Sprintf("%s%s-%d", constants.ClusterCollectorGroupPrefix, collectorset.Spec.ClusterName, ordinal)
		filter := fmt.Sprintf("collectorGroupId:%v,description:\"%v\"", groupID, name)
		params := lm.NewGetCollectorListParams()
		params.SetFilter(&filter)
		restResponse, err := client.LM.GetCollectorList(params)
		if err != nil {
			return nil, err
		}
		var id int32
		if restResponse.Payload == nil || restResponse.Payload.Total == 0 {
			log.Infof("Adding collector with description %q", name)
			kubernetesLabelApp := constants.CustomPropertyKubernetesLabelApp
			kubernetesLabelAppValue := constants.CustomPropertyKubernetesLabelAppValue
			autoClusterName := constants.CustomPropertyAutoClusterName
			AutoClusterNameValue := collectorset.Spec.ClusterName
			customProperties := []*models.NameAndValue{
				{Name: &kubernetesLabelApp, Value: &kubernetesLabelAppValue},
				{Name: &autoClusterName, Value: &AutoClusterNameValue},
			}

			body := &models.Collector{
				Description:                   name,
				CollectorGroupID:              groupID,
				NeedAutoCreateCollectorDevice: false,
				CustomProperties:              customProperties,
			}
			id, err = addCollector(client, body)
			if err != nil {
				return nil, err
			}

			// update the escalating chain id, if failed the value will be the default value
			// the default value of this option param is 0, which means disable notification
			collector, err := getCollectorByID(client, id)
			if err != nil || collector == nil {
				log.Warnf("Failed to get the collector, err: %v", err)
				collector = body
				collector.ID = id
			}

			collector.EscalatingChainID = collectorset.Spec.EscalationChainID
			err = updateCollector(client, collector)
			if err != nil {
				log.Warnf("Failed to update the escalation chain id. The default value will be used. %v", err)
			}
		} else {
			id = restResponse.Payload.Items[0].ID
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func getResourceRequirements(size string) apiv1.ResourceRequirements {
	resourceList := apiv1.ResourceList{}
	var quantity *resource.Quantity
	switch size {
	case "nano":
		quantity = resource.NewQuantity(2*1024*1024*1024, resource.BinarySI)
	case "small":
		quantity = resource.NewQuantity(2*1024*1024*1024, resource.BinarySI)
	case "medium":
		quantity = resource.NewQuantity(4*1024*1024*1024, resource.BinarySI)
	case "large":
		quantity = resource.NewQuantity(8*1024*1024*1024, resource.BinarySI)
	default:
		break
	}
	resourceList[apiv1.ResourceMemory] = *quantity

	return apiv1.ResourceRequirements{
		Limits: resourceList,
	}
}

func addCollector(client *client.LMSdkGo, body *models.Collector) (int32, error) {
	params := lm.NewAddCollectorParams()
	params.SetBody(body)
	restResponse, err := client.LM.AddCollector(params)
	if err != nil {
		return -1, err
	}
	return restResponse.Payload.ID, nil
}

func getCollectorByID(client *client.LMSdkGo, id int32) (*models.Collector, error) {
	params := lm.NewGetCollectorByIDParams()
	params.SetID(id)
	restResponse, err := client.LM.GetCollectorByID(params)
	if err != nil {
		return nil, err
	}
	return restResponse.Payload, nil
}

func updateCollector(client *client.LMSdkGo, body *models.Collector) error {
	params := lm.NewUpdateCollectorByIDParams()
	params.SetBody(body)
	params.SetID(body.ID)
	_, err := client.LM.UpdateCollectorByID(params)
	if err != nil {
		return err
	}

	return nil
}

func updateCollectorBackupAgent(client *client.LMSdkGo, id, backupID int32) error {
	// Get all the fields before updating to prevent setting default values to the other fields
	restResponse, err := getCollectorByID(client, id)
	if err != nil || restResponse == nil {
		return fmt.Errorf("failed to get the collector: %v", err)
	}

	collector := restResponse
	collector.EnableFailBack = true
	collector.BackupAgentID = backupID
	updateErr := updateCollector(client, collector)
	if updateErr != nil {
		return fmt.Errorf("failed to update the collector: %v", updateErr)
	}
	return nil
}
