package command

import (
	"errors"
	"time"

	"github.com/gitana/internal/gitmanager"
	"github.com/sirupsen/logrus"
)

type Sync struct {
	DashboardLabels           string
	Namespace                 string
	Repository                gitmanager.Repository
	DashboardFolderAnnotation string
	SyncTimer                 time.Duration
	LogLevel                  string
	KubeConfig                string
}

func (pcmd *Sync) Validate() error {
	logrus.Info("Validating flags")

	if pcmd.Namespace == "" {
		return errors.New("namespace flag can not be null")
	}

	if pcmd.DashboardFolderAnnotation == "" {
		return errors.New("dashboard.folder-annotation flag can not be null")
	}

	if err := pcmd.Repository.Validate(); err != nil {
		return err
	}

	return nil
}
