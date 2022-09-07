package kube

import (
	"context"
	"fmt"

	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CreateUpdateDeployment(clientset *kubernetes.Clientset, ns string, deploymentName string, numOfReplicas int32, ports []int32, protocol string, upstreamServiceAddress string, includeHey bool, dryRun bool, saveYaml bool) error {
	var serviceAccountName = deploymentName + "-serviceaccount"
	containers := []corev1.Container{}
	containersInfo := []ContainerInfo{}

	for index, port := range ports {
		ports := []corev1.ContainerPort{
			{
				Name:          protocol,
				Protocol:      corev1.ProtocolTCP,
				ContainerPort: port,
			},
		}
		envVar := []corev1.EnvVar{
			{
				Name:  "LISTEN_ADDR",
				Value: "0.0.0.0:" + fmt.Sprint(port),
			},
			{
				Name:  "MESSAGE",
				Value: "I AM ALIVE",
			},
			{
				Name:  "UPSTREAM_URIS",
				Value: upstreamServiceAddress,
			},
			{
				Name: "NAME",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: "metadata.name",
					},
				},
			},
		}
		containers = append(containers, corev1.Container{
			Name:            deploymentName + "-" + fmt.Sprint(index),
			Image:           "nicholasjackson/fake-service:v0.7.8",
			ImagePullPolicy: "Always",
			Ports:           ports,
			Env:             envVar,
		})

		containersInfo = append(containersInfo, ContainerInfo{
			Name:            deploymentName + "-" + fmt.Sprint(index),
			Image:           "nicholasjackson/fake-service:v0.7.8",
			ImagePullPolicy: "Always",
			Env: []EnvVar{
				{
					Name:  "LISTEN_ADDR",
					Value: "0.0.0.0:" + fmt.Sprint(port),
				},
				{
					Name:  "MESSAGE",
					Value: "I AM ALIVE",
				},
				{
					Name:  "UPSTREAM_URIS",
					Value: upstreamServiceAddress,
				},
				{
					Name: "NAME",
					ValueFrom: FieldRef{
						FieldRef: FieldPath{
							FieldPath: "metadata.name",
						},
					},
				},
			},
			Ports: []ContainerPort{
				{
					Name:          protocol,
					Protocol:      "TCP",
					ContainerPort: port,
				},
			},
		})
	}

	if includeHey {
		if !dryRun {
			log.Infof("Including hey container in %s deployment.", deploymentName)
		}
		containers = append(containers, corev1.Container{
			Name:            "hey",
			Image:           "docker.io/pszeto/hey",
			ImagePullPolicy: "Always",
		})
		containersInfo = append(containersInfo, ContainerInfo{
			Name:            "hey",
			Image:           "docker.io/pszeto/hey",
			ImagePullPolicy: "Always",
		})
	}

	deploymentSpec := appsv1.DeploymentSpec{
		Replicas: &numOfReplicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app": deploymentName,
			},
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"app": deploymentName,
				},
			},
			Spec: corev1.PodSpec{
				ServiceAccountName: serviceAccountName,
				Containers:         containers,
			},
		},
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Namespace: ns,
			Labels: map[string]string{
				"app": deploymentName,
			},
		},
		Spec: deploymentSpec,
	}

	deploymentInfo := &Deployment{
		ApiVersion: "apps/v1",
		Kind:       "Deployment",
		Metadata: MetaData{
			Name:      deploymentName,
			Namespace: ns,
			Labels: map[string]string{
				"app": deploymentName,
			},
		},
		Spec: DeploymentSpec{
			Replicas: numOfReplicas,
			Selector: MatchLabelSelector{
				MatchLabels: map[string]string{
					"app": deploymentName,
				},
			},
			Template: DeploymentTemplate{
				ObjectMetadata: ObjectMetaDataLabels{
					Labels: map[string]string{
						"app": deploymentName,
					},
				},
				Spec: ContainerSpec{
					ServiceAccountName: serviceAccountName,
					Containers:         containersInfo,
				},
			},
		},
	}

	if dryRun {
		yamlData, err := yaml.Marshal(&deploymentInfo)
		if err != nil {
			log.Errorf("Error while Marshaling. %v", err)
			return err
		}

		fmt.Println("---------- Deployment Yaml----------")
		fmt.Println(string(yamlData))
	} else {
		_, err := clientset.AppsV1().Deployments(ns).Get(context.Background(), deploymentName, metav1.GetOptions{})
		if err == nil {
			log.Infof("Deployment %s already exist in %s namespace.", deploymentName, ns)
			log.Infof("Deleting %s deployment in %s namespace.", deploymentName, ns)
			err = clientset.AppsV1().Deployments(ns).Delete(context.Background(), deploymentName, *metav1.NewDeleteOptions(0))
			if err != nil {
				log.Errorf("Failed to delete %s deployment in %s namespace.", deploymentName, ns)
				return err
			}
		}
		log.Infof("Creating %s deployment in %s namespace.", deploymentName, ns)
		_, err = clientset.AppsV1().Deployments(ns).Create(context.TODO(), deployment, metav1.CreateOptions{})
		if err != nil {
			log.Errorf("Failed to create %s deployment in %s namespace.", deploymentName, ns)
			return err
		} else {
			log.Infof("Successfully created %s deployment in %s namespace.", deploymentName, ns)
		}
	}

	if saveYaml {
		file, err := os.OpenFile(deploymentName+"-"+ns+"-deployment.yaml", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			log.Fatalf("error opening/creating file: %v", err)
			return err
		}
		defer file.Close()

		enc := yaml.NewEncoder(file)

		err = enc.Encode(deploymentInfo)
		if err != nil {
			log.Fatalf("error encoding: %v", err)
			return err
		}
	}

	return nil
}
