package kube

import (
	"context"
	"strings"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func convertNamespaceLabels(labels string) map[string]string {
	labelMap := make(map[string]string)
	if labels != "" {
		labelsSplice := strings.Split(labels, ",")
		for _, label := range labelsSplice {
			temp := strings.Split(label, "=")
			labelMap[temp[0]] = temp[1]
		}
	}
	return labelMap
}

func EnsureNamespaceExists(clientset *kubernetes.Clientset, namespace string, nsLabels string, dryRun bool) error {
	if !dryRun {
		_, err := clientset.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
		ns := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name:   namespace,
				Labels: convertNamespaceLabels(nsLabels),
			},
		}
		if err != nil {
			log.Infof("Creating %s namespace.", namespace)
			_, err := clientset.CoreV1().Namespaces().Create(context.Background(), ns, metav1.CreateOptions{})
			if err != nil {
				log.Errorf("Failed to create %s namespace.", namespace)
				return err
			} else {
				log.Infof("Successfully created %s namespace.", namespace)
			}
		} else {
			log.Infof("Namespace %s already exist.", namespace)
			log.Infof("Updating %s namespace.", namespace)
			_, err := clientset.CoreV1().Namespaces().Update(context.Background(), ns, metav1.UpdateOptions{})
			if err != nil {
				log.Errorf("Failed to update %s namespace.", namespace)
				return err
			} else {
				log.Infof("Successfully updated %s namespace.", namespace)
			}
		}
	}
	return nil
}
