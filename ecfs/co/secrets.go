package co

import (
	"github.com/elastifile/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetSecret(namespace string, secretName string) (data map[string][]byte, err error) {
	clientSet := clientSet()
	secret, err := clientSet.CoreV1().Secrets(namespace).Get(secretName, metav1.GetOptions{})
	data = secret.Data
	return
}

func GetSecretValue(namespace string, secretName string, key string) (value string, err error) {
	secrets, err := GetSecret(namespace, secretName)
	if err != nil {
		err = errors.WrapPrefix(err, "Failed to get CO secrets", 0)
		return
	}

	data, ok := secrets[key]
	if !ok {
		err = errors.Errorf("Key %v not found in config map %v (namespace: %v)", key, secretName, namespace)
	}

	value = string(data)
	return
}
