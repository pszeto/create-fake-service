package kube

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func buildConfigFromFlags(context, kubeconfigPath string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
}

func GetClientSet(kubectx string) (*kubernetes.Clientset, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		log.Errorf("Error getting user home dir: %v", err)
		return nil, err
	}
	kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")
	log.Infof("Using kubeconfig: %v", kubeConfigPath)
	var kubeConfig *rest.Config
	if kubectx != "" {
		log.Infof("Using context: %v", kubectx)
		kubeConfig, err = buildConfigFromFlags(kubectx, kubeConfigPath)
	} else {
		kubeConfig, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	}

	if err != nil {
		log.Errorf("Error getting Kubernetes config: %v", err)
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		log.Errorf("Error getting Kubernetes clientset: %v", err)
		return nil, err
	}
	return clientset, nil
}
