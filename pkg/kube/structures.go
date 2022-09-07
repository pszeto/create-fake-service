package kube

type ContainerInfo struct {
	Name            string          `yaml:"name"`
	Image           string          `yaml:"image"`
	ImagePullPolicy string          `yaml:"imagePullPolicy"`
	Ports           []ContainerPort `yaml:"ports,omitempty"`
	Env             []EnvVar        `yaml:"env,omitempty"`
}

type ContainerPort struct {
	Name          string `yaml:"name"`
	Protocol      string `yaml:"protocol"`
	Port          int32  `yaml:"port,omitempty"`
	ContainerPort int32  `yaml:"containerPort,omitempty"`
	TargetPort    int32  `yaml:"targetPort,omitempty"`
}

type ContainerSpec struct {
	ServiceAccountName string          `yaml:"serviceAccountName"`
	Containers         []ContainerInfo `yaml:"containers"`
}

type Deployment struct {
	ApiVersion string         `yaml:"apiVersion"`
	Kind       string         `yaml:"kind"`
	Metadata   MetaData       `yaml:"metadata"`
	Spec       DeploymentSpec `yaml:"spec"`
}

type DeploymentSpec struct {
	Selector MatchLabelSelector `yaml:"selector"`
	Template DeploymentTemplate `yaml:"template"`
}

type DeploymentTemplate struct {
	ObjectMetadata ObjectMetaDataLabels `yaml:"metadata"`
	Spec           ContainerSpec        `yaml:"spec"`
}

type EnvVar struct {
	Name      string   `yaml:"name"`
	Value     string   `yaml:"value,omitempty"`
	ValueFrom FieldRef `yaml:"valueFrom,omitempty"`
}

type FieldPath struct {
	FieldPath string `yaml:"fieldPath"`
}

type FieldRef struct {
	FieldRef FieldPath `yaml:"fieldRef"`
}

type MatchLabels struct {
	MatchLabels map[string]string `yaml:"matchLabels"`
}

type MatchLabelSelector struct {
	MatchLabels map[string]string `yaml:"matchLabels"`
}

type MetaData struct {
	Name      string            `yaml:"name"`
	Namespace string            `yaml:"namespace"`
	Labels    map[string]string `yaml:"labels,omitempty"`
}

type ObjectMetaDataLabels struct {
	Labels map[string]string `yaml:"labels"`
}

type Service struct {
	ApiVersion string      `yaml:"apiVersion"`
	Kind       string      `yaml:"kind"`
	Metadata   MetaData    `yaml:"metadata"`
	Spec       ServiceSpec `yaml:"spec"`
}

type ServiceAccount struct {
	ApiVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   MetaData `yaml:"metadata"`
}

type ServiceSpec struct {
	Selector map[string]string `yaml:"selector"`
	Ports    []ContainerPort   `yaml:"ports"`
}
