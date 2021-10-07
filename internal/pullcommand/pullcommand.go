package pullcommand

import (
	"errors"
	"time"

	"github.com/gitana/internal/gitmanager"
	"github.com/sirupsen/logrus"
)

type Command struct {
	DashboardLabels           string
	Namespace                 string
	Repository                gitmanager.Repository
	DashboardFolderAnnotation string
	SyncTimer                 time.Duration
}

func (pcmd *Command) Validate() error {
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
