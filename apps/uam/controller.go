package uam

import (
	"github.com/airdb/xadmin-api/apps/data"
	uamv1 "github.com/airdb/xadmin-api/genproto/uam/v1"
	"github.com/go-masonry/mortar/interfaces/log"
	"go.uber.org/fx"
)

// UamServiceController responsible for the business logic of our UamService
type UamServiceController interface {
	uamv1.ServiceServer
}

type uamControllerDeps struct {
	fx.In

	Logger   log.Logger
	UserRepo data.UserRepo
}

type uamController struct {
	uamv1.UnimplementedServiceServer

	log    log.Fields
	deps   uamControllerDeps
	conver *uamConvert
}

// CreateUamServiceController is a constructor for Fx
func CreateUamServiceController(deps uamControllerDeps) UamServiceController {
	return &uamController{
		log:    deps.Logger.WithField("controller", "uam"),
		deps:   deps,
		conver: newUamConvert(),
	}
}
