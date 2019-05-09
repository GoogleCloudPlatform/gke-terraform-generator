project_id = "{{.ProjectID}}"
project_name = "{{.ProjectName}}"
region = "{{.Region}}"
zones = [{{- range $i, $e := .Zones -}}{{if $i}}, {{end}}"{{$e}}"{{end}}]
service_account_id = "{{.ServiceAccountID}}"
