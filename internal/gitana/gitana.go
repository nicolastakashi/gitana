package gitana

import (
	"context"
	"errors"
	"time"

	"github.com/gitana/internal/command"
	"github.com/gitana/internal/dashboardloader"
	"github.com/gitana/internal/k8sclient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
)

var lastSuccessfulSync = promauto.NewGauge(prometheus.GaugeOpts{
	Namespace: "gitana",
	Name:      "last_success_sync_timestamp_seconds",
	Help:      "Unix timestamp of the last successful dashboard sync in seconds",
})

var syncLatency = promauto.NewSummary(
	prometheus.SummaryOpts{
		Namespace: "gitana",
		Name:      "sync_time_seconds",
		Help:      "Time taken by the sync operation",
	},
)

var syncSuccessTotal = promauto.NewCounter(
	prometheus.CounterOpts{
		Namespace: "gitana",
		Name:      "sync_total_success",
		Help:      "Total number of successful sync operations",
	},
)

var syncErrorTotal = promauto.NewCounter(
	prometheus.CounterOpts{
		Namespace: "gitana",
		Name:      "sync_total_error",
		Help:      "Total number of sync operations with errors",
	},
)

func Start(ctx context.Context, pcmd command.Sync) error {

	t := time.NewTimer(1 * time.Millisecond)

	for {
		select {
		case <-t.C:
			logrus.Info("start new sync operation")
			if err := start(ctx, pcmd); err != nil {
				syncErrorTotal.Inc()
				return err
			} else {
				logrus.Info("dashboards sync operation is done")
				syncSuccessTotal.Inc()
				lastSuccessfulSync.SetToCurrentTime()
			}

			t.Reset(pcmd.SyncTimer)
		case <-ctx.Done():
			logrus.Info("shut down gitana syncer")
			return nil
		}
	}
}

func start(ctx context.Context, pcmd command.Sync) error {
	timer := prometheus.NewTimer(syncLatency)

	client, err := k8sclient.New(pcmd.KubeConfig)

	if err != nil {
		return err
	}

	if pcmd.Repository.Auth.AuthSecretName != "" {
		secret, err := client.GetSecret(pcmd.Namespace, pcmd.Repository.Auth.AuthSecretName)
		if err != nil {
			return err
		}

		secretData := secret.Data["auth.yaml"]

		if secretData == nil {
			return errors.New("auth secret there is no auth.yaml")
		}

		err = yaml.Unmarshal(secretData, &pcmd.Repository.Auth)

		if err != nil {
			logrus.Errorf("error to unmarshal auth secret %v", err)
			return err
		}
	}

	_, err = pcmd.Repository.Get(ctx)

	if err != nil {
		return err
	}

	dashboards, err := dashboardloader.Load(pcmd.Repository.GetPath())

	if err != nil {
		return err
	}

	if len(dashboards) == 0 {
		logrus.Warn("no dashboards found")
		return nil
	}

	configMaps, err := client.GetConfigMaps(pcmd.Namespace)

	if err != nil {
		return err
	}

	for _, dashboard := range dashboards {
		cm, err := dashboard.ToConfigMap(pcmd.Namespace, pcmd.DashboardLabels, pcmd.DashboardFolderAnnotation)

		if err != nil {
			logrus.Errorf("error to convert dashboard %v to config map", dashboard.FileName)
			continue
		}

		if ccm, ok := configMaps[dashboard.Name]; !ok {
			createConfigMap(client, cm, dashboard)
		} else if dashboard.NeedsToUpdate(ccm, cm) {
			updateConfigMap(client, ccm, cm, dashboard)
		}
	}

	for _, cm := range configMaps {
		if dashboard, ok := dashboards[cm.Name]; !ok {
			deleteConfigMap(client, cm, &dashboard)
		}
	}

	timer.ObserveDuration()

	return nil
}

func createConfigMap(client *k8sclient.K8sClient, cm v1.ConfigMap, dashboard dashboardloader.Dashboard) {
	logrus.Debugf("creating dashboard %v", dashboard.FileName)

	_, err := client.CreateConfigMap(cm)

	if err != nil {
		logrus.Errorf("error to create dashboard %v", dashboard.FileName)
	}
}

func updateConfigMap(client *k8sclient.K8sClient, ccm v1.ConfigMap, cm v1.ConfigMap, dashboard dashboardloader.Dashboard) {
	logrus.Debugf("updating dashboard %v", dashboard.FileName)

	_, err := client.UpdateConfigMap(ccm, cm)

	if err != nil {
		logrus.Errorf("error to update dashboard %v", dashboard.FileName)
	}
}

func deleteConfigMap(client *k8sclient.K8sClient, cm v1.ConfigMap, dashboard *dashboardloader.Dashboard) {
	err := client.DeleteConfigMap(cm)

	if err != nil {
		logrus.Errorf("error to delete dashboard %v", dashboard.FileName)
	}
}
