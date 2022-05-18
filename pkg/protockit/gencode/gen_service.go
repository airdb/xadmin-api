package gencode

import (
	"fmt"
	"strings"

	"github.com/airdb/xadmin-api/pkg/protockit/util"
	"google.golang.org/protobuf/compiler/protogen"
)

type serviceGenerator struct {
	gen     *protogen.Plugin
	file    *protogen.File
	service *protogen.Service
	name    string
}

func NewServiceGenerator(
	gen *protogen.Plugin,
	file *protogen.File,
	service *protogen.Service,
) *serviceGenerator {
	return &serviceGenerator{
		gen:     gen,
		file:    file,
		service: service,
	}
}

func (r *serviceGenerator) Run() error {
	if r.service.GoName[len(r.service.GoName)-7:] != "Service" {
		return fmt.Errorf("%s should end with Service", r.service.GoName)
	}

	r.name = strings.ToLower(r.service.GoName[0 : len(r.service.GoName)-7])

	if len(r.service.Methods) == 0 {
		return nil
	}

	filename := r.file.GeneratedFilenamePrefix + `.go`

	g := r.gen.NewGeneratedFile(filename, servicesPackage)

	g.P("package services")

	if err := r.genDeps(g); err != nil {
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

func (r serviceGenerator) genDeps(g *protogen.GeneratedFile) error {
	g.P()
	g.P("type ", util.LcFirst(r.service.GoName), "Deps struct {")
	g.P(fxPackage.Ident("In"))
	g.P()
	g.P("Config ", logPackage.Ident("Config"))
	g.P("Logger ", logPackage.Ident("Logger"))
	g.P("Metrics ", monitorPackage.Ident("Metrics"))
	g.P("}")

	return nil
}

func (r serviceGenerator) genImpl(g *protogen.GeneratedFile) error {
	g.P()
	g.P("type ", r.name, "Impl struct {")
	g.P("*", g.QualifiedGoIdent(
		r.file.GoImportPath.Ident("Unimplemented"+r.service.GoName+"Server")))
	g.P()
	g.P("deps ", util.LcFirst(r.service.GoName)+"Deps")
	g.P("log ", logPackage.Ident("Fields"))
	g.P("}")

	return nil
}

func (r serviceGenerator) genNew(g *protogen.GeneratedFile) error {
	g.P()
	g.P("func Create", r.service.GoName, "Server(deps ", util.LcFirst(r.service.GoName), "Deps) ",
		g.QualifiedGoIdent(r.file.GoImportPath.Ident(r.service.GoName+"Server")),
		" {")
	g.P("return &", r.name, "Impl{")
	g.P("deps: deps,")
	g.P(`log: deps.Logger.WithField("service", "`, r.name, `"),`)
	g.P("}")
	g.P("}")

	return nil
}

func (r serviceGenerator) genMethods(g *protogen.GeneratedFile) error {
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

func (r serviceGenerator) genServiceMethod(g *protogen.GeneratedFile, method *protogen.Method, index int) {
	g.P("func (c *", r.name, "Impl) ", methodSignature(g, method), "{")
	g.P("out := new(", method.Output.GoIdent, ")")
	g.P("return out, nil")
	g.P("}")
	g.P()
	return
}
