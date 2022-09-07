package create

type App struct {
	config *Config
}

type Config struct {
	Namespace          string
	Deployment         string
	DeploymentReplicas int32
	Kubectx            string
	Ports              string
	UpstreamUris       string
	Protocol           string
	IncludeHey         bool
	DryRun             bool
	SaveYaml           bool
}
