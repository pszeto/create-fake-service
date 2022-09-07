# create-fake-service

    Usage:
    create-fake-service [flags]

    Flags:
        --deployment string           Specify name of deployment. Defaults to temporary
        --deployment-replicas int32   Specify the number of replicas for the deployment (default 1)
        --dry-run                     Specify if it's a dry run. Default false
    -h, --help                        help for create-fake-service
        --include-hey                 Specify whether to include hey container in deployment. Default to false
        --kube-context string         Specify which kube context to use.
        --namespace string            Specify namespace. default namespace is used if not set
        --namespace-labels string     Specify the labels for the namespaces, defaults to none. Comma seperated.  Example: istio-injection=enabled,istio.io/rev=1-12        
        --ports string                Specify ports to expose for the service/deployment. Defaults to 8080
        --protocol string             Specify protocol for for the service/deployment. Defaults to http
        --save-yaml                   Specify if it should save the yamls to file. Default false
        --upstream-uris string        Specify the upstream service addresses for the fake service. Comma seperated.  Example: http://some-app.default:8080

Example
`go run main.go --namespace my-sleep --deployment project --ports 8080,443 --kube-context pszeto-cluster1-eks --upstream-uris http://httpbin.default:8000 --protocol http --include-hey true --dry-run true --save-yaml true`