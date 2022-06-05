package controllers

import (
	"github.com/airdb/xadmin-api/app/data"
	ucmv1 "github.com/airdb/xadmin-api/genproto/ucm/v1"
	"github.com/go-masonry/mortar/interfaces/log"
	"go.uber.org/fx"
)

// UcmServiceController responsible for the business logic of our UcmService
type UcmServiceController interface {
	ucmv1.ServiceServer
}

type ucmControllerDeps struct {
	fx.In

	Logger      log.Logger
	LostRepo    data.LostRepo
	ProjectRepo data.ProjectRepo
	IssueRepo   data.IssueRepo
}

type ucmInfoController struct {
	ucmv1.UnimplementedServiceServer

	log    log.Fields
	deps   ucmControllerDeps
	conver *ucmConvert
}

// CreateUcmServiceController is a constructor for Fx
func CreateUcmServiceController(deps ucmControllerDeps) UcmServiceController {
	return &ucmInfoController{
		log:    deps.Logger.WithField("controller", "ucm"),
		deps:   deps,
		conver: newUcmConvert(),
	}
}
