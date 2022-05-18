package gencode

import (
	"fmt"
	"strings"

	"github.com/airdb/xadmin-api/pkg/protockit/util"
	"github.com/gobeam/stringy"
	"google.golang.org/protobuf/compiler/protogen"
)

type dataGenerator struct {
	gen     *protogen.Plugin
	file    *protogen.File
	service *protogen.Service
	name    string
}

func NewDataGenerator(
	gen *protogen.Plugin,
	file *protogen.File,
	service *protogen.Service,
) *dataGenerator {
	return &dataGenerator{
		gen:     gen,
		file:    file,
		service: service,
	}
}

func (r *dataGenerator) Run() error {
	if r.service.GoName[len(r.service.GoName)-7:] != "Service" {
		return fmt.Errorf("%s should end with Service", r.service.GoName)
	}

	r.name = strings.ToLower(r.service.GoName[0 : len(r.service.GoName)-7])

	if len(r.service.Methods) == 0 {
		return nil
	}

	for _, entity := range guessEntity(r.service) {
		filename := stringy.New(entity.GoIdent.GoName).SnakeCase().LcFirst() + `.go`
		g := r.gen.NewGeneratedFile(filename, dataPackage)

		g.P("package data")

		if err := r.genInterface(g, entity); err != nil {
			return err
		}

		if err := r.genDeps(g, entity); err != nil {
			return err
		}

		if err := r.genImpl(g, entity); err != nil {
			return err
		}

		if err := r.genNew(g, entity); err != nil {
			return err
		}
	}

	return nil
}

func (r dataGenerator) genInterface(g *protogen.GeneratedFile, message *protogen.Message) error {
	g.P()
	g.P("type ", message.GoIdent.GoName, "Repo interface {")
	g.P(g.QualifiedGoIdent(repoPackage.Ident("Repo")),
		r.genericSignature(g, message))
	g.P("}")

	return nil
}

func (r dataGenerator) genDeps(g *protogen.GeneratedFile, message *protogen.Message) error {
	g.P()
	g.P("type ", util.LcFirst(message.GoIdent.GoName), "RepoDeps struct {")
	g.P(fxPackage.Ident("In"))
	g.P("}")

	return nil
}

func (r dataGenerator) genImpl(g *protogen.GeneratedFile, message *protogen.Message) error {
	g.P()
	g.P("type ", util.LcFirst(message.GoIdent.GoName), "Repo struct {")
	g.P("*", g.QualifiedGoIdent(repoKitPackage.Ident("Repo")),
		r.genericSignature(g, message))
	g.P()
	g.P("deps ", util.LcFirst(message.GoIdent.GoName), "RepoDeps")
	g.P("}")

	return nil
}

func (r dataGenerator) genNew(g *protogen.GeneratedFile, message *protogen.Message) error {
	g.P()
	g.P("func New", message.GoIdent.GoName,
		"Repo(deps ", util.LcFirst(message.GoIdent.GoName), "RepoDeps) ",
		message.GoIdent.GoName+"Repo {")
	g.P("repo := ", util.LcFirst(message.GoIdent.GoName), "Repo {")
	g.P("deps: deps,")
	g.P("}")
	g.P("repo.Repo = ",
		repoKitPackage.Ident("NewRepo"), r.genericSignature(g, message))
	g.P()
	g.P("return repo")
	g.P("}")

	return nil
}

func (r dataGenerator) genericSignature(g *protogen.GeneratedFile, message *protogen.Message) string {
	s := "["
	s += message.GoIdent.GoName
	s += "Entity, uint, "
	s += g.QualifiedGoIdent(r.file.GoImportPath.Ident(
		"List" + message.GoIdent.GoName + "Request"))
	s += "]"

	return s
}
