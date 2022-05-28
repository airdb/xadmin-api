package gencode

import (
	"path"
	"strings"

	"github.com/airdb/xadmin-api/pkg/protockit/util"
	"github.com/gobeam/stringy"
	"google.golang.org/protobuf/compiler/protogen"
)

type repoGenerator struct {
	gen     *protogen.Plugin
	file    *protogen.File
	message *protogen.Message
	name    string
}

func NewRepoGenerator(
	gen *protogen.Plugin,
	file *protogen.File,
	message *protogen.Message,
) *repoGenerator {
	return &repoGenerator{
		gen:     gen,
		file:    file,
		message: message,
	}
}

func (r *repoGenerator) Run() error {
	r.name = strings.ToLower(r.message.GoIdent.GoName)

	filename := path.Join("app", "data",
		stringy.New(r.message.GoIdent.GoName).SnakeCase().LcFirst()+`.go`)
	g := r.gen.NewGeneratedFile(filename, dataPackage)

	g.P("package data")

	if err := r.genInterface(g); err != nil {
		return err
	}

	if err := r.genDeps(g); err != nil {
		return err
	}

	if err := r.genImpl(g); err != nil {
		return err
	}

	if err := r.genNew(g); err != nil {
		return err
	}

	return nil
}

func (r repoGenerator) genInterface(g *protogen.GeneratedFile) error {
	g.P()
	g.P("type ", r.message.GoIdent.GoName, "Repo interface {")
	g.P(g.QualifiedGoIdent(repoPackage.Ident("Repo")),
		r.genericSignature(g))
	g.P("}")

	return nil
}

func (r repoGenerator) genDeps(g *protogen.GeneratedFile) error {
	g.P()
	g.P("type ", util.LcFirst(r.message.GoIdent.GoName), "RepoDeps struct {")
	g.P(fxPackage.Ident("In"))
	g.P("}")

	return nil
}

func (r repoGenerator) genImpl(g *protogen.GeneratedFile) error {
	g.P()
	g.P("type ", util.LcFirst(r.message.GoIdent.GoName), "Repo struct {")
	g.P("*", g.QualifiedGoIdent(repoKitPackage.Ident("Repo")),
		r.genericSignature(g))
	g.P()
	g.P("deps ", util.LcFirst(r.message.GoIdent.GoName), "RepoDeps")
	g.P("}")

	return nil
}

func (r repoGenerator) genNew(g *protogen.GeneratedFile) error {
	g.P()
	g.P("func New", r.message.GoIdent.GoName,
		"Repo(deps ", util.LcFirst(r.message.GoIdent.GoName), "RepoDeps) ",
		r.message.GoIdent.GoName+"Repo {")
	g.P("repo := ", util.LcFirst(r.message.GoIdent.GoName), "Repo {")
	g.P("deps: deps,")
	g.P("}")
	g.P("repo.Repo = ",
		repoKitPackage.Ident("NewRepo"), r.genericSignature(g))
	g.P()
	g.P("return repo")
	g.P("}")

	return nil
}

func (r repoGenerator) genericSignature(g *protogen.GeneratedFile) string {
	s := "["
	s += g.QualifiedGoIdent(
		dataPackage.Ident(r.message.GoIdent.GoName + "Entity"))
	s += ", uint, "
	s += g.QualifiedGoIdent(r.file.GoImportPath.Ident(
		"List" + util.Pluralize.Plural(r.message.GoIdent.GoName) + "Request"))
	s += "]"

	return s
}
