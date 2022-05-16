package generror

import (
	"bytes"
	"text/template"
)

var errorsTemplate = `
{{ range . }}

func Is{{.CamelValue}}(err error) bool {
	if err == nil {
		return false
	}
	e := errorskit.FromError(err)
	return e.Reason == {{.Name}}_{{.Value}}.String() && e.Code == {{.HTTPCode}} 
}

func Error{{.CamelValue}}(format string, args ...interface{}) *errorskit.Error {
	 return errorskit.New({{.HTTPCode}}, {{.Name}}_{{.Value}}.String(), fmt.Sprintf(format, args...))
}
{{- end }}
`

type valueInfo struct {
	Name       string
	Value      string
	HTTPCode   int32
	CamelValue string
}

type valueInfos []*valueInfo

func (e *valueInfos) execute() string {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("errors").Parse(errorsTemplate)
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, e); err != nil {
		panic(err)
	}
	return buf.String()
}
