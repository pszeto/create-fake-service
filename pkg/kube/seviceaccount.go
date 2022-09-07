package kube

import (
	"context"
	"fmt"

	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func EnsureServiceAccountExists(clientset *kubernetes.Clientset, ns string, serviceAccount string, dryRun bool, saveYaml bool) error {
	var serviceAccountName = serviceAccount + "-serviceaccount"
	saInfo := ServiceAccount{
		ApiVersion: "v1",
		Kind:       "ServiceAccount",
		Metadata: MetaData{
			Name:      serviceAccountName,
			Namespace: ns,
		},
	}

	if dryRun {
		yamlData, err := yaml.Marshal(&saInfo)
		if err != nil {
			log.Errorf("Error while Marshaling. %v", err)
			return err
		}

		fmt.Println("----------ServiceAccount Yaml----------")
		fmt.Println(string(yamlData))
	} else {
		_, err := clientset.CoreV1().ServiceAccounts(ns).Get(context.Background(), serviceAccountName, metav1.GetOptions{})
		if err != nil {
			log.Infof("Creating %s serviceaccount in %s namespace.", serviceAccountName, ns)
			sa := &corev1.ServiceAccount{
				ObjectMeta: metav1.ObjectMeta{
					Name:      serviceAccountName,
					Namespace: ns,
				},
			}
			_, err := clientset.CoreV1().ServiceAccounts(ns).Create(context.Background(), sa, metav1.CreateOptions{})
			if err != nil {
				log.Errorf("Failed to create %s serviceaccount in %s namespace.", serviceAccountName, ns)
				return err
			} else {
				log.Infof("Successfully created %s serviceaccount in %s namespace.", serviceAccountName, ns)
			}
		} else {
			log.Infof("ServiceAccount %s already exist in %s namespace.", serviceAccountName, ns)
		}
	}

	if saveYaml {
		file, err := os.OpenFile(serviceAccountName+"-"+ns+"-serviceaccount.yaml", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			log.Fatalf("error opening/creating file: %v", err)
			return err
		}
		defer file.Close()

		enc := yaml.NewEncoder(file)

		err = enc.Encode(saInfo)
		if err != nil {
			log.Fatalf("error encoding: %v", err)
			return err
		}
	}

	return nil
}
