package co

/*
	Container Orchestrator support
*/

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// getClientSet returns Container Orchestration configuration
func getClientSet() (clientSet *kubernetes.Clientset) {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return
}
