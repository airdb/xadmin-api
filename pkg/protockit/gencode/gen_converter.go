package gencode

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/airdb/xadmin-api/pkg/protockit/util"
	"github.com/gobeam/stringy"
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
	r.name = strings.ToLower(r.service.GoName[0 : len(r.service.GoName)-7])

	if len(r.service.Methods) == 0 {
		return nil
	}

	filename := path.Join("app", "controllers",
		stringy.New(r.name).SnakeCase().LcFirst()+`_converter.go`)

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
	g.P("type ", util.LcFirst(r.name), "Convert struct {")
	g.P("}")

	return nil
}

func (r converterGenerator) genNew(g *protogen.GeneratedFile) error {
	g.P()
	g.P("func new", util.UcFirst(r.name), "Convert() *",
		r.name+"Convert",
		" {")
	g.P("return &", util.LcFirst(r.name), "Convert{")
	g.P("}")
	g.P("}")

	return nil
}

func (r converterGenerator) genMethods(g *protogen.GeneratedFile) error {
	g.P()
	createMsg := findCreateRequest(r.service)
	for _, entity := range guessEntity(r.service) {
		r.genConvertProtoToModel(g, entity)
		if createMsg != nil {
			r.genConvertProtoCreateToModel(g, entity, createMsg)
		}
	}

	return nil
}

func (r converterGenerator) genConvertProtoToModel(
	g *protogen.GeneratedFile, entity *protogen.Message,
) {
	protoIdent := g.QualifiedGoIdent(entity.GoIdent)
	modelIdent := g.QualifiedGoIdent(dataPackage.Ident(
		entity.GoIdent.GoName + "Entity"))

	pmFunc := fmt.Sprintf("FromProto%sToModel%s",
		entity.GoIdent.GoName, entity.GoIdent.GoName)
	pmSign := pmFunc
	pmSign += "(in *" + protoIdent
	pmSign += ") *" + modelIdent

	g.P("func (c ", util.LcFirst(r.name), "Convert) ", pmSign, "{")
	g.P("if in == nil {")
	g.P("return nil")
	g.P("}")
	g.P("")
	g.P("return &", modelIdent, "{}")
	g.P("}")
	g.P()

	pmsSign := fmt.Sprintf("FromProto%sToModel%s",
		util.Pluralize.Plural(entity.GoIdent.GoName),
		util.Pluralize.Plural(entity.GoIdent.GoName))
	pmsSign += "(in []*" + protoIdent
	pmsSign += ") []*" + modelIdent

	g.P("func (c ", util.LcFirst(r.name), "Convert) ", pmsSign, "{")
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
		entity.GoIdent.GoName, entity.GoIdent.GoName)
	mpSign := mpFunc
	mpSign += "(in *" + modelIdent
	mpSign += ") *" + protoIdent

	g.P("func (c ", util.LcFirst(r.name), "Convert) ", mpSign, "{")
	g.P("if in == nil {")
	g.P("return nil")
	g.P("}")
	g.P("")
	g.P("return &", protoIdent, "{}")
	g.P("}")
	g.P()

	mpsSign := fmt.Sprintf("FromModel%sToProto%s",
		util.Pluralize.Plural(entity.GoIdent.GoName),
		util.Pluralize.Plural(entity.GoIdent.GoName))
	mpsSign += "(in []*" + modelIdent
	mpsSign += ") []*" + protoIdent

	g.P("func (c ", util.LcFirst(r.name), "Convert) ", mpsSign, "{")
	g.P("if in == nil || len(in) == 0 {")
	g.P("return nil")
	g.P("}")
	g.P("")
	g.P("res := make([]*", protoIdent, ", len(in))")
	g.P("for i := 0; i < len(in); i++ {")
	g.P("res[i] = c.", mpFunc, "(in[i])")
	g.P("}")
	g.P("return res")
	g.P("")
	g.P("}")
	g.P()
}

func (r converterGenerator) genConvertProtoCreateToModel(
	g *protogen.GeneratedFile, entity *protogen.Message, create *protogen.Message,
) {
	modelIdent := g.QualifiedGoIdent(dataPackage.Ident(
		entity.GoIdent.GoName + "Entity"))
	createIdent := g.QualifiedGoIdent(create.GoIdent)

	pmCreateFunc := fmt.Sprintf("FromProtoCreate%sToModel%s",
		entity.GoIdent.GoName, entity.GoIdent.GoName)
	pmCreateSign := pmCreateFunc
	pmCreateSign += "(in *" + createIdent
	pmCreateSign += ") *" + modelIdent

	g.P("func (c ", util.LcFirst(r.name), "Convert) ", pmCreateSign, "{")
	g.P("if in == nil {")
	g.P("return nil")
	g.P("}")
	g.P("")
	g.P("return &", modelIdent, "{}")
	g.P("}")
	g.P()
}
