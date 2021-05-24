package config

import (
	"fmt"
	"sync"
	"time"

	"github.com/logicmonitor/k8s-argus/pkg/constants"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	k8sClientSet *kubernetes.Clientset
	w            *watcher
)

const defaultConfigResyncDuration = 10 * time.Minute

// GetClientSet get k8sClient object
func GetClientSet() *kubernetes.Clientset {
	return k8sClientSet
}

type watcher struct {
	labelSelector labels.Set
	configData    map[string]string
	mu            *sync.Mutex
}

// Init initialises k8sClient and creates watcher object
func Init(kubeConfigFile string) error {
	var err error
	k8sClientSet, err = newK8sClient(kubeConfigFile)
	if err != nil {
		return fmt.Errorf("failed to initialise kubernetes rest clientset with error: %w", err)
	}
	w = NewConfigWatcher()

	return nil
}

// IConfig config w interface
type IConfig interface {
	Set(key string, value string)
	Get(name string) (string, error)
	Run()
}

// Set save new key value in configmap data cache
func (cw *watcher) Set(key string, value string) {
	cw.mu.Lock()
	defer cw.mu.Unlock()
	cw.configData[key] = value
}

// NewConfigWatcher creates new config watcher
func NewConfigWatcher() *watcher { // nolint: golint,revive
	return &watcher{
		labelSelector: labels.Set{"chart": "argus"},
		configData:    make(map[string]string),
		mu:            new(sync.Mutex),
	}
}

// Load graceful load to on init to load config
func Load() error {
	if w == nil {
		return fmt.Errorf("could not initialise config/watcher")
	}
	cmList, err := k8sClientSet.CoreV1().ConfigMaps(metav1.NamespaceAll).List(metav1.ListOptions{
		LabelSelector: w.labelSelector.String(),
	})
	if err != nil {
		return fmt.Errorf("failed to load application configmap with selector %s: %w", w.labelSelector, err)
	}
	if len(cmList.Items) == 0 {
		return fmt.Errorf("could not find configmap with selector %s", w.labelSelector)
	}
	logrus.Infof("Configmap Data: %v", cmList.Items[0].Data)
	// configmap cannot be shared across namespace so even if multiple presents, namespace of any configmap represents it
	w.Set("namespace", cmList.Items[0].Namespace)
	for k, v := range cmList.Items[0].Data {
		w.Set(k, v)
	}

	return nil
}

// Run start watcher to listen for config map changes
func Run() {
	optionsModifier := func(options *metav1.ListOptions) {
		options.FieldSelector = fields.Everything().String()
		options.LabelSelector = labels.Set{"chart": "argus"}.String()
	}
	watchlist := cache.NewFilteredListWatchFromClient(k8sClientSet.CoreV1().RESTClient(), "configmaps", corev1.NamespaceAll, optionsModifier)
	_, controller := cache.NewInformer(
		watchlist,
		&corev1.ConfigMap{}, // nolint: exhaustivestruct
		defaultConfigResyncDuration,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    AddFunc(),
			DeleteFunc: DeleteFunc(),
			UpdateFunc: UpdateFunc(),
		},
	)

	stop := make(chan struct{})
	go controller.Run(stop)
}

// Get returns config value if present else error
func (cw *watcher) Get(name string) (string, error) {
	cw.mu.Lock()
	defer cw.mu.Unlock()
	val, ok := cw.configData[name]
	if !ok {
		return "", fmt.Errorf("no config present for %s", name)
	}

	return val, nil
}

// AddFunc is a function that implements the w interface.
func AddFunc() func(obj interface{}) {
	return func(obj interface{}) {
		o := obj.(*corev1.ConfigMap) // nolint: forcetypeassert
		o.ManagedFields = make([]metav1.ManagedFieldsEntry, 0)
		logrus.Infof("Add Configmap: %v", o)

		for k, v := range o.Data {
			w.Set(k, v)
		}
	}
}

// UpdateFunc is a function that implements the w interface.
func UpdateFunc() func(oldObj, newObj interface{}) {
	return func(oldObj, newObj interface{}) {
		o := oldObj.(*corev1.ConfigMap) // nolint: forcetypeassert
		n := newObj.(*corev1.ConfigMap) // nolint: forcetypeassert
		o.ManagedFields = make([]metav1.ManagedFieldsEntry, 0)
		n.ManagedFields = make([]metav1.ManagedFieldsEntry, 0)
		logrus.Infof("Update Configmap: %v changed to new %v", o, n)
		for k, v := range n.Data {
			w.Set(k, v)
		}
	}
}

// DeleteFunc is a function that implements the w interface.
func DeleteFunc() func(obj interface{}) {
	return func(obj interface{}) {
		o := obj.(*corev1.ConfigMap) // nolint: forcetypeassert
		// Deliberately delete not implemented, deleting config map could abnormally stop application.
		// user may delete configmap while updating any param - delete and recreate.
		logrus.Warnf("Delete Configmap does not change loaded application as application may stop abnormally: %v", o)
	}
}

func newK8sClient(filePath string) (*kubernetes.Clientset, error) {
	restClientConfig, err := clientcmd.BuildConfigFromFlags("", filePath)
	if err != nil {
		return nil, err
	}

	restClientConfig.UserAgent = constants.UserAgentBase + constants.Version

	clientSet, err := kubernetes.NewForConfig(restClientConfig)
	if err != nil {
		return nil, err
	}

	return clientSet, nil
}

// GetWatchConfig returns config value if present else error
func GetWatchConfig(name string) (string, error) {
	return w.Get(name)
}
