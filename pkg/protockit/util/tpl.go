package util

import (
	"bytes"
	"text/template"

	"github.com/gobeam/stringy"
	"google.golang.org/protobuf/compiler/protogen"
)

func NewTpl(name string) *template.Template {
	tpl := template.New(name).Funcs(template.FuncMap{
		"sliceJoin": func(items []string) string {
			buf := new(bytes.Buffer)
			for k, v := range items {
				buf.WriteString(`"` + v + `"`)
				if k < len(items) {
					buf.WriteString(`,`)
				}
			}
			return buf.String()
		},
		"lcFirst": LcFirst,
		"outputParams": func(output *protogen.Message) string {
			return ""
		},
	})

	return tpl
}

func LcFirst(s string) string {
	return stringy.New(s).LcFirst()
}

func UcFirst(s string) string {
	return stringy.New(s).UcFirst()
}
