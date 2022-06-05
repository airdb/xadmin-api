package services

import (
	"github.com/airdb/xadmin-api/app/controllers"
	"github.com/airdb/xadmin-api/app/validations"
	uamv1 "github.com/airdb/xadmin-api/genproto/uam/v1"
	"github.com/go-masonry/mortar/interfaces/cfg"
	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/go-masonry/mortar/interfaces/monitor"
	"go.uber.org/fx"
)

type uamServiceDeps struct {
	fx.In

	Logger      log.Logger
	Config      cfg.Config
	Controller  controllers.UamServiceController
	Validations validations.UamServiceValidations
	Metrics     monitor.Metrics `optional:"true"`
}

type uamImpl struct {
	uamv1.UnimplementedServiceServer // if keep this one added even when you change your interface this code will compile

	deps uamServiceDeps
	log  log.Fields
}

func CreateUamServiceService(deps uamServiceDeps) uamv1.ServiceServer {
	return &uamImpl{
		deps: deps,
		log:  deps.Logger.WithField("service", "uam"),
	}
}
