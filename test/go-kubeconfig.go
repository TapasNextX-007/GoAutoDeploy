// Checks cluster health or retrieves cluster nodes to test the connectivity

package main

import (
	"context"
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"path/filepath"
)

func main() {
	clientset, err := InitializeKubeClient()
	if err != nil {
		log.Fatalf("error initializing kube client: %v", err)
		return
	}
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		log.Fatalf("error listing nodes: %v", err)
		return
	}

	fmt.Printf("Cluster nodes:\n\t%v\n", nodes)
	for _, node := range nodes.Items {
		fmt.Printf("Name: %s, Version: %s\n", node.Name, node.Status.NodeInfo.KubeletVersion)
	}
}

// Initializing and building Kubernetes client configuration in Go
func InitializeKubeClient() (*kubernetes.Clientset, error) {
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
