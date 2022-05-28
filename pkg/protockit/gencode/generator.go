package gencode

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	annov1 "github.com/airdb/xadmin-api/genproto/annotation/v1"
	"github.com/airdb/xadmin-api/pkg/protockit/util"
	"github.com/samber/lo"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	contextPackage = protogen.GoImportPath("context")
	errorsPackage  = protogen.GoImportPath("errors")
	fxPackage      = protogen.GoImportPath("go.uber.org/fx")
	cfgPackage     = protogen.GoImportPath("github.com/go-masonry/mortar/interfaces/cfg")
	logPackage     = protogen.GoImportPath("github.com/go-masonry/mortar/interfaces/log")
	monitorPackage = protogen.GoImportPath("github.com/go-masonry/mortar/interfaces/monitor")
)

var (
	repoPackage        protogen.GoImportPath
	repoKitPackage     protogen.GoImportPath
	queryKitPackage    protogen.GoImportPath
	dataPackage        protogen.GoImportPath
	validationsPackage protogen.GoImportPath
	controllersPackage protogen.GoImportPath
	servicesPackage    protogen.GoImportPath
)

var (
	actionsWords = []string{"List", "Get", "Create", "Update", "Delete"}
	reMethod     = regexp.MustCompile(
		fmt.Sprintf(`^(%s)([A-Z].*)$`, strings.Join(actionsWords, "|")),
	)
	reInOut = regexp.MustCompile(`^([A-Z].*)(Request|Response)$`)
)

func guessGoimportPath(s string, dirs ...string) protogen.GoImportPath {
	segs := []string{}
	for _, seg := range strings.Split(strings.Trim(s, "\""), "/") {
		if seg == "genproto" {
			break
		}
		segs = append(segs, seg)
	}

	segs = append(segs, dirs...)

	return protogen.GoImportPath(strings.Join(segs, "/"))
}

func methodSignature(g *protogen.GeneratedFile, method *protogen.Method) string {
	s := method.GoName + "(ctx " + g.QualifiedGoIdent(contextPackage.Ident("Context"))
	s += ", in *" + g.QualifiedGoIdent(method.Input.GoIdent)
	s += ") ("
	s += "*" + g.QualifiedGoIdent(method.Output.GoIdent)
	s += ", error)"
	return s
}

func validationMethodSignature(g *protogen.GeneratedFile, method *protogen.Method) string {
	s := method.GoName + "(ctx " + g.QualifiedGoIdent(contextPackage.Ident("Context"))
	s += ", in *" + g.QualifiedGoIdent(method.Input.GoIdent)
	s += ") error"
	return s
}

func guessEntity(service *protogen.Service) []*protogen.Message {
	entities := []*protogen.Message{}
	for _, method := range service.Methods {
		retMethod := reMethod.FindStringSubmatch(method.GoName)
		if retMethod == nil || len(retMethod) < 1 {
			continue
		}
		if !lo.Contains([]string{"Get"}, retMethod[1]) {
			continue
		}

		retOutput := reInOut.FindStringSubmatch(method.Output.GoIdent.GoName)
		if retOutput == nil {
			entities = append(entities, method.Output)
			continue
		}
		if len(retOutput) != 3 {
			continue
		}
		for _, field := range method.Output.Fields {
			option, err := util.KitParser[annov1.FieldDescriptor](
				service.Comments.Leading.String())
			if err != nil {
				log.Printf("can not parse %s kit option: (%s)", service.GoName, err)
			}
			if util.KitGencodeLayerEmpty(&option) {
				continue
			}
			if field.Desc.Kind() == protoreflect.MessageKind {
				entities = append(entities, field.Message)
			}
		}
	}

	return entities
}

func findCreateRequest(service *protogen.Service) *protogen.Message {
	for _, method := range service.Methods {
		retMethod := reMethod.FindStringSubmatch(method.GoName)
		if retMethod == nil || len(retMethod) < 1 {
			continue
		}
		if !lo.Contains([]string{"Create"}, retMethod[1]) {
			continue
		}

		return method.Input
	}

	return nil
}
