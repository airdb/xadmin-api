package bchm

import (
	bchmv1 "github.com/airdb/xadmin-api/genproto/bchm/v1"
	"github.com/go-masonry/mortar/interfaces/cfg"
	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/go-masonry/mortar/interfaces/monitor"
	"go.uber.org/fx"
)

type serviceDeps struct {
	fx.In

	Logger      log.Logger
	Config      cfg.Config
	Controller  Controller
	Validations ServiceValidations
	Metrics     monitor.Metrics `optional:"true"`
}

type serviceImpl struct {
	bchmv1.UnimplementedServiceServer // if keep this one added even when you change your interface this code will compile

	deps serviceDeps
	log  log.Fields
}

func CreateServiceService(deps serviceDeps) bchmv1.ServiceServer {
	return &serviceImpl{
		deps: deps,
		log:  deps.Logger.WithField("service", "bchm"),
	}
}
