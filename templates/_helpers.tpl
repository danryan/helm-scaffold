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
