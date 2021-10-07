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
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gitana/internal/gitana"
	"github.com/gitana/internal/pullcommand"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull grafana dashboards from Git repository and creates the required configMap",
	Long:  `The pull command pulls the Grafana dashboards from a Git repository and foreach dashboard it will creates a config map for that dashboard:`,
	Run: func(cmd *cobra.Command, args []string) {
		os.Setenv("KUBERNETES_SERVICE_HOST", "127.0.0.1")
		os.Setenv("KUBERNETES_SERVICE_PORT", "33295")

		logrus.SetLevel(logrus.DebugLevel)

		logrus.Info("Welcome to gitana...")

		if err := pcmd.Validate(); err != nil {
			os.Exit(1)
		}

		ctx, cancel := context.WithCancel(context.Background())
		wg, ctx := errgroup.WithContext(ctx)

		term := make(chan os.Signal, 1)
		signal.Notify(term, os.Interrupt, syscall.SIGTERM)

		srv := createHttpServer(serverPort)

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

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	srv := &http.Server{
		Addr:     port,
		Handler:  mux,
		ErrorLog: &log.Logger{},
	}
	return srv
}

var pcmd = &pullcommand.Command{}
var serverPort = ""

func init() {
	pullCmd.Flags().StringVar(&pcmd.Repository.Url, "repository.url", "", "git repository url")
	pullCmd.Flags().StringVar(&pcmd.Repository.Path, "repository.path", "", "path to clone the git repository")
	pullCmd.Flags().StringVar(&pcmd.Repository.Branch, "repository.branch", "main", "path to clone the git repository")
	pullCmd.Flags().StringVar(&pcmd.Namespace, "namespace", "default", "namespace that will store the dashboard config map")
	pullCmd.Flags().StringVar(&pcmd.DashboardLabels, "dashboard.labels", "grafana_dashboard=nil", "dashboard label selector")
	pullCmd.Flags().StringVar(&pcmd.DashboardFolderAnnotation, "dashboard.folder-annotation", "", "dashboard folder annotation")
	pullCmd.Flags().StringVar(&serverPort, "http.port", ":9754", "listem port for http endpoints")
	pullCmd.Flags().DurationVar(&pcmd.SyncTimer, "sync-timer", 10*time.Second, "interval to pull and sync dashboards")
	rootCmd.AddCommand(pullCmd)
}
