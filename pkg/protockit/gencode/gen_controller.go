package gencode

import (
	"fmt"
	"path"
	"strings"

	"github.com/airdb/xadmin-api/pkg/protockit/util"
	"github.com/gobeam/stringy"
	"github.com/samber/lo"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type controllerGenerator struct {
	gen     *protogen.Plugin
	file    *protogen.File
	service *protogen.Service
	name    string
}

func NewControllerGenerator(
	gen *protogen.Plugin,
	file *protogen.File,
	service *protogen.Service,
) *controllerGenerator {
	return &controllerGenerator{
		gen:     gen,
		file:    file,
		service: service,
	}
}

func (r *controllerGenerator) Run() error {
	if r.service.GoName[len(r.service.GoName)-7:] != "Service" {
		return fmt.Errorf("%s should end with Service", r.service.GoName)
	}

	r.name = strings.ToLower(r.service.GoName[0 : len(r.service.GoName)-7])

	if len(r.service.Methods) == 0 {
		return nil
	}

	filename := path.Join("app", "controllers",
		stringy.New(r.name).SnakeCase().LcFirst()+`.go`)

	g := r.gen.NewGeneratedFile(filename, controllersPackage)

	g.P("package controllers")

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

	if err := r.genMethods(g); err != nil {
		return err
	}

	return nil
}

func (r controllerGenerator) genInterface(g *protogen.GeneratedFile) error {
	g.P()
	g.P("type ", r.service.GoName, "Controller interface {")
	g.P(g.QualifiedGoIdent(r.file.GoImportPath.Ident(r.service.GoName + "Server")))
	g.P("}")

	return nil
}

func (r controllerGenerator) genDeps(g *protogen.GeneratedFile) error {
	entityMethods := map[string][]string{}
	for _, method := range r.service.Methods {
		ret := reMethod.FindStringSubmatch(method.GoName)
		if ret == nil || len(ret) < 1 {
			continue
		}
		name := util.Pluralize.Singular(ret[2])
		if _, ok := entityMethods[name]; ok {
			entityMethods[name] = []string{}
		}
		entityMethods[name] = append(entityMethods[name], name)
	}
	g.P()
	g.P("type ", util.LcFirst(r.service.GoName), "ControllerDeps struct {")
	g.P(fxPackage.Ident("In"))
	g.P()
	g.P("Config ", cfgPackage.Ident("Config"))
	g.P("Logger ", logPackage.Ident("Logger"))
	g.P("Metrics ", monitorPackage.Ident("Metrics"))
	g.P()
	for entityMethod := range entityMethods {
		g.P(entityMethod+"Repo ", dataPackage.Ident(entityMethod+"Repo"))
	}
	g.P("}")

	return nil
}

func (r controllerGenerator) genImpl(g *protogen.GeneratedFile) error {
	g.P()
	g.P("type ", util.LcFirst(r.service.GoName), "Controller struct {")
	g.P(g.QualifiedGoIdent(
		r.file.GoImportPath.Ident("Unimplemented" + r.service.GoName + "Server")))
	g.P()
	g.P("deps ", util.LcFirst(r.service.GoName)+"ControllerDeps")
	g.P("log ", logPackage.Ident("Fields"))
	g.P("convert *", r.name+"Convert")
	g.P("}")

	return nil
}

func (r controllerGenerator) genNew(g *protogen.GeneratedFile) error {
	g.P()
	g.P("func New", r.service.GoName, "Controller(deps ",
		util.LcFirst(r.service.GoName), "ControllerDeps) ",
		r.service.GoName+"Controller",
		" {")
	g.P("return &", util.LcFirst(r.service.GoName), "Controller{")
	g.P("deps: deps,")
	g.P(`log: deps.Logger.WithField("service", "`, r.name, `"),`)
	g.P(`convert: new`, util.UcFirst(r.name), "Convert(),")
	g.P("}")
	g.P("}")

	return nil
}

func (r controllerGenerator) genMethods(g *protogen.GeneratedFile) error {
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

func (r controllerGenerator) genServiceMethod(g *protogen.GeneratedFile, method *protogen.Method, index int) {
	ret := reMethod.FindStringSubmatch(method.GoName)
	if ret == nil || len(ret) < 1 {
		r.genServiceMethodDefault(g, method, index)
		return
	}

	switch ret[1] {
	case "List":
		r.genServiceMethodList(g, method, index)
	case "Get":
		r.genServiceMethodGet(g, method, index)
	case "Create":
		r.genServiceMethodCreate(g, method, index)
	case "Update":
		r.genServiceMethodUpdate(g, method, index)
	case "Delete":
		r.genServiceMethodDelete(g, method, index)
	default:
		r.genServiceMethodDefault(g, method, index)
	}

	return
}
func (r controllerGenerator) genServiceMethodDefault(g *protogen.GeneratedFile, method *protogen.Method, index int) {
	g.P("func (c *", util.LcFirst(r.service.GoName), "Controller) ",
		methodSignature(g, method), "{")
	g.P("out := new(", method.Output.GoIdent, ")")
	g.P("return out, nil")
	g.P("}")
	g.P()
}

func (r controllerGenerator) genServiceMethodList(g *protogen.GeneratedFile, method *protogen.Method, index int) {
	ret := reMethod.FindStringSubmatch(method.GoName)
	entityStr := util.Pluralize.Singular(ret[2])
	repoStr := entityStr + "Repo"

	g.P("func (c *", util.LcFirst(r.service.GoName), "Controller) ",
		methodSignature(g, method), "{")
	g.P(`c.log.Debug(ctx, "`, method.GoName, ` accepted")`)
	g.P()
	g.P("total, filtered, err := c.deps.", repoStr+".Count(ctx, in)")
	g.P("if err != nil {")
	errorListStr := `"` + method.GoName + ` count error"`
	g.P("c.log.WithError(err).Debug(ctx, ", errorListStr, ")")
	g.P("return nil, ", errorsPackage.Ident("New"), "(", errorListStr, ")")
	g.P("}")
	g.P("if total == 0 {")
	errorEmptyStr := `"` + r.name + "s is empty" + `"`
	g.P("return nil, ", errorsPackage.Ident("New"), "(", errorEmptyStr, ")")
	g.P("}")
	g.P("")
	g.P("items, err := c.deps.", repoStr, ".List(ctx, in)")
	g.P("if err != nil {")
	errorListQuery := `"` + method.GoName + ` error"`
	g.P("c.log.WithError(err).Debug(ctx, ", errorListQuery, ")")
	g.P("return nil, ", errorsPackage.Ident("New"), "(", errorListQuery, ")")
	g.P("}")
	g.P("")
	r.genOut(g, method.Output, index)
	g.P("}")
	g.P()
	return
}

func (r controllerGenerator) genServiceMethodGet(g *protogen.GeneratedFile, method *protogen.Method, index int) {
	ret := reMethod.FindStringSubmatch(method.GoName)
	entityStr := util.Pluralize.Singular(ret[2])
	repoStr := entityStr + "Repo"

	g.P("func (c *", util.LcFirst(r.service.GoName), "Controller) ",
		methodSignature(g, method), "{")
	g.P(`c.log.Debug(ctx, "`, method.GoName, ` accepted")`)
	g.P()
	g.P("item, err := c.deps.", repoStr, ".Get(ctx, uint(in.GetId()))")
	g.P("if err != nil {")
	errorListStr := `"` + method.GoName + ` error"`
	g.P("c.log.WithError(err).Debug(ctx, ", errorListStr, ")")
	g.P("return nil, ", errorsPackage.Ident("New"), "(", errorListStr, ")")
	g.P("}")
	g.P("")
	r.genOut(g, method.Output, index)
	g.P("}")
	g.P()
	return
}

func (r controllerGenerator) genServiceMethodCreate(g *protogen.GeneratedFile, method *protogen.Method, index int) {
	ret := reMethod.FindStringSubmatch(method.GoName)
	entityStr := util.Pluralize.Singular(ret[2])
	repoStr := entityStr + "Repo"

	g.P("func (c *", util.LcFirst(r.service.GoName), "Controller) ",
		methodSignature(g, method), "{")
	g.P(`c.log.Debug(ctx, "`, method.GoName, ` accepted")`)
	g.P()
	entityName := util.Pluralize.Singular(entityStr)
	convertFuncName := fmt.Sprintf("FromProtoCreate%sToModel%s", entityName, entityName)
	g.P("item := c.convert.", convertFuncName, "(in)")
	g.P("err := c.deps.", repoStr, ".Create(ctx, item)")
	g.P("if err != nil {")
	errorListStr := `"` + method.GoName + ` error"`
	g.P("c.log.WithError(err).Debug(ctx, ", errorListStr, ")")
	g.P("return nil, ", errorsPackage.Ident("New"), "(", errorListStr, ")")
	g.P("}")
	g.P("")
	r.genOut(g, method.Output, index)
	g.P("}")
	g.P()
	return
}

func (r controllerGenerator) genServiceMethodUpdate(g *protogen.GeneratedFile, method *protogen.Method, index int) {
	ret := reMethod.FindStringSubmatch(method.GoName)
	entityStr := util.Pluralize.Singular(ret[2])
	repoStr := entityStr + "Repo"

	g.P("func (c *", util.LcFirst(r.service.GoName), "Controller) ",
		methodSignature(g, method), "{")
	g.P(`c.log.Debug(ctx, "`, method.GoName, ` accepted")`)
	g.P()
	entityName := util.Pluralize.Singular(entityStr)
	convertFuncName := fmt.Sprintf("FromProto%sToModel%s", entityName, entityName)
	g.P("data := c.convert.", convertFuncName, "(in.GetItem())")
	g.P()
	g.P("fm := ", queryKitPackage.Ident("NewField"),
		`(in.GetUpdateMask(), in.GetItem()).WithAction("update")`)
	g.P()
	g.P("err := c.deps.", repoStr, ".Update(ctx, data.ID, data, fm)")
	g.P("if err != nil {")
	errorStr := `"` + method.GoName + ` failed"`
	g.P("c.log.WithError(err).Debug(ctx, ", errorStr, ")")
	g.P("return nil, ", errorsPackage.Ident("New"), "(", errorStr, ")")
	g.P("}")
	g.P("")
	g.P("item, err := c.deps.", repoStr, ".Get(ctx, uint(data.ID))")
	g.P("if err != nil {")
	errorStr = `"` + method.GoName + ` item not exist"`
	g.P("c.log.WithError(err).Debug(ctx, ", errorStr, ")")
	g.P("return nil, ", errorsPackage.Ident("New"), "(", errorStr, ")")
	g.P("}")
	g.P("")
	r.genOut(g, method.Output, index)
	g.P("}")
	g.P()
	return
}

func (r controllerGenerator) genServiceMethodDelete(g *protogen.GeneratedFile, method *protogen.Method, index int) {
	ret := reMethod.FindStringSubmatch(method.GoName)
	entityStr := util.Pluralize.Singular(ret[2])
	repoStr := entityStr + "Repo"

	g.P("func (c *", util.LcFirst(r.service.GoName), "Controller) ",
		methodSignature(g, method), "{")
	g.P(`c.log.Debug(ctx, "`, method.GoName, ` accepted")`)
	g.P()
	g.P("err := c.deps.", repoStr, ".Delete(ctx, uint(in.GetId()))")
	g.P("if err != nil {")
	errorListStr := `"` + method.GoName + ` failed"`
	g.P("c.log.WithError(err).Debug(ctx, ", errorListStr, ")")
	g.P("return nil, ", errorsPackage.Ident("New"), "(", errorListStr, ")")
	g.P("}")
	g.P("")
	r.genOut(g, method.Output, index)
	g.P("}")
	g.P()
	return
}

func (r controllerGenerator) genOut(g *protogen.GeneratedFile, msg *protogen.Message, index int) {
	ret := reInOut.FindStringSubmatch(msg.GoIdent.GoName)
	if ret == nil {
		if strings.Contains(msg.GoIdent.GoImportPath.String(), "genproto") {
			entityName := util.Pluralize.Singular(msg.GoIdent.GoName)
			convertFuncName := fmt.Sprintf("FromModel%sToProto%s", entityName, entityName)
			g.P("return ", "c.convert.", convertFuncName, "(item), nil")
			return
		}
		g.P("return &", g.QualifiedGoIdent(msg.GoIdent), "{}, nil")
		return
	}

	g.P("return &", g.QualifiedGoIdent(msg.GoIdent), "{")
	intKinds := []protoreflect.Kind{protoreflect.Int32Kind, protoreflect.Int64Kind}
	for _, field := range msg.Fields {
		if field.Desc.IsList() {
			entityName := util.Pluralize.Plural(field.Message.GoIdent.GoName)
			convertFuncName := fmt.Sprintf("FromModel%sToProto%s", entityName, entityName)
			g.P(field.GoName, ": c.convert.", convertFuncName, "(items),")
		} else if lo.Contains(intKinds, field.Desc.Kind()) {
			switch field.GoName {
			case "TotalSize":
				g.P(field.GoName, ": total,")
				continue
			case "FilteredSize":
				g.P(field.GoName, ": filtered,")
				continue
			}
		} else if field.Desc.Kind() == protoreflect.MessageKind {
			entityName := util.Pluralize.Singular(field.Message.GoIdent.GoName)
			convertFuncName := fmt.Sprintf("FromModel%sToProto%s", entityName, entityName)
			g.P(field.GoName, ": c.convert.", convertFuncName, "(item),")
			continue
		}
	}
	g.P("}, nil")
}
