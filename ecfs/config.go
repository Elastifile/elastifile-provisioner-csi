package main

import (
	"github.com/elastifile/errors"
	"github.com/golang/glog"
)

type config struct {
	NFSServer       string `parameter:"nfsServer"`
	EmanageURL      string `parameter:"emanageURL"`
	Username        string `parameter:"username"`
	Password        string
	SecretName      string `parameter:"secretName"`
	SecretNamespace string `parameter:"secretNamespace"`
}

//// InClusterConfig
//func readConfigMap() {
//	config, err := rest.InClusterConfig()
//	if err != nil {
//		panic(err.Error())
//	}
//
//	// creates the clientset
//	clientset, err := kubernetes.NewForConfig(config)
//	if err != nil {
//		panic(err.Error())
//	}
//
//	clientset.CoreV1()
//
//	for {
//		pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
//		if err != nil {
//			panic(err.Error())
//		}
//		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
//
//		// Examples for error handling:
//		// - Use helper functions like e.g. errors.IsNotFound()
//		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
//		_, err = clientset.CoreV1().Pods("default").Get("example-xxxxx", metav1.GetOptions{})
//		if errors.IsNotFound(err) {
//			fmt.Printf("Pod not found\n")
//		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
//			fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
//		} else if err != nil {
//			panic(err.Error())
//		} else {
//			fmt.Printf("Found pod\n")
//		}
//
//		time.Sleep(10 * time.Second)
//	}
//}
//
//func readSecrets(conf *config, namespace string) error {
//	if kubeClient == nil {
//		return nil, fmt.Errorf("Cannot get kube client")
//	}
//	secrets, err := kubeClient.CoreV1().Secrets(secretNs).Get(secretName, metav1.GetOptions{})
//	if err != nil {
//		err = fmt.Errorf("Couldn't get secret %v/%v err: %v", secretNs, secretName, err)
//		return nil, err
//	}
//	for name, data := range secrets.Data {
//		secret = string(data)
//		glog.V(4).Infof("found ceph secret info: %s", name)
//	}
//
//	secret, err := secrets.Get(conf.SecretName)
//	if err != nil {
//		return fmt.Errorf("secret %q: %v", conf.SecretName, err)
//	}
//
//	const passwordKey = "password.txt"
//	password, ok := secret.Data[passwordKey]
//	if !ok {
//		return fmt.Errorf("secret %q: no value found for key %q", conf.SecretName, passwordKey)
//	}
//
//	conf.Password = string(password)
//	return nil
//}

func GetProvisionerSettings() (configMap map[string]string, secret map[string]string, err error) {
	// TODO: Implement fetching cluster configuration and credentials from configmap+secret - these are fake values
	glog.Warning("Config map and secrets are not yet supported - using hard-coded values!!!")

	configMap = make(map[string]string)
	secret = make(map[string]string)

	configMap[managementAddress] = "https://35.241.144.0"
	configMap[managementUserName] = "admin"
	configMap[nfsAddress] = "172.28.0.5"
	secret[managementPassword] = "changeme"

	return
}

func newPluginConfig() (conf *config, err error) {
	configMap, secret, err := GetProvisionerSettings()
	if err != nil {
		err = errors.Wrap(err, 0)
	}

	// TODO: Check key availability
	conf = &config{
		NFSServer:  configMap[nfsAddress],
		EmanageURL: configMap[managementAddress],
		Username:   configMap[managementUserName],
		Password:   secret[managementPassword],
	}

	// TODO: Mask the password in the log
	glog.Infof("Parsed config map and secrets %+v", conf)
	return
}
