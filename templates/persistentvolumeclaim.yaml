apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: {{ template "(( .Chart.Name )).fullname" . }}
  labels:
    app: {{ template "(( .Chart.Name )).name" . }}
    chart: {{ template "(( .Chart.Name )).chartref" . }}
    heritage: {{ .Release.Service | quote }}
    release: {{ .Release.Name | quote }}
spec:
  accessModes:
    - {{ .Values.persistence.accessMode | quote }}
  resources:
    requests:
      storage: {{ .Values.persistence.size | quote }}
{{- if .Values.persistence.storageClass }}
{{- if (eq "-" .Values.persistence.storageClass) }}
  storageClassName: ""
{{- else }}
  storageClassName: "{{ .Values.persistence.storageClass }}"
{{- end }}
{{- end }}
{{- end -}}
{{- define "common.persistentvolumeclaim" -}}
{{- $top := first . -}}
{{- if and $top.Values.persistence.enabled (not $top.Values.persistence.existingClaim) -}}
{{- template "common.util.merge" (append . "common.persistentvolumeclaim.tpl") -}}
{{- end -}}
