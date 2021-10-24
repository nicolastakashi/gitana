# gitana

![Version: 1.0.0](https://img.shields.io/badge/Version-1.0.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 1.16.0](https://img.shields.io/badge/AppVersion-1.16.0-informational?style=flat-square)

## Get Repo Info

```console
helm repo add gitana https://nicolastakashi.github.io/gitana
helm repo update
```

_See [helm repo](https://helm.sh/docs/helm/helm_repo/) for command documentation._

## Installing the Chart

To install the chart with the release name `my-release`:

```console
helm install my-release nicolastakashi/gitana
```

## Uninstalling the Chart

To uninstall/delete the my-release deployment:

```console
helm delete my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Values

| Key                                    | Type   | Default                                        | Description                                                                                                                               |
|----------------------------------------|--------|------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------|
| affinity                               | object | `{}`                                           |                                                                                                                                           |
| flags.dashboard.folderAnnotation       | string | `"dashboard-folder"`                           | ref: https://github.com/grafana/helm-charts/tree/main/charts/grafana#configuration sidecar.dashboards.folderAnnotation                    |
| flags.dashboard.labels                 | list   | `[{"name":"grafana-dashboard","value":"nil"}]` | ref: https://github.com/grafana/helm-charts/tree/main/charts/grafana#configuration sidecar.dashboards.label sidecar.dashboards.labelValue |
| flags.kubeconfig                       | string | `""`                                           |                                                                                                                                           |
| flags.log.level                        | string | `"info"`                                       |                                                                                                                                           |
| flags.namespace                        | string | `"gitana"`                                     |                                                                                                                                           |
| flags.repository.auth                  | object | `{}`                                           |                                                                                                                                           |
| flags.repository.branch                | string | `"main"`                                       |                                                                                                                                           |
| flags.repository.url                   | string | `"https://github.com/nicolastakashi/poc"`      |                                                                                                                                           |
| flags.syncTimer                        | string | `"5m"`                                         |                                                                                                                                           |
| fullnameOverride                       | string | `""`                                           |                                                                                                                                           |
| image.pullPolicy                       | string | `"IfNotPresent"`                               |                                                                                                                                           |
| image.repository                       | string | `"ntakashi/gitana"`                            |                                                                                                                                           |
| image.tag                              | string | `"0.1.0"`                                      |                                                                                                                                           |
| imagePullSecrets                       | list   | `[]`                                           |                                                                                                                                           |
| nameOverride                           | string | `""`                                           |                                                                                                                                           |
| nodeSelector                           | object | `{}`                                           |                                                                                                                                           |
| podAnnotations                         | object | `{}`                                           |                                                                                                                                           |
| podSecurityContext                     | object | `{}`                                           |                                                                                                                                           |
| resources                              | object | `{}`                                           |                                                                                                                                           |
| securityContext.readOnlyRootFilesystem | bool   | `true`                                         |                                                                                                                                           |
| service.port                           | int    | `80`                                           |                                                                                                                                           |
| service.type                           | string | `"ClusterIP"`                                  |                                                                                                                                           |
| serviceAccount.annotations             | object | `{}`                                           |                                                                                                                                           |
| serviceAccount.create                  | bool   | `true`                                         |                                                                                                                                           |
| serviceAccount.name                    | string | `""`                                           | If not set and create is true, a name is generated using the fullname template                                                            |
| serviceMonitor.enabled                 | bool   | `false`                                        |                                                                                                                                           |
| serviceMonitor.interval                | string | `""`                                           | ref: https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#endpoint                                              |
| serviceMonitor.scrapeTimeout           | string | `""`                                           | ref: https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#endpoint                                              |
| tolerations                            | list   | `[]`                                           |                                                                                                                                           |

