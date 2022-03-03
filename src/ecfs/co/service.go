package co

import (
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/go-errors/errors"
)

func GetService(namespace string, serviceName string) (service *corev1.Service, err error) {
	clientSet := getClientSet()
	opts := metav1.GetOptions{}
	service, err = clientSet.CoreV1().Services(namespace).Get(serviceName, opts)
	if err != nil {
		err = errors.Wrap(err, 0)
		return
	}

	return
}

func GetServiceWithRetries(namespace string, serviceName string, timeout time.Duration) (service *corev1.Service, err error) {
	select {
	case <-time.After(timeout):
		err = errors.WrapPrefix(err, fmt.Sprintf("Timed out getting service %v in namespace %v after %v",
			serviceName, namespace, timeout), 0)
		return
	default:
		service, err = GetService(namespace, serviceName)
	}

	return
}
