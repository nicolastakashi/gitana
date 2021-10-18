# Overview
Gitana is a lightweight application that will help you sync Grafana dashboards from a Git repository to Kubernetes ConfigMap and leverages the dashboard sidecar on the [Grafana helm chart](https://github.com/grafana/helm-charts/tree/main/charts/grafana) that provisions dashboard ConfigMaps created by Gitana into Grafana.

# Contributing
Contributions are very welcome! See our [CONTRIBUTING.md](CONTRIBUTING.md) for more information.

## Docker images

Docker images are available on [Quay.io](https://quay.io) or [Docker Hub](https://hub.docker.com/).

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
