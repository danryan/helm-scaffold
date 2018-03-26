{{- define "common.ingress.tpl" -}}
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: {{ template "(( .Chart.Name )).fullname" . }}
  labels:
    app: {{ template "(( .Chart.Name )).name" . }}
    chart: {{ template "(( .Chart.Name )).chartref" . }}
    heritage: {{ .Release.Service | quote }}
    release: {{ .Release.Name | quote }}
  {{- if .Values.ingress.annotations }}
  annotations:
    {{ include "common.annote" .Values.ingress.annotations | indent 4 }}
  {{- end }}
spec:
  rules:
  {{- range $host := .Values.ingress.hosts }}
  - host: {{ $host }}
    http:
      paths:
      - path: /
        backend:
          serviceName: {{ template "(( .Chart.Name )).fullname" $ }}
          servicePort: 80
  {{- end }}
  {{- if .Values.ingress.tls }}
  tls:
{{ toYaml .Values.ingress.tls | indent 4 }}
  {{- end -}}