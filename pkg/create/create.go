package create

import (
	"strconv"
	"strings"

	kube "github.com/pszeto/create-fake-service/pkg/kube"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

func New(config *Config) *App {
	return &App{
		config: config,
	}
}

func (app *App) Entry() error {
	log.SetFormatter(&log.JSONFormatter{})

	var clientset *kubernetes.Clientset
	var err error
	var portsAsInt []int32

	if !app.config.DryRun {
		clientset, err = kube.GetClientSet(app.config.Kubectx)
		if err != nil {
			log.Fatalf("Exiting program. Fatal Error : %v", err)
		}
	}

	if app.config.Namespace == "" {
		log.Infoln("Namespace not specified. Using default namespace default.")
		app.config.Namespace = "default"
	}

	if app.config.Ports == "" {
		log.Infoln("Ports not specified. Using default port 8080.")
		portsAsInt = append(portsAsInt, 8080)
	} else {
		portsArray := strings.Split(app.config.Ports, ",")
		for _, port := range portsArray {
			intVar, err := strconv.ParseInt(port, 10, 32)
			if err != nil {
				log.Errorln("Invalid ports specified.")
				log.Fatalln(err)
				return err
			} else {
				portsAsInt = append(portsAsInt, int32(intVar))
			}
		}
	}

	if app.config.Protocol == "" {
		log.Infoln("protocol not specified. Using default protocol http.")
		app.config.Protocol = "http"
	}

	if app.config.DeploymentReplicas != 1 {
		log.Infoln("Deployment Replicas:", app.config.DeploymentReplicas)
		app.config.Protocol = "http"
	}

	if !app.config.IncludeHey {
		log.Infoln("include-hey not specified. Defaulting to false.")
	}

	if !app.config.SaveYaml {
		log.Infoln("save-yaml not specified. Defaulting to false.")
	}

	if app.config.UpstreamUris == "" {
		log.Infoln("Upstream Service URIs not specified.  Defaulting to http://httpbin.default:8080.")
		app.config.UpstreamUris = "http://httpbin.default:8080"
	}

	if app.config.Deployment == "" {
		log.Infoln("Deployment name not specified. Using default deployment name temporary.")
		app.config.Deployment = "temporary"
	}

	err = kube.EnsureNamespaceExists(clientset, app.config.Namespace, app.config.DryRun)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	err = kube.EnsureServiceAccountExists(clientset, app.config.Namespace, app.config.Deployment, app.config.DryRun, app.config.SaveYaml)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	err = kube.CreateUpdateDeployment(clientset, app.config.Namespace, app.config.Deployment, app.config.DeploymentReplicas, portsAsInt, app.config.Protocol, app.config.UpstreamUris, app.config.IncludeHey, app.config.DryRun, app.config.SaveYaml)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	err = kube.CreateUpdateService(clientset, app.config.Namespace, app.config.Deployment, portsAsInt, app.config.Protocol, app.config.DryRun, app.config.SaveYaml)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	return nil
}
