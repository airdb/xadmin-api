package services

import (
	"context"

	"github.com/airdb/xadmin-api/app/controllers"
	"github.com/airdb/xadmin-api/app/validations"
	apiv1 "github.com/airdb/xadmin-api/genproto/v1"
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
	apiv1.UnimplementedPassportServiceServer // if keep this one added even when you change your interface this code will compile

	deps passportServiceDeps
}

func CreatePassportServiceService(deps passportServiceDeps) apiv1.PassportServiceServer {
	return &passportImpl{
		deps: deps,
	}
}

func (w *passportImpl) LoginUrl(ctx context.Context, request *apiv1.LoginUrlRequest) (*apiv1.LoginUrlResponse, error) {
	w.deps.Logger.WithField("request", request).Debug(ctx, "hello request")
	return w.deps.Controller.LoginUrl(ctx, request)
}

func (w *passportImpl) Login(ctx context.Context, request *apiv1.LoginRequest) (*apiv1.LoginResponse, error) {
	if err := w.deps.Validations.Login(ctx, request); err != nil {
		return nil, err
	}
	w.deps.Logger.WithField("request", request).Debug(ctx, "login request")
	return w.deps.Controller.Login(ctx, request)
}

func (w *passportImpl) LoginCallback(ctx context.Context, request *apiv1.LoginCallbackRequest) (*apiv1.LoginCallbackResponse, error) {
	if err := w.deps.Validations.LoginCallback(ctx, request); err != nil {
		return nil, err
	}
	w.deps.Logger.WithField("request", request).Debug(ctx, "login callback request")
	return w.deps.Controller.LoginCallback(ctx, request)
}

func (w *passportImpl) Profile(ctx context.Context, request *apiv1.ProfileRequest) (result *apiv1.ProfileResponse, err error) {
	err = w.deps.Validations.Profile(ctx, request)
	if err == nil {
		result, err = w.deps.Controller.Profile(ctx, request)
	}
	w.deps.Logger.WithError(err).Debug(ctx, "logout request")
	return
}

func (w *passportImpl) Logout(ctx context.Context, request *apiv1.LogoutRequest) (result *empty.Empty, err error) {
	err = w.deps.Validations.Logout(ctx, request)
	if err == nil {
		result, err = w.deps.Controller.Logout(ctx, request)
	}
	w.deps.Logger.WithError(err).Debug(ctx, "logout request")
	return
}
