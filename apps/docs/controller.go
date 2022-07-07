package docs

import (
	"context"
	"errors"

	docsv1 "github.com/airdb/xadmin-api/genproto/docs/v1"
	"github.com/airdb/xadmin-api/output"
	"github.com/go-masonry/mortar/interfaces/cfg"
	"github.com/go-masonry/mortar/interfaces/log"
	"go.uber.org/fx"
	"google.golang.org/genproto/googleapis/api/httpbody"
)

// Controller responsible for the business logic of our DocsService
type Controller interface {
	docsv1.ServiceServer
}

type controllerDeps struct {
	fx.In

	Config cfg.Config
	Logger log.Logger
}

type controller struct {
	docsv1.UnimplementedServiceServer

	deps controllerDeps
	log  log.Fields
}

// CreateController is a constructor for Fx
func CreateController(deps controllerDeps) Controller {
	return &controller{
		deps: deps,
		log:  deps.Logger.WithField("controller", "docs"),
	}
}

func (c *controller) GetSwagger(ctx context.Context, request *docsv1.GetSwaggerRequest) (*httpbody.HttpBody, error) {
	c.log.Debug(ctx, "get swagger accepted")

	bytes, err := output.FS.ReadFile("apis.swagger.yaml")
	if err != nil {
		return nil, errors.New("swagger file not exist")
	}

	return &httpbody.HttpBody{
		ContentType: "text/plain",
		Data:        bytes,
	}, nil
}
