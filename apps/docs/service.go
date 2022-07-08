package docs

import (
	"context"

	docsv1 "github.com/airdb/xadmin-api/genproto/docs/v1"
	"github.com/go-masonry/mortar/interfaces/cfg"
	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/go-masonry/mortar/interfaces/monitor"
	"go.uber.org/fx"
	"google.golang.org/genproto/googleapis/api/httpbody"
)

type serviceDeps struct {
	fx.In

	Logger     log.Logger
	Config     cfg.Config
	Controller Controller
	Metrics    monitor.Metrics `optional:"true"`
}

type serviceImpl struct {
	docsv1.UnimplementedServiceServer // if keep this one added even when you change your interface this code will compile

	deps serviceDeps
	log  log.Fields
}

func CreatePassportServiceService(deps serviceDeps) docsv1.ServiceServer {
	return &serviceImpl{
		deps: deps,
		log:  deps.Logger.WithField("service", "docs"),
	}
}

func (w *serviceImpl) GetSwagger(ctx context.Context, request *docsv1.GetSwaggerRequest) (*httpbody.HttpBody, error) {
	w.log.WithField("request", request).Debug(ctx, "hello request")
	return w.deps.Controller.GetSwagger(ctx, request)
}
