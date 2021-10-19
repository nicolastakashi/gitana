# Overview
Gitana is a lightweight application that will help you sync Grafana dashboards from a Git repository to Kubernetes ConfigMap and leverages the dashboard sidecar on the [Grafana helm chart](https://github.com/grafana/helm-charts/tree/main/charts/grafana) that provisions dashboard ConfigMaps created by Gitana into Grafana.

# Sync Command Flags

```bash
./gitana sync --help

The sync command pulls the Grafana dashboards from a Git repository and foreach dashboard it will creates a config map for that dashboard:

Usage:
  gitana sync [flags]

Flags:
      --dashboard.folder-annotation string   dashboard folder annotation
      --dashboard.labels string              dashboard label selector (default "grafana_dashboard=nil")
  -h, --help                                 help for sync
      --http.port string                     listem port for http endpoints (default ":9754")
      --kubeconfig string                    (optional) absolute path to the kubeconfig file
      --log.level string                     listem port for http endpoints (default "info")
      --namespace string                     namespace that will store the dashboard config map (default "default")
      --repository.auth.password string      password to perform authentication
      --repository.auth.username string      username to perform authentication
      --repository.branch string             path to clone the git repository (default "main")
      --repository.path string               path to clone the git repository
      --repository.url string                git repository url
      --sync-timer duration                  interval to sync and sync dashboards (default 5m0s)

```

# Contributing
Contributions are very welcome! See our [CONTRIBUTING.md](CONTRIBUTING.md) for more information.

## Docker images

Docker images are available on [Docker Hub](https://hub.docker.com/r/ntakashi/gitana).

## Building from source

To build Prometheus from source code, first ensure that you have a working
Go environment with [version 1.16 or greater installed](https://golang.org/doc/install).

To build the source code you can use the `make build`, which will compile in
the assets so that Gitana can be run from anywhere:

```bash

    $ mkdir -p $GOPATH/src/github.com/gitana
    $ cd $GOPATH/src/github.com/gitana
    $ git clone https://github.com/nicolastakashi/gitana.git
    $ cd gitana
    $ make build
    $ ./gitana sync <args>
```

The Makefile provides several targets:

  * *build*: build the `gitana`
  * *fmt*: format the source code
  * *vet*: check the source code for common errors
  <!-- * *test*: run the tests -->
  <!-- * *test-short*: run the short tests -->
