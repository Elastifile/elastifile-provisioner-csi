package co

import (
	"github.com/go-errors/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetConfigMap(namespace string, configMapName string) (data map[string]string, err error) {
	clientSet := getClientSet()
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

func UpdateConfigMap(namespace string, configMapName string, data map[string]string) (err error) {
	clientSet := getClientSet()
	configMap, err := clientSet.CoreV1().ConfigMaps(namespace).Get(configMapName, metav1.GetOptions{})
	if err != nil {
		err = errors.Wrap(err, 0)
		return
	}

	for key, value := range data {
		configMap.Data[key] = value
	}

	_, err = clientSet.CoreV1().ConfigMaps(namespace).Update(configMap)
	if err != nil {
		err = errors.Wrap(err, 0)
		return
	}

	return
}

func CreateConfigMap(namespace string, configMapName string, data map[string]string) (err error) {
	clientSet := getClientSet()

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      configMapName,
		},
		Data: data,
	}

	_, err = clientSet.CoreV1().ConfigMaps(namespace).Create(configMap)
	if err != nil {
		err = errors.Wrap(err, 0)
		return
	}

	return
}

func DeleteConfigMap(namespace string, configMapName string) (err error) {
	clientSet := getClientSet()
	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}

	err = clientSet.CoreV1().ConfigMaps(namespace).Delete(configMapName, deleteOptions)
	if err != nil {
		err = errors.Wrap(err, 0)
		return
	}

	return
}
