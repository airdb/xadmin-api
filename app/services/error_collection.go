package services

import (
	"context"

	"github.com/airdb/xadmin-api/app/controllers"
	"github.com/airdb/xadmin-api/app/validations"
	errorCollectionv1 "github.com/airdb/xadmin-api/genproto/error_collection/v1"
	"github.com/go-masonry/mortar/interfaces/cfg"
	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/go-masonry/mortar/interfaces/monitor"
	"go.uber.org/fx"
)

type errorCollectionServiceDeps struct {
	fx.In

	Logger      log.Logger
	Config      cfg.Config
	Controller  controllers.ErrorCollectionServiceController
	Validations validations.ErrorCollectionServiceValidations
	Metrics     monitor.Metrics `optional:"true"`
}

type errorCollectionImpl struct {
	errorCollectionv1.UnimplementedErrorCollectionServiceServer // if keep this one added even when you change your interface this code will compile

	deps errorCollectionServiceDeps
	log  log.Fields
}

func CreateErrorCollectionServiceService(deps errorCollectionServiceDeps) errorCollectionv1.ErrorCollectionServiceServer {
	return &errorCollectionImpl{
		deps: deps,
		log:  deps.Logger.WithField("service", "errorCollection"),
	}
}

func (w *errorCollectionImpl) Collect(ctx context.Context, request *errorCollectionv1.CreateErrorCollectionRequest) (*errorCollectionv1.CreateErrorCollectionResponse, error) {
	w.log.WithField("request", request).Debug(ctx, "collect request")

	return w.deps.Controller.CreateErrorCollection(ctx, request)
}
