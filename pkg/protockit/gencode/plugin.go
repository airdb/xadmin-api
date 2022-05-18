package gencode

import (
	"context"
	"fmt"
	"log"

	annov1 "github.com/airdb/xadmin-api/genproto/annotation/v1"
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
	queryKitPackage = protogen.GoImportPath(
		guessGoimportPath(file.GoImportPath.String(), "pkg", "querykit"))
	dataPackage = protogen.GoImportPath(
		guessGoimportPath(file.GoImportPath.String(), "app", "data"))
	validationsPackage = protogen.GoImportPath(
		guessGoimportPath(file.GoImportPath.String(), "app", "validations"))
	controllersPackage = protogen.GoImportPath(
		guessGoimportPath(file.GoImportPath.String(), "app", "controllers"))
	servicesPackage = protogen.GoImportPath(
		guessGoimportPath(file.GoImportPath.String(), "app", "services"))

	services := map[*protogen.Service]*annov1.MessageDescriptor{}
	entiryNames := []string{}
	entiryMessages := []*protogen.Message{}

	for _, service := range file.Services {
		option, err := util.KitParser[annov1.MessageDescriptor](
			service.Comments.Leading.String())
		if err != nil {
			log.Printf("can not parse %s kit option: (%s)", service.GoName, err)
		}
		if !util.KitGencodeValid(&option) {
			continue
		}
		if service.GoName[len(service.GoName)-7:] != "Service" {
			return ctx, fmt.Errorf("%s should end with Service", service.GoName)
		}

		services[service] = &option

		if util.KitGencodeLayerValid(&option, "*", "repo") {
			entities := guessEntity(service)
			entiryMessages = append(entiryMessages, entities...)
			for _, entity := range entities {
				option, err := util.KitParser[annov1.MessageDescriptor](
					entity.Comments.Leading.String())
				if err != nil {
					log.Printf("can not parse %s kit option: (%s)", entity.GoIdent.GoName, err)
				}
				if util.KitGencodeLayerEmpty(&option) {
					continue
				}
				entiryNames = append(entiryNames, entity.GoIdent.GoName)
			}
		}
	}

	var rangeErr error
	for _, entity := range entiryMessages {
		if err := NewRepoGenerator(gen, file, entity).Run(); err != nil {
			rangeErr = err
		}
	}

	for service := range services {
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
