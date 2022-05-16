package controllers

import (
	"context"
	"errors"

	"github.com/airdb/xadmin-api/app/data"
	teamworkv1 "github.com/airdb/xadmin-api/genproto/teamwork/v1"
	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/fx"
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
	// conver *teamworkConvert
}

// CreateTeamworkServiceController is a constructor for Fx
func CreateTeamworkServiceController(deps teamworkInfoControllerDeps) TeamworkServiceController {
	return &teamworkInfoController{
		deps: deps,
		log:  deps.Logger.WithField("controller", "teamwork"),
		// conver: newTeamworkConvert(),
	}
}

func (c *teamworkInfoController) ListOnduty(ctx context.Context, request *teamworkv1.ListOndutyRequest) (*teamworkv1.ListOndutyResponse, error) {
	c.log.Debug(ctx, "list onduty accepted")
	return nil, nil
}

func (c *teamworkInfoController) ListTaskByProject(ctx context.Context, request *teamworkv1.ListTaskByProjectRequest) (*teamworkv1.ListTaskByProjectResponse, error) {
	c.log.Debug(ctx, "list task by project accepted")
	return nil, nil
}

func (c *teamworkInfoController) ListTaskByUser(ctx context.Context, request *teamworkv1.ListTaskByUserRequest) (*teamworkv1.ListTaskByUserResponse, error) {
	c.log.Debug(ctx, "list task by user accepted")
	return nil, nil
}

func (c *teamworkInfoController) ListLosts(ctx context.Context, request *teamworkv1.ListLostsRequest) (*teamworkv1.ListLostsResponse, error) {
	c.log.Debug(ctx, "list losts accepted")
	return nil, nil

	/*
		total, filtered, err := c.deps.LostRepo.Count(ctx, request)
		if err != nil {
			c.log.WithError(err).Debug(ctx, "list losts count error")
			return nil, errors.New("list losts count error")
		}
		if total == 0 {
			return nil, errors.New("losts is empty")
		}

		items, err := c.deps.LostRepo.List(ctx, request)
		if err != nil {
			c.log.WithError(err).Debug(ctx, "list losts error")
			return nil, errors.New("list losts error")
		}

		return &teamworkv1.ListLostsResponse{
			TotalSize:    total,
			FilteredSize: filtered,
			Items: func() []*teamworkv1.Lost {
				res := make([]*teamworkv1.Lost, len(items))
				for i := 0; i < len(items); i++ {
					res[i] = c.conver.FromModelLostToProtoLost(items[i])
				}
				return res
			}(),
		}, nil
	*/
}

func (c *teamworkInfoController) GetLost(ctx context.Context, request *teamworkv1.GetLostRequest) (*teamworkv1.GetLostResponse, error) {
	c.log.Debug(ctx, "get lost accepted")
	return nil, errors.New("")

	/*
		item, err := c.deps.LostRepo.Get(ctx, uint(request.GetId()))
		if err != nil {
			c.log.WithError(err).Debug(ctx, "get lost error")
			return nil, errors.New("name not exist")
		}

		return &teamworkv1.GetLostResponse{
			Item: c.conver.FromModelLostToProtoLost(item),
		}, err
	*/
}

func (c *teamworkInfoController) CreateLost(ctx context.Context, request *teamworkv1.CreateLostRequest) (*teamworkv1.CreateLostResponse, error) {
	c.log.Debug(ctx, "create lost accepted")
	return nil, errors.New("")
	/*
		item := c.conver.FromProtoCreateLostToModelLost(request)
		err := c.deps.LostRepo.Create(ctx, item)
		if err != nil {
			c.log.WithError(err).Debug(ctx, "create lost item failed")
			return nil, errors.New("create lost item failed")
		}

		return &teamworkv1.CreateLostResponse{
			Item: c.conver.FromModelLostToProtoLost(item),
		}, err
	*/
}

func (c *teamworkInfoController) UpdateLost(ctx context.Context, request *teamworkv1.UpdateLostRequest) (*teamworkv1.UpdateLostResponse, error) {
	c.log.Debug(ctx, "update lost accepted")
	return nil, errors.New("")

	/*
		data := c.conver.FromProtoLostToModelLost(request.GetItem())

		fm := querykit.NewField(request.GetUpdateMask(), request.GetItem()).WithAction("update")

		err := c.deps.LostRepo.Update(ctx, data.ID, data, fm)
		if err != nil {
			c.log.WithError(err).Debug(ctx, "update lost item failed")
			return nil, errors.New("update lost item failed")
		}

		item, err := c.deps.LostRepo.Get(ctx, uint(data.ID))
		if err != nil || item == nil {
			c.log.WithError(err).Debug(ctx, "update lost item not exist")
			return nil, errors.New("update lost item not exist")
		}

		return &teamworkv1.UpdateLostResponse{
			Item: c.conver.FromModelLostToProtoLost(item),
		}, err
	*/
}

func (c *teamworkInfoController) DeleteLost(ctx context.Context, request *teamworkv1.DeleteLostRequest) (*empty.Empty, error) {
	c.log.Debug(ctx, "delete lost accepted")

	err := c.deps.LostRepo.Delete(ctx, uint(request.GetId()))
	if err != nil {
		c.log.WithError(err).Debug(ctx, "delete lost item failed")
		return nil, errors.New("delete lost item failed")
	}

	return &empty.Empty{}, nil
}
