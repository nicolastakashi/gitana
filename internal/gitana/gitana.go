package gitana

import (
	"context"
	"os"

	"github.com/gitana/internal/dashboardloader"
	"github.com/gitana/internal/k8sclient"
	"github.com/gitana/internal/pullcommand"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

func Start(ctx context.Context, pcmd pullcommand.Command) {
	_, err := pcmd.Repository.Get(ctx)

	if err != nil {
		os.Exit(1)
	}

	dashboards := dashboardloader.Load(pcmd.Repository.Path)

	if len(dashboards) == 0 {
		logrus.Warn("no dashboards found")
		os.Exit(0)
	}

	client, err := k8sclient.New()

	if err != nil {
		os.Exit(1)
	}

	configMaps, err := client.GetConfigMaps(pcmd.Namespace)

	if err != nil {
		os.Exit(1)
	}

	for _, dashboard := range dashboards {
		cm, err := dashboard.ToConfigMap(pcmd.Namespace, pcmd.DashboardLabels, pcmd.DashboardFolderAnnotation)

		if err != nil {
			logrus.Errorf("error to convert dashboard %v to config map", dashboard.FileName)
			continue
		}

		if ccm, ok := configMaps[dashboard.Name]; !ok {
			createConfigMap(client, &cm, &dashboard)
		} else if dashboard.NeedsToUpdate(&ccm, &cm) {
			updateConfigMap(client, &ccm, &cm, &dashboard)
		}
	}

	for _, cm := range configMaps {
		if dashboard, ok := dashboards[cm.Name]; !ok {
			deleteConfigMap(client, &cm, &dashboard)
		}
	}

	logrus.Info("dashboards sync is done")
}

func createConfigMap(client *k8sclient.K8sClient, cm *v1.ConfigMap, dashboard *dashboardloader.Dashboard) {
	logrus.Debugf("creating dashboard %v", dashboard.FileName)

	_, err := client.CreateConfigMap(cm)

	if err != nil {
		logrus.Errorf("error to create dashboard %v", dashboard.FileName)
	}
}

func updateConfigMap(client *k8sclient.K8sClient, ccm *v1.ConfigMap, cm *v1.ConfigMap, dashboard *dashboardloader.Dashboard) {
	logrus.Debugf("updating dashboard %v", dashboard.FileName)

	_, err := client.UpdateConfigMap(ccm, cm)

	if err != nil {
		logrus.Errorf("error to update dashboard %v", dashboard.FileName)
	}
}

func deleteConfigMap(client *k8sclient.K8sClient, cm *v1.ConfigMap, dashboard *dashboardloader.Dashboard) {
	err := client.DeleteConfigMap(cm)

	if err != nil {
		logrus.Errorf("error to delete dashboard %v", dashboard.FileName)
	}
}
