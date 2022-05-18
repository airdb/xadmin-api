package gencode

import (
	"fmt"
	"strings"

	"github.com/airdb/xadmin-api/pkg/protockit/util"
	"google.golang.org/protobuf/compiler/protogen"
)

type validationGenerator struct {
	gen     *protogen.Plugin
	file    *protogen.File
	service *protogen.Service
	name    string
}

func NewValidationGenerator(
	gen *protogen.Plugin,
	file *protogen.File,
	service *protogen.Service,
) *validationGenerator {
	return &validationGenerator{
		gen:     gen,
		file:    file,
		service: service,
	}
}

func (r *validationGenerator) Run() error {
	if r.service.GoName[len(r.service.GoName)-7:] != "Service" {
		return fmt.Errorf("%s should end with Service", r.service.GoName)
	}

	r.name = strings.ToLower(r.service.GoName[0 : len(r.service.GoName)-7])

	if len(r.service.Methods) == 0 {
		return nil
	}

	filename := r.file.GeneratedFilenamePrefix + `.go`

	g := r.gen.NewGeneratedFile(filename, validationsPackage)

	g.P("package validations")

	if err := r.genInterface(g); err != nil {
		return err
	}

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

func (r validationGenerator) genInterface(g *protogen.GeneratedFile) error {
	g.P()
	g.P("type ", util.LcFirst(r.service.GoName), "Validations interface {")
	for _, method := range r.service.Methods {
		g.P(method.GoName, validationMethodSignature(g, method))
	}
	g.P("}")

	return nil
}

func (r validationGenerator) genImpl(g *protogen.GeneratedFile) error {
	g.P()
	g.P("type ", util.LcFirst(r.service.GoName), "Validations struct {")
	g.P("}")

	return nil
}

func (r validationGenerator) genNew(g *protogen.GeneratedFile) error {
	g.P()
	g.P("func Create", r.service.GoName, "Validations() ",
		g.QualifiedGoIdent(r.file.GoImportPath.Ident(r.service.GoName+"Validations")),
		" {")
	g.P("return new(", util.LcFirst(r.service.GoName), "Validations)")
	g.P("}")

	return nil
}

func (r validationGenerator) genMethods(g *protogen.GeneratedFile) error {
	g.P()
	var methodIndex, streamIndex int

	for _, method := range r.service.Methods {
		if !method.Desc.IsStreamingServer() && !method.Desc.IsStreamingClient() {
			// Unary RPC method
			r.genServiceMethod(g, method, methodIndex)
			methodIndex++
		} else {
			// Streaming RPC method
			r.genServiceMethod(g, method, streamIndex)
			streamIndex++
		}
	}

	return nil
}

func (r validationGenerator) genServiceMethod(g *protogen.GeneratedFile, method *protogen.Method, index int) {
	g.P("func (c *", util.LcFirst(r.service.GoName), "Validations) ",
		validationMethodSignature(g, method), "{")
	g.P("return nil")
	g.P("}")
	g.P()
	return
}
