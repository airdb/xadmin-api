package genextends

import (
	"bytes"
	"fmt"
	"text/template"

	annov1 "github.com/airdb/xadmin-api/genproto/annotation/v1"
	"google.golang.org/protobuf/compiler/protogen"
	// "github.com/stoewer/go-strcase"
	// "go.einride.tech/aip/reflect/aipreflect"
	// "go.einride.tech/aip/resourcename"
	// "google.golang.org/protobuf/compiler/protogen"
	// "google.golang.org/protobuf/reflect/protoreflect"
	// "google.golang.org/protobuf/reflect/protoregistry"
)

const (
	loPackage = protogen.GoImportPath("github.com/samber/lo")
)

type messageGenerator struct {
	message *protogen.Message
	option  *annov1.MessageDescriptor
	options map[string]*annov1.FieldDescriptor

	messageName string

	actionName    string
	actionContent []any

	maskMapName    string
	maskMapContent []any
}

func NewMessageGenerator(
	message *protogen.Message,
	option *annov1.MessageDescriptor,
	options map[string]*annov1.FieldDescriptor,
) messageGenerator {
	return messageGenerator{
		message: message,
		option:  option,
		options: options,
	}
}

func (r messageGenerator) Run(g *protogen.GeneratedFile) error {
	if r.message == nil || r.options == nil || len(r.options) == 0 {
		return nil
	}

	if len(r.message.Fields) == 0 {
		return nil
	}

	g.QualifiedGoIdent(loPackage.Ident(""))

	return r.genExtend(g)
}

var messageExtendTemplate = `
var {{.actionName}} = map[string][]string {
	{{- range .actionContent }}
	"{{.name}}": {{"{"}}{{ sliceJoin .actions }}{{"}"}},
	{{- end }}
}
var {{.maskMapName}} = map[string]string {
	{{- range .maskMapContent }}
	"{{.k}}": "{{.v}}",
	{{- end }}
}

func (x *{{.messageName}}) KitFieldsByActions(action ...string) []string {
	fields := []string{}
	for field, actions := range {{.actionName}} {
		if len(lo.Intersect(actions, action)) > 0 ||
			lo.IndexOf(actions, "*") >= 0 {
			fields = append(fields, field)
		}
	}
	return fields
}

func (x *{{.messageName}}) KitMaskMap(field string) string {
	if val, ok := {{.maskMapName}}[field]; ok {
		return val
	}
	return field
}
`

func (r messageGenerator) genExtend(g *protogen.GeneratedFile) error {
	name := r.message.GoIdent.GoName
	actionName := fmt.Sprintf("KIT_MESSAEG_%s_ACTIONS", name)
	maskMapName := fmt.Sprintf("KIT_MESSAEG_%s_MASKMAP", name)

	buf := new(bytes.Buffer)
	tmpl, err := template.New("messsage_extend").
		Funcs(template.FuncMap{
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
		}).
		Parse(messageExtendTemplate)
	if err != nil {
		panic(err)
	}

	data := map[string]any{}

	data["messageName"] = name
	data["actionName"] = actionName
	data["actionContent"] = func() []any {
		res := []any{}
		for _, field := range r.message.Fields {
			option, ok := r.options[field.GoIdent.GoName]
			if !ok || len(option.Actions) == 0 {
				continue
			}
			res = append(res, map[string]any{
				"name":    field.Desc.Name(),
				"actions": option.Actions,
			})
		}
		return res
	}()
	data["maskMapName"] = maskMapName
	data["maskMapContent"] = func() []any {
		res := []any{}
		if r.option.MaskMap != nil && len(r.option.MaskMap) > 0 {
			for k, v := range r.option.MaskMap {
				res = append(res, map[string]any{
					"k": k,
					"v": v,
				})
			}
		}
		return res
	}()
	if err := tmpl.Execute(buf, data); err != nil {
		panic(err)
	}
	g.P(buf.String())

	return nil
}
