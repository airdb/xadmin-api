package gencode

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/airdb/xadmin-api/pkg/protockit/util"
	"google.golang.org/protobuf/compiler/protogen"
)

type converterGenerator struct {
	gen      *protogen.Plugin
	file     *protogen.File
	service  *protogen.Service
	name     string
	reMethod *regexp.Regexp
	reInOut  *regexp.Regexp
}

func NewConverterGenerator(
	gen *protogen.Plugin,
	file *protogen.File,
	service *protogen.Service,
) *converterGenerator {
	return &converterGenerator{
		gen:     gen,
		file:    file,
		service: service,
		reMethod: regexp.MustCompile(
			fmt.Sprintf(`^(%s)([A-Z].*)$`, strings.Join(actionsWords, "|")),
		),
		reInOut: regexp.MustCompile(`^([A-Z].*)(Request|Response)$`),
	}
}

func (r *converterGenerator) Run() error {
	if r.service.GoName[len(r.service.GoName)-7:] != "Service" {
		return fmt.Errorf("%s should end with Service", r.service.GoName)
	}

	r.name = strings.ToLower(r.service.GoName[0 : len(r.service.GoName)-7])

	if len(r.service.Methods) == 0 {
		return nil
	}

	filename := r.file.GeneratedFilenamePrefix + `_converter.go`

	g := r.gen.NewGeneratedFile(filename, controllersPackage)

	g.P("package controllers")

	if err := r.genImpl(g); err != nil {
		return err
	}

	if err := r.genNew(g); err != nil {
		return err
	}

	if err := r.genMethods(g); err != nil {
		return err
	}

	return nil
}

func (r converterGenerator) genImpl(g *protogen.GeneratedFile) error {
	g.P()
	g.P("type ", util.LcFirst(r.service.GoName), "Converter struct {")
	g.P("}")

	return nil
}

func (r converterGenerator) genNew(g *protogen.GeneratedFile) error {
	g.P()
	g.P("func new", r.service.GoName, "Converter() ",
		g.QualifiedGoIdent(r.file.GoImportPath.Ident(r.service.GoName+"Converter")),
		" {")
	g.P("return &", util.LcFirst(r.service.GoName), "Converter{")
	g.P("}")
	g.P("}")

	return nil
}

func (r converterGenerator) genMethods(g *protogen.GeneratedFile) error {
	g.P()
	for _, entity := range guessEntity(r.service) {
		r.genConvertProtoToModel(g, entity)
	}

	return nil
}

func (r converterGenerator) genConvertProtoToModel(
	g *protogen.GeneratedFile, message *protogen.Message,
) {
	protoIdent := g.QualifiedGoIdent(message.GoIdent)
	modelIdent := g.QualifiedGoIdent(dataPackage.Ident(
		message.GoIdent.GoName + "Entity"))

	pmFunc := fmt.Sprintf("FromProto%sToModel%s",
		message.GoIdent.GoName, message.GoIdent.GoName)
	pmSign := pmFunc
	pmSign += "(in *" + protoIdent
	pmSign += ") *" + modelIdent

	g.P("func (c ", util.LcFirst(r.service.GoName), "Converter) ", pmSign, "{")
	g.P("if in == nil {")
	g.P("return nil")
	g.P("}")
	g.P("")
	g.P("return &", modelIdent, "{}")
	g.P("}")
	g.P()

	pmsSign := fmt.Sprintf("FromProto%sToModel%s",
		util.Pluralize.Plural(message.GoIdent.GoName),
		util.Pluralize.Plural(message.GoIdent.GoName))
	pmsSign += "(in []*" + protoIdent
	pmsSign += ") []*" + modelIdent

	g.P("func (c ", util.LcFirst(r.service.GoName), "Converter) ", pmsSign, "{")
	g.P("if in == nil || len(in) == 0 {")
	g.P("return nil")
	g.P("}")
	g.P("")
	g.P("res := make([]*", modelIdent, ", len(in))")
	g.P("for i := 0; i < len(in); i++ {")
	g.P("res[i] = c.", pmFunc, "(in[i])")
	g.P("}")
	g.P("")
	g.P("return res")
	g.P("}")
	g.P()

	mpFunc := fmt.Sprintf("FromModel%sToProto%s",
		message.GoIdent.GoName, message.GoIdent.GoName)
	mpSign := mpFunc
	mpSign += "(in *" + modelIdent
	mpSign += ") *" + protoIdent

	g.P("func (c ", util.LcFirst(r.service.GoName), "Converter) ", mpSign, "{")
	g.P("if in == nil {")
	g.P("return nil")
	g.P("}")
	g.P("")
	g.P("return &", protoIdent, "{}")
	g.P("}")
	g.P()

	mpsSign := fmt.Sprintf("FromModel%sToProto%s",
		util.Pluralize.Plural(message.GoIdent.GoName),
		util.Pluralize.Plural(message.GoIdent.GoName))
	mpsSign += "(in []*" + protoIdent
	mpsSign += ") []*" + modelIdent

	g.P("func (c ", util.LcFirst(r.service.GoName), "Converter) ", mpsSign, "{")
	g.P("if in == nil || len(in) == 0 {")
	g.P("return nil")
	g.P("}")
	g.P("")
	g.P("res := make([]*", modelIdent, ", len(in))")
	g.P("for i := 0; i < len(in); i++ {")
	g.P("res[i] = c.", mpFunc, "(in[i])")
	g.P("}")
	g.P("return res")
	g.P("")
	g.P("}")
	g.P()
}
