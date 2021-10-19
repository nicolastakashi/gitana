/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gitana/internal/command"
	"github.com/gitana/internal/gitana"
	"github.com/gitana/internal/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync grafana dashboards from Git repository and to configMap",
	Long:  `The sync command pulls the Grafana dashboards from a Git repository and foreach dashboard it will creates a config map for that dashboard:`,
	Run: func(cmd *cobra.Command, args []string) {

		if err := logging.Configure(pcmd.LogLevel); err != nil {
			os.Exit(1)
		}

		logrus.Info("Welcome to gitana...")

		if err := pcmd.Validate(); err != nil {
			logrus.Error(err)
			os.Exit(1)
		}

		ctx, cancel := context.WithCancel(context.Background())
		wg, ctx := errgroup.WithContext(ctx)

		term := make(chan os.Signal, 1)
		signal.Notify(term, os.Interrupt, syscall.SIGTERM)

		srv := createHttpServer(serverPort)

		err := prometheus.DefaultRegisterer.Register(version.NewCollector("gitana"))

		if err != nil {
			logrus.Errorf("error to register version collector %v", err)
		}

		logrus.Info("listen on " + serverPort)

		wg.Go(func() error {
			if err := srv.ListenAndServe(); err != http.ErrServerClosed {
				logrus.Error("msg", "http server error", err, err)
				return err
			}
			return nil
		})

		wg.Go(func() error {
			if err := gitana.Start(ctx, *pcmd); err != nil {
				return err
			}
			return nil
		})

		select {
		case <-term:
			logrus.Info("received SIGTERM, exiting gracefully...")
		case <-ctx.Done():
		}

		if err := srv.Shutdown(ctx); err != nil {
			logrus.Error("server shutdown error ", err)
		}

		cancel()

		if err := wg.Wait(); err != nil {
			logrus.Error("unhandled error received. Exiting...", err)
			os.Exit(1)
		}

		os.Exit(0)
	},
}

func createHttpServer(port string) *http.Server {
	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/-/health", func(rw http.ResponseWriter, _ *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Header().Set("Content-Type", "application/json")

		resp := map[string]string{
			"message": "Healthy",
		}

		jsonResp, err := json.Marshal(resp)

		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}

		rw.Write(jsonResp)
	})

	srv := &http.Server{
		Addr:     port,
		Handler:  mux,
		ErrorLog: &log.Logger{},
	}
	return srv
}

var pcmd = &command.Sync{}
var serverPort = ""

func init() {
	syncCmd.Flags().StringVar(&serverPort, "http.port", ":9754", "listem port for http endpoints")
	syncCmd.Flags().StringVar(&pcmd.Repository.Url, "repository.url", "", "git repository url")
	syncCmd.Flags().StringVar(&pcmd.Repository.Path, "repository.path", "", "path to clone the git repository")
	syncCmd.Flags().StringVar(&pcmd.Repository.Auth.Username, "repository.auth.username", "", "username to perform authentication")
	syncCmd.Flags().StringVar(&pcmd.Repository.Auth.Password, "repository.auth.password", "", "password to perform authentication")
	syncCmd.Flags().StringVar(&pcmd.Repository.Branch, "repository.branch", "main", "path to clone the git repository")
	syncCmd.Flags().StringVar(&pcmd.Namespace, "namespace", "default", "namespace that will store the dashboard config map")
	syncCmd.Flags().StringVar(&pcmd.DashboardLabels, "dashboard.labels", "grafana_dashboard=nil", "dashboard label selector")
	syncCmd.Flags().StringVar(&pcmd.DashboardFolderAnnotation, "dashboard.folder-annotation", "", "dashboard folder annotation")
	syncCmd.Flags().DurationVar(&pcmd.SyncTimer, "sync-timer", 300*time.Second, "interval to sync and sync dashboards")
	syncCmd.Flags().StringVar(&pcmd.LogLevel, "log.level", logrus.InfoLevel.String(), "listem port for http endpoints")
	syncCmd.Flags().StringVar(&pcmd.KubeConfig, "kubeconfig", "", "(optional) absolute path to the kubeconfig file")
	rootCmd.AddCommand(syncCmd)
}
