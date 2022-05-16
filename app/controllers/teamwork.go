package controllers

import (
	"context"

	"github.com/airdb/xadmin-api/app/data"
	teamworkv1 "github.com/airdb/xadmin-api/genproto/teamwork/v1"
	"github.com/go-masonry/mortar/interfaces/log"
	"go.uber.org/fx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TeamworkServiceController responsible for the business logic of our TeamworkService
type TeamworkServiceController interface {
	teamworkv1.TeamworkServiceServer
}

type teamworkInfoControllerDeps struct {
	fx.In

	DB       data.LostRepo
	LostRepo data.LostRepo
	Logger   log.Logger
}

type teamworkInfoController struct {
	teamworkv1.UnimplementedTeamworkServiceServer

	deps teamworkInfoControllerDeps
	log  log.Fields
}

// CreateTeamworkServiceController is a constructor for Fx
func CreateTeamworkServiceController(deps teamworkInfoControllerDeps) TeamworkServiceController {
	return &teamworkInfoController{
		deps: deps,
		log:  deps.Logger.WithField("controller", "teamwork"),
	}
}

func (c *teamworkInfoController) ListOnduty(ctx context.Context, request *teamworkv1.ListOndutyRequest) (*teamworkv1.ListOndutyResponse, error) {
	c.log.Debug(ctx, "list onduty accepted")
	return nil, status.Errorf(codes.Unimplemented, "method ListOnduty not implemented")
}

func (c *teamworkInfoController) ListTaskByProject(ctx context.Context, request *teamworkv1.ListTaskByProjectRequest) (*teamworkv1.ListTaskByProjectResponse, error) {
	c.log.Debug(ctx, "list task by project accepted")
	return nil, status.Errorf(codes.Unimplemented, "method ListTaskByProject not implemented")
}

func (c *teamworkInfoController) ListTaskByUser(ctx context.Context, request *teamworkv1.ListTaskByUserRequest) (*teamworkv1.ListTaskByUserResponse, error) {
	c.log.Debug(ctx, "list task by user accepted")
	return nil, status.Errorf(codes.Unimplemented, "method ListTaskByUser not implemented")
}
