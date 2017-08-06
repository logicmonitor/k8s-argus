package types

import (
	"github.com/logicmonitor/k8s-argus/pkg/config"
	lm "github.com/logicmonitor/lm-sdk-go"
	"k8s.io/client-go/kubernetes"
)

// Base is a struct for embedding
type Base struct {
	LMClient  *lm.DefaultApi
	K8sClient *kubernetes.Clientset
	Config    *config.Config
}
