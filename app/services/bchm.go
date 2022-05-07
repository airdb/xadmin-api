package services

import (
	"context"

	"github.com/airdb/xadmin-api/app/controllers"
	"github.com/airdb/xadmin-api/app/validations"
	bchmv1 "github.com/airdb/xadmin-api/genproto/bchm/v1"
	"github.com/go-masonry/mortar/interfaces/cfg"
	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/go-masonry/mortar/interfaces/monitor"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/fx"
)

type bchmServiceDeps struct {
	fx.In

	Logger      log.Logger
	Config      cfg.Config
	Controller  controllers.BchmServiceController
	Validations validations.BchmServiceValidations
	Metrics     monitor.Metrics `optional:"true"`
}

type bchmImpl struct {
	bchmv1.UnimplementedBchmServiceServer // if keep this one added even when you change your interface this code will compile

	deps bchmServiceDeps
	log  log.Fields
}

func CreateBchmServiceService(deps bchmServiceDeps) bchmv1.BchmServiceServer {
	return &bchmImpl{
		deps: deps,
		log:  deps.Logger.WithField("service", "bchm"),
	}
}

func (w *bchmImpl) ListLosts(ctx context.Context, request *bchmv1.ListLostsRequest) (*bchmv1.ListLostsResponse, error) {
	w.log.WithField("request", request).Debug(ctx, "list lost request")
	return w.deps.Controller.ListLosts(ctx, request)
}

func (w *bchmImpl) GetLost(ctx context.Context, request *bchmv1.GetLostRequest) (*bchmv1.GetLostResponse, error) {
	w.log.WithField("request", request).Debug(ctx, "get lost request")
	if err := w.deps.Validations.GetLost(ctx, request); err != nil {
		return nil, err
	}
	return w.deps.Controller.GetLost(ctx, request)
}

func (w *bchmImpl) CreateLost(ctx context.Context, request *bchmv1.CreateLostRequest) (*bchmv1.CreateLostResponse, error) {
	w.log.WithField("request", request).Debug(ctx, "create lost request")
	if err := w.deps.Validations.CreateLost(ctx, request); err != nil {
		return nil, err
	}
	return w.deps.Controller.CreateLost(ctx, request)
}

func (w *bchmImpl) UpdateLost(ctx context.Context, request *bchmv1.UpdateLostRequest) (result *bchmv1.UpdateLostResponse, err error) {
	w.log.WithField("request", request).Debug(ctx, "update lost request")
	err = w.deps.Validations.UpdateLost(ctx, request)
	if err == nil {
		result, err = w.deps.Controller.UpdateLost(ctx, request)
	}
	w.log.WithError(err).Debug(ctx, "update lost done")
	return
}

func (w *bchmImpl) DeleteLost(ctx context.Context, request *bchmv1.DeleteLostRequest) (result *empty.Empty, err error) {
	w.log.WithField("request", request).Debug(ctx, "delete lost request")
	err = w.deps.Validations.DeleteLost(ctx, request)
	if err == nil {
		result, err = w.deps.Controller.DeleteLost(ctx, request)
	}
	w.log.WithError(err).Debug(ctx, "delete lost done")
	return
}
