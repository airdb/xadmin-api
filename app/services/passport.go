package services

import (
	"context"

	"github.com/airdb/xadmin-api/app/controllers"
	"github.com/airdb/xadmin-api/app/validations"
	passportv1 "github.com/airdb/xadmin-api/genproto/passport/v1"
	"github.com/go-masonry/mortar/interfaces/cfg"
	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/go-masonry/mortar/interfaces/monitor"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/fx"
)

type passportServiceDeps struct {
	fx.In

	Logger      log.Logger
	Config      cfg.Config
	Controller  controllers.PassportServiceController
	Validations validations.PassportServiceValidations
	Metrics     monitor.Metrics `optional:"true"`
}

type passportImpl struct {
	passportv1.UnimplementedPassportServiceServer // if keep this one added even when you change your interface this code will compile

	deps passportServiceDeps
	log  log.Fields
}

func CreatePassportServiceService(deps passportServiceDeps) passportv1.PassportServiceServer {
	return &passportImpl{
		deps: deps,
		log:  deps.Logger.WithField("service", "passport"),
	}
}

func (w *passportImpl) Preset(ctx context.Context, request *passportv1.PresetRequest) (*passportv1.PresetResponse, error) {
	w.log.WithField("request", request).Debug(ctx, "hello request")
	return w.deps.Controller.Preset(ctx, request)
}

func (w *passportImpl) Login(ctx context.Context, request *passportv1.LoginRequest) (*passportv1.LoginResponse, error) {
	if err := w.deps.Validations.Login(ctx, request); err != nil {
		return nil, err
	}
	w.log.WithField("request", request).Debug(ctx, "login request")
	return w.deps.Controller.Login(ctx, request)
}

func (w *passportImpl) Callback(ctx context.Context, request *passportv1.CallbackRequest) (*passportv1.CallbackResponse, error) {
	if err := w.deps.Validations.Callback(ctx, request); err != nil {
		return nil, err
	}
	w.log.WithField("request", request).Debug(ctx, "login callback request")
	return w.deps.Controller.Callback(ctx, request)
}

func (w *passportImpl) Profile(ctx context.Context, request *passportv1.ProfileRequest) (result *passportv1.ProfileResponse, err error) {
	err = w.deps.Validations.Profile(ctx, request)
	if err == nil {
		result, err = w.deps.Controller.Profile(ctx, request)
	}
	w.log.WithError(err).Debug(ctx, "logout request")
	return
}

func (w *passportImpl) Logout(ctx context.Context, request *passportv1.LogoutRequest) (result *empty.Empty, err error) {
	err = w.deps.Validations.Logout(ctx, request)
	if err == nil {
		result, err = w.deps.Controller.Logout(ctx, request)
	}
	w.log.WithError(err).Debug(ctx, "logout request")
	return
}
