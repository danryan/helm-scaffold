{{/*
Expand the name of the chart.
*/}}
{{- define "(( .Chart.Name )).name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "(( .Chart.Name )).fullname" -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "(( .Chart.Name )).service.fullname" -}}
{{- .Values.service.nameOverride | default .Chart.Name }}
{{- end -}}

{{- /*
(( .Chart.Name )).chartref prints a chart name and version.
It does minimal escaping for use in Kubernetes labels.
Example output:
  zookeeper-1.2.3
  wordpress-3.2.1_20170219
*/ -}}
{{- define "(( .Chart.Name )).chartref" -}}
  {{- replace "+" "_" .Chart.Version | printf "%s-%s" .Chart.Name -}}
{{- end -}}