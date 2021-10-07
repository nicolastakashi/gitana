package k8sclient

import (
	"context"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type K8sClient struct {
	client kubernetes.Clientset
}

func New() (*K8sClient, error) {
	config, err := rest.InClusterConfig()

	if err != nil {
		logrus.Error("error while get in cluster config: %v", err)
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		logrus.Error("error while creating k8s client: %v", err)
		return nil, err
	}

	return &K8sClient{
		client: *clientset,
	}, nil
}

func (kc *K8sClient) GetConfigMaps(namespace string) (map[string]v1.ConfigMap, error) {

	configMaps, err := kc.client.CoreV1().ConfigMaps(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/managed-by=gitana",
	})
	if err != nil {
		logrus.Error("error to get dashboard config maps %v", err)
		return nil, err
	}

	cmMap := map[string]v1.ConfigMap{}

	for _, item := range configMaps.Items {
		cmMap[item.Name] = item
	}

	return cmMap, nil
}

func (kc *K8sClient) CreateConfigMap(cm v1.ConfigMap) (*v1.ConfigMap, error) {
	ncm, err := kc.client.CoreV1().ConfigMaps(cm.Namespace).Create(context.TODO(), &cm, metav1.CreateOptions{})
	if err != nil {
		logrus.Error("error to create config map %v", err)
		return &v1.ConfigMap{}, err
	}
	return ncm, nil
}

func (kc *K8sClient) UpdateConfigMap(ccm v1.ConfigMap, dcm v1.ConfigMap) (*v1.ConfigMap, error) {
	dcm.SetResourceVersion(ccm.GetResourceVersion())

	cm, err := kc.client.CoreV1().ConfigMaps(ccm.Namespace).Update(context.TODO(), &dcm, metav1.UpdateOptions{})
	if err != nil {
		logrus.Error("error to update config map %v", err)
		return cm, err
	}
	return cm, nil
}

func (kc *K8sClient) DeleteConfigMap(cm v1.ConfigMap) error {
	err := kc.client.CoreV1().ConfigMaps(cm.Namespace).Delete(context.TODO(), cm.Name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}
