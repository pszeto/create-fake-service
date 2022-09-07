package cmd

import (
	"os"

	"github.com/pszeto/create-fake-service/pkg/create"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cfg = &create.Config{}

var rootCmd = &cobra.Command{
	Use:   "create-fake-service",
	Short: "configures and deploys fake-service on a cluster",
	Long:  "configures and deploys fake-service on a cluster",
	Run: func(cmd *cobra.Command, args []string) {
		create.New(cfg).Entry()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&cfg.Namespace, "namespace", "", "Specify namespace. default namespace is used if not set")
	rootCmd.Flags().StringVar(&cfg.NamespaceLabels, "namespace-labels", "", "Specify the labels for the namespaces, defaults to none. Comma seperated.  Example: istio-injection=enabled,istio.io/rev=1-12")
	rootCmd.Flags().StringVar(&cfg.Deployment, "deployment", "", "Specify name of deployment. Defaults to temporary")
	rootCmd.Flags().StringVar(&cfg.Ports, "ports", "", "Specify ports to expose for the service/deployment. Defaults to 8080")
	rootCmd.Flags().StringVar(&cfg.Protocol, "protocol", "", "Specify protocol for for the service/deployment. Defaults to http")
	rootCmd.Flags().StringVar(&cfg.UpstreamUris, "upstream-uris", "", "Specify the upstream service addresses for the fake service. Comma seperated.  Example: http://some-app.default:8080")
	rootCmd.Flags().BoolVar(&cfg.IncludeHey, "include-hey", false, "Specify whether to include hey container in deployment. Default to false")
	rootCmd.Flags().Int32Var(&cfg.DeploymentReplicas, "deployment-replicas", 1, "Specify the number of replicas for the deployment")
	rootCmd.Flags().StringVar(&cfg.Kubectx, "kube-context", "", "Specify which kube context to use.")
	rootCmd.Flags().BoolVar(&cfg.DryRun, "dry-run", false, "Specify if it's a dry run. Default false")
	rootCmd.Flags().BoolVar(&cfg.SaveYaml, "save-yaml", false, "Specify if it should save the yamls to file. Default false")
}
