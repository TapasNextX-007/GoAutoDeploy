package main

import (
	"flag"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// Accepting domain name from the user
	var domain string
	flag.StringVar(&domain, "domain", "", "The domain to be used for deployment")
	flag.Parse()

	if domain == "" {
		log.Fatalf("Please provide a domain using the -domain flag")
		return
	}

	values := map[string]interface{}{
		"ingress": map[string]interface{}{
			"hosts": []map[string]interface{}{
				{
					"host": domain,
					"paths": []map[string]interface{}{
						{
							"path":     "/",
							"pathType": "ImplementationSpecific",
						},
					},
				},
			},
		},
	}

	// Providing helm chart details
	chartName := "my-nginx-chart"
	releaseName := "my-nginx-release"
	_, err := initializeKubeClient()
	if err != nil {
		log.Fatalf("error initializing kube client: %v", err)
		return
	}
	// Instantiating all the env for helm
	settings := cli.New()
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		log.Printf("%+v", err)
		os.Exit(1)
	}

	if err != nil {
		log.Fatalf("error downloading index file: %v", err)
		return
	}
	// Leverage install action performed by helm
	installAction := action.NewInstall(actionConfig)
	installAction.ReleaseName = releaseName

	// Locating the chart details from reo
	cp, err := installAction.ChartPathOptions.LocateChart(chartName, settings)
	if err != nil {
		log.Fatalf("Unable to locate the chart %v", err)
		return
	}

	//Loading the chart
	chartReq, err := loader.Load(cp)
	_, err = installAction.Run(chartReq, values)
	if err != nil {
		log.Fatalf("error installing chart: %v", err)
		return
	}
	log.Println("Chart installed successfully!")
}

// Building the kube configuration
func initializeKubeClient() (*kubernetes.Clientset, error) {
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}
