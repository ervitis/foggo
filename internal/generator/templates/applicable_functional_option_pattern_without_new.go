package templates

const AFOPTemplateWithoutNew = TemplateBase + `
type {{ .structName }}Option interface {
	apply(*{{ .structName }})
}
{{ range .fields }}{{ if ne .Ignore true}}
type {{ .Name }}Option struct {
	{{ .Name }} {{ .Type }}
}

func (o {{ .Name }}Option) apply(s *{{ $.structName }}) {
	s.{{ .Name }} = o.{{ .Name }}
}

func With{{ .Name }}({{ .Name }} {{ .Type }}) {{ $.structName }}Option {
	return func(args *{{ $.structName }}) {
		args.{{ .Name }} = {{ .Name }}
	}
}
{{ end }}{{ end }}`
