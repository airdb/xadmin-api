package teamwork

import (
	"context"

	teamworkv1 "github.com/airdb/xadmin-api/genproto/teamwork/v1"
	"github.com/go-masonry/mortar/interfaces/cfg"
	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/go-masonry/mortar/interfaces/monitor"
	"go.uber.org/fx"
)

type teamworkServiceDeps struct {
	fx.In

	Logger      log.Logger
	Config      cfg.Config
	Controller  TeamworkServiceController
	Validations TeamworkServiceValidations
	Metrics     monitor.Metrics `optional:"true"`
}

type teamworkImpl struct {
	teamworkv1.UnimplementedTeamworkServiceServer // if keep this one added even when you change your interface this code will compile

	deps teamworkServiceDeps
	log  log.Fields
}

func CreateTeamworkServiceService(deps teamworkServiceDeps) teamworkv1.TeamworkServiceServer {
	return &teamworkImpl{
		deps: deps,
		log:  deps.Logger.WithField("service", "teamwork"),
	}
}

func (w *teamworkImpl) ListOnduty(ctx context.Context, request *teamworkv1.ListOndutyRequest) (*teamworkv1.ListOndutyResponse, error) {
	w.log.WithField("request", request).Debug(ctx, "list onduty request")

	return w.deps.Controller.ListOnduty(ctx, request)
}

func (w *teamworkImpl) ListTaskByProject(ctx context.Context, request *teamworkv1.ListTaskByProjectRequest) (*teamworkv1.ListTaskByProjectResponse, error) {
	w.log.WithField("request", request).Debug(ctx, "list task by project request")

	return w.deps.Controller.ListTaskByProject(ctx, request)
}
