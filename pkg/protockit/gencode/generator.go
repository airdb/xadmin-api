package gencode

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/samber/lo"
	"google.golang.org/protobuf/compiler/protogen"
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
	dataPackage        protogen.GoImportPath
	validationsPackage protogen.GoImportPath
	controllersPackage protogen.GoImportPath
	servicesPackage    protogen.GoImportPath
	queryKitPackage    protogen.GoImportPath
)

var (
	actionsWords = []string{"List", "Get", "Create", "Update", "Delete"}
	reMethod     = regexp.MustCompile(
		fmt.Sprintf(`^(%s)([A-Z].*)$`, strings.Join(actionsWords, "|")),
	)
	reInOut = regexp.MustCompile(`^([A-Z].*)(Request|Response)$`)
)

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

func guessEntity(service *protogen.Service) map[string]*protogen.Message {
	entities := map[string]*protogen.Message{}
	for _, method := range service.Methods {
		ret := reMethod.FindStringSubmatch(method.GoName)
		if ret == nil || len(ret) < 1 {
			continue
		}
		if !lo.Contains([]string{"Get"}, ret[1]) {
			continue
		}

		entities[method.Output.GoIdent.GoName] = method.Output
	}

	return entities
}
