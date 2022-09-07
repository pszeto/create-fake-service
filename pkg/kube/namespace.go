package kube

import (
	"context"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func EnsureNamespaceExists(clientset *kubernetes.Clientset, ns string, dryRun bool) error {
	if !dryRun {
		_, err := clientset.CoreV1().Namespaces().Get(context.Background(), ns, metav1.GetOptions{})
		if err != nil {
			log.Infof("Creating %s namespace.", ns)
			nsName := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: ns,
				},
			}
			_, err := clientset.CoreV1().Namespaces().Create(context.Background(), nsName, metav1.CreateOptions{})
			if err != nil {
				log.Errorf("Failed to create %s namespace.", ns)
				return err
			} else {
				log.Infof("Successfully created %s namespace.", ns)
			}
		} else {
			log.Infof("Namespace %s already exist.", ns)
		}
	}
	return nil
}
