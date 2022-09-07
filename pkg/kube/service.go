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

func CreateUpdateService(clientset *kubernetes.Clientset, ns string, serviceName string, ports []int32, protocol string, dryRun bool, saveYaml bool) error {
	servicePorts := []corev1.ServicePort{}
	servicePortsInfo := []ContainerPort{}

	for index, port := range ports {
		var portName string
		if port == 443 {
			portName = "https"
		} else {
			portName = protocol + "-" + fmt.Sprint(index)
		}
		servicePorts = append(servicePorts, corev1.ServicePort{
			Name:     portName,
			Port:     port,
			Protocol: "TCP",
		})
		servicePortsInfo = append(servicePortsInfo, ContainerPort{
			Name:     portName,
			Port:     port,
			Protocol: "TCP",
		})
	}

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: ns,
			Labels: map[string]string{
				"app": serviceName,
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": serviceName,
			},
			Ports: servicePorts,
		},
	}

	serviceInfo := Service{
		ApiVersion: "v1",
		Kind:       "Service",
		Metadata: MetaData{
			Name:      serviceName,
			Namespace: ns,
			Labels: map[string]string{
				"app": serviceName,
			},
		},
		Spec: ServiceSpec{
			Selector: map[string]string{
				"app": serviceName,
			},
			Ports: servicePortsInfo,
		},
	}
	if dryRun {
		yamlData, err := yaml.Marshal(&serviceInfo)
		if err != nil {
			log.Errorf("Error while Marshaling. %v", err)
			return err
		}

		fmt.Println("----------Service Yaml----------")
		fmt.Println(string(yamlData))
	} else {
		_, err := clientset.CoreV1().Services(ns).Get(context.Background(), serviceName, metav1.GetOptions{})
		if err == nil {
			log.Infof("Service %s already exist in %s namespace.", serviceName, ns)
			log.Infof("Deleting %s service in %s namespace.", serviceName, ns)
			err = clientset.CoreV1().Services(ns).Delete(context.Background(), serviceName, *metav1.NewDeleteOptions(0))
			if err != nil {
				log.Errorf("Failed to delete %s service in %s namespace.", serviceName, ns)
				return err
			}
		}
		log.Infof("Creating %s service in %s namespace.", serviceName, ns)
		_, err = clientset.CoreV1().Services(ns).Create(context.Background(), service, metav1.CreateOptions{})
		if err != nil {
			log.Errorf("Failed to create %s service in %s namespace.", serviceName, ns)
			return err
		} else {
			log.Infof("Successfully created %s service in %s namespace.", serviceName, ns)
		}
	}

	if saveYaml {
		file, err := os.OpenFile(serviceName+"-"+ns+"-service.yaml", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			log.Fatalf("error opening/creating file: %v", err)
			return err
		}
		defer file.Close()

		enc := yaml.NewEncoder(file)

		err = enc.Encode(serviceInfo)
		if err != nil {
			log.Fatalf("error encoding: %v", err)
			return err
		}
	}

	return nil
}
