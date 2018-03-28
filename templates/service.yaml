apiVersion: v1
kind: Service
metadata:
  name: {{ template "(( .Chart.Name )).fullname" . }}
  labels:
    app: {{ template "(( .Chart.Name )).name" . }}
    chart: {{ template "(( .Chart.Name )).chartref" . }}
    heritage: {{ .Release.Service | quote }}
    release: {{ .Release.Name | quote }}
spec:
  type: {{ .Values.service.type }}
  ports:
  - name: http
    port: 80
    targetPort: http
  selector:
    app: {{ template "(( .Chart.Name )).name" . }}
    release: {{ .Release.Name | quote }}
{{- end -}}
