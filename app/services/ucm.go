package services

import (
	"github.com/airdb/xadmin-api/app/controllers"
	"github.com/airdb/xadmin-api/app/validations"
	ucmv1 "github.com/airdb/xadmin-api/genproto/ucm/v1"
	"github.com/go-masonry/mortar/interfaces/cfg"
	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/go-masonry/mortar/interfaces/monitor"
	"go.uber.org/fx"
)

type ucmServiceDeps struct {
	fx.In

	Logger      log.Logger
	Config      cfg.Config
	Controller  controllers.UcmServiceController
	Validations validations.UcmServiceValidations
	Metrics     monitor.Metrics `optional:"true"`
}

type ucmImpl struct {
	ucmv1.UnimplementedServiceServer // if keep this one added even when you change your interface this code will compile

	deps ucmServiceDeps
	log  log.Fields
}

func CreateServiceService(deps ucmServiceDeps) ucmv1.ServiceServer {
	return &ucmImpl{
		deps: deps,
		log:  deps.Logger.WithField("service", "ucm"),
	}
}
