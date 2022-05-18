package gencode

import (
	"context"
	"strings"

	"github.com/airdb/xadmin-api/pkg/protockit/util"
	"google.golang.org/protobuf/compiler/protogen"
)

func Process(ctx context.Context, file *protogen.File) (context.Context, error) {
	gen := util.FromContextGen(ctx)
	if gen == nil {
		panic("plugin not in context")
	}

	repoPackage = protogen.GoImportPath(
		guessGoimportPath(file.GoImportPath.String(), "pkg", "interfaces", "repo"))
	repoKitPackage = protogen.GoImportPath(
		guessGoimportPath(file.GoImportPath.String(), "pkg", "repoKit"))
	dataPackage = protogen.GoImportPath(
		guessGoimportPath(file.GoImportPath.String(), "app", "data"))
	validationsPackage = protogen.GoImportPath(
		guessGoimportPath(file.GoImportPath.String(), "app", "validations"))
	controllersPackage = protogen.GoImportPath(
		guessGoimportPath(file.GoImportPath.String(), "app", "controllers"))
	servicesPackage = protogen.GoImportPath(
		guessGoimportPath(file.GoImportPath.String(), "app", "services"))
	queryKitPackage = protogen.GoImportPath(
		guessGoimportPath(file.GoImportPath.String(), "pkg", "querykit"))

	var rangeErr error

	for _, service := range file.Services {
		if err := NewDataGenerator(gen, file, service).Run(); err != nil {
			rangeErr = err
		}
		if err := NewValidationGenerator(gen, file, service).Run(); err != nil {
			rangeErr = err
		}
		if err := NewControllerGenerator(gen, file, service).Run(); err != nil {
			rangeErr = err
		}
		if err := NewConverterGenerator(gen, file, service).Run(); err != nil {
			rangeErr = err
		}
		if err := NewServiceGenerator(gen, file, service).Run(); err != nil {
			rangeErr = err
		}
	}
	if rangeErr != nil {
		return ctx, rangeErr
	}
	return ctx, nil
}

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
