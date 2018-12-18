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

func UpdateSecrets(namespace string, secretsName string, data map[string][]byte) (err error) {
	clientSet := clientSet()
	secrets, err := clientSet.CoreV1().Secrets(namespace).Get(secretsName, metav1.GetOptions{})
	if err != nil {
		err = errors.Wrap(err, 0)
		return
	}

	for key, value := range data {
		secrets.Data[key] = value
	}

	_, err = clientSet.CoreV1().Secrets(namespace).Update(secrets)
	if err != nil {
		err = errors.Wrap(err, 0)
		return
	}

	return
}

func CreateSecrets(namespace string, secretsName string, data map[string][]byte) (err error) {
	clientSet := clientSet()
	secrets, err := clientSet.CoreV1().Secrets(namespace).Get(secretsName, metav1.GetOptions{})
	if err != nil {
		err = errors.Wrap(err, 0)
		return
	}

	for key, value := range data {
		secrets.Data[key] = value
	}

	_, err = clientSet.CoreV1().Secrets(namespace).Create(secrets)
	if err != nil {
		err = errors.Wrap(err, 0)
		return
	}

	return
}
