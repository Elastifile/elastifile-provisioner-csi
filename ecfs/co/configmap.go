package co

import (
	"github.com/elastifile/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetConfigMap(namespace string, configMapName string) (data map[string]string, err error) {
	clientSet := clientSet()
	configMap, err := clientSet.CoreV1().ConfigMaps(namespace).Get(configMapName, metav1.GetOptions{})
	if err != nil {
		err = errors.Wrap(err, 0)
		return
	}
	data = configMap.Data
	return
}

func GetConfigMapValue(namespace string, configMapName string, key string) (value string, err error) {
	configMap, err := GetConfigMap(namespace, configMapName)
	if err != nil {
		err = errors.Wrap(err, 0)
		return
	}

	data, ok := configMap[key]
	if !ok {
		err = errors.Errorf("Key %v not found in config map %v (namespace: %v)", key, configMapName, namespace)
	}
	value = string(data)
	return
}
