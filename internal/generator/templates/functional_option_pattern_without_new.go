package templates

const FOPTemplateWithoutNew = TemplateBase + `
type {{ .structName }}Option func(*{{ .structName }})

{{ range .fields }}{{ if ne .Ignore true}}
func With{{ .Name }}({{ .Name }} {{ .Type }}) {{ $.structName }}Option {
	return func(args *{{ $.structName }}) {
		args.{{ .Name }} = {{ .Name }}
	}
}
{{ end }}{{ end }}`
