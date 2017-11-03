package controller

import (
	"fmt"
	"strings"

	crv1alpha1 "github.com/logicmonitor/k8s-collectorset-controller/pkg/apis/v1alpha1"
	"github.com/logicmonitor/k8s-collectorset-controller/pkg/constants"
	"github.com/logicmonitor/k8s-collectorset-controller/pkg/utilities"
	lm "github.com/logicmonitor/lm-sdk-go"
	log "github.com/sirupsen/logrus"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	clientset "k8s.io/client-go/kubernetes"
)

// CreateOrUpdateCollectorSet creates a replicaset for each collector in
// a CollectorSet
func CreateOrUpdateCollectorSet(collectorset *crv1alpha1.CollectorSet, lmClient *lm.DefaultApi, client clientset.Interface) ([]int32, error) {
	groupID, err := getCollectorGroupID(lmClient, collectorset.Name)
	if err != nil {
		return nil, err
	}
	log.Printf("Collector group %q has ID %d", strings.Title(collectorset.Name), groupID)

	ids, err := getCollectorIDs(lmClient, groupID, collectorset)
	if err != nil {
		return nil, err
	}

	secretIsOptional := false

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
							Image:           "logicmonitor/k8s-collector:develop",
							ImagePullPolicy: apiv1.PullAlways,
							Env: []apiv1.EnvVar{
								{
									Name: "ACCOUNT",
									ValueFrom: &apiv1.EnvVarSource{
										SecretKeyRef: &apiv1.SecretKeySelector{
											LocalObjectReference: apiv1.LocalObjectReference{
												Name: constants.ArgusSecretName,
											},
											Key:      "account",
											Optional: &secretIsOptional,
										},
									},
								},
								{
									Name: "ACCESS_ID",
									ValueFrom: &apiv1.EnvVarSource{
										SecretKeyRef: &apiv1.SecretKeySelector{
											LocalObjectReference: apiv1.LocalObjectReference{
												Name: constants.ArgusSecretName,
											},
											Key:      "accessID",
											Optional: &secretIsOptional,
										},
									},
								},
								{
									Name: "ACCESS_KEY",
									ValueFrom: &apiv1.EnvVarSource{
										SecretKeyRef: &apiv1.SecretKeySelector{
											LocalObjectReference: apiv1.LocalObjectReference{
												Name: constants.ArgusSecretName,
											},
											Key:      "accessKey",
											Optional: &secretIsOptional,
										},
									},
								},
								{
									Name:  "COLLECTOR_SIZE",
									Value: collectorset.Spec.Size,
								},
								{
									Name:  "COLLECTOR_IDS",
									Value: strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ids)), ","), "[]"),
								},
							},
							Resources: getResourceRequirements(collectorset.Spec.Size),
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

	if _, _err := client.AppsV1beta1().StatefulSets(statefulset.ObjectMeta.Namespace).Create(&statefulset); _err != nil {
		if !apierrors.IsAlreadyExists(_err) {
			return nil, _err
		}
		if _, _err := client.AppsV1beta1().StatefulSets(statefulset.ObjectMeta.Namespace).Update(&statefulset); _err != nil {
			return nil, _err
		}
	}

	collectorset.Status.IDs = ids

	err = updateCollectors(lmClient, collectorset, ids, groupID)
	if err != nil {
		log.Warnf("Failed to set collector backup agents: %v", err)
	}
	return collectorset.Status.IDs, nil
}

func updateCollectors(client *lm.DefaultApi, collectorset *crv1alpha1.CollectorSet, ids []int32, groupID int32) error {
	// force an even number of elements in the slice
	even := len(ids)%2 != 0
	if even {
		ids = ids[:len(ids)-1]
	}

	for i := 0; i < len(ids); i++ {
		name := fmt.Sprintf("%s-%d", collectorset.Name, i)
		// We are at the beginning of a pair
		if i%2 == 0 {
			err := updateCollectorBackupAgent(client, groupID, ids[i], ids[i+1], name)
			if err != nil {
				return err
			}
		} else {
			err := updateCollectorBackupAgent(client, groupID, ids[i], ids[i-1], name)
			if err != nil {
				return err
			}
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

func getCollectorGroupID(client *lm.DefaultApi, name string) (int32, error) {
	name = strings.Title(name)
	restResponse, apiResponse, err := client.GetCollectorGroupList("", 1, 0, "name:"+name)
	if _err := utilities.CheckAllErrors(restResponse, apiResponse, err); _err != nil {
		return -1, _err
	}
	if restResponse.Data.Total == 0 {
		log.Printf("Adding collector group with name %q", name)
		return addCollectorGroup(client, name)
	}
	if restResponse.Data.Total == 1 {
		return restResponse.Data.Items[0].Id, err
	}
	return -1, fmt.Errorf("Failed to get collector group ID")
}

func addCollectorGroup(client *lm.DefaultApi, name string) (int32, error) {
	group := lm.RestCollectorGroup{
		Name: name,
	}
	restResponse, apiResponse, err := client.AddCollectorGroup(group)
	if _err := utilities.CheckAllErrors(restResponse, apiResponse, err); _err != nil {
		return -1, _err
	}
	return restResponse.Data.Id, nil
}

// $(statefulset name)-$(ordinal)
func getCollectorIDs(client *lm.DefaultApi, groupID int32, collectorset *crv1alpha1.CollectorSet) ([]int32, error) {
	var ids []int32
	for ordinal := int32(0); ordinal < *collectorset.Spec.Replicas; ordinal++ {
		name := fmt.Sprintf("%s-%d", collectorset.Name, ordinal)
		restResponse, apiResponse, err := client.GetCollectorList("", 1, 0, "description:"+name)
		if _err := utilities.CheckAllErrors(restResponse, apiResponse, err); _err != nil {
			return nil, _err
		}
		var id int32
		if restResponse.Data.Total == 0 {
			log.Printf("Adding collector with description %q", name)
			id, err = addCollector(client, groupID, name)
			if err != nil {
				return nil, err
			}
		} else {
			id = restResponse.Data.Items[0].Id
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
		break
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

func addCollector(client *lm.DefaultApi, groupID int32, description string) (int32, error) {
	collector := lm.RestCollector{
		Description:                   description,
		CollectorGroupId:              groupID,
		NeedAutoCreateCollectorDevice: false,
	}
	restResponse, apiResponse, err := client.AddCollector(collector)
	if _err := utilities.CheckAllErrors(restResponse, apiResponse, err); _err != nil {
		return -1, _err
	}
	return restResponse.Data.Id, nil
}

func updateCollectorBackupAgent(client *lm.DefaultApi, groupID, id, backupID int32, description string) error {
	collector := lm.RestCollector{
		Description:                   description,
		CollectorGroupId:              groupID,
		NeedAutoCreateCollectorDevice: false,
		EnableFailBack:                true,
		BackupAgentId:                 backupID,
	}
	_, _, err := client.UpdateCollectorById(id, collector)

	return err
}
