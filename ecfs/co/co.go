package co

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// clientSet returns Container orchestration configuration
func clientSet() (clientSet *kubernetes.Clientset) {
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
