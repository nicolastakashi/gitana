{{/*
Expand the name of the chart.
*/}}
{{- define "gitana.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "gitana.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "gitana.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "gitana.labels" -}}
helm.sh/chart: {{ include "gitana.chart" . }}
{{ include "gitana.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "gitana.selectorLabels" -}}
app.kubernetes.io/name: {{ include "gitana.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "gitana.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "gitana.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{- define "dashboard.labels" -}}
{{- if .Values.flags.dashboard.labels }}
{{- range .Values.flags.dashboard.labels }}{{(print .name "=" .value ) }},{{- end }}
{{- end }}
{{- end }}

{{- define "gitana.authSecretName" -}}
{{- if .Values.authSecret.secretname }}
{{- .Values.authSecret.secretname }}
{{- else }}
{{- printf "%s-auth-secret" (include "gitana.fullname" .) }}
{{- end }}
{{- end }}