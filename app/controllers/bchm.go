package controllers

import (
	"context"
	"errors"

	"github.com/airdb/xadmin-api/app/data"
	bchmv1 "github.com/airdb/xadmin-api/genproto/bchm/v1"
	"github.com/airdb/xadmin-api/pkg/querykit"
	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/fx"
)

// BchmServiceController responsible for the business logic of our BchmService
type BchmServiceController interface {
	bchmv1.BchmServiceServer
}

type bchmInfoControllerDeps struct {
	fx.In

	Logger   log.Logger
	LostRepo data.LostRepo
}

type bchmInfoController struct {
	bchmv1.UnimplementedBchmServiceServer

	deps   bchmInfoControllerDeps
	log    log.Fields
	conver *bchmConvert
}

// CreateBchmServiceController is a constructor for Fx
func CreateBchmServiceController(deps bchmInfoControllerDeps) BchmServiceController {
	return &bchmInfoController{
		deps:   deps,
		log:    deps.Logger.WithField("controller", "bchm"),
		conver: newBchmConvert(),
	}
}

func (c *bchmInfoController) ListLosts(ctx context.Context, request *bchmv1.ListLostsRequest) (*bchmv1.ListLostsResponse, error) {
	c.log.Debug(ctx, "list losts accepted")

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

	return &bchmv1.ListLostsResponse{
		TotalSize:    total,
		FilteredSize: filtered,
		Items: func() []*bchmv1.Lost {
			res := make([]*bchmv1.Lost, len(items))
			for i := 0; i < len(items); i++ {
				res[i] = c.conver.FromModelLostToProtoLost(items[i])
			}
			return res
		}(),
	}, nil
}

func (c *bchmInfoController) GetLost(ctx context.Context, request *bchmv1.GetLostRequest) (*bchmv1.GetLostResponse, error) {
	c.log.Debug(ctx, "get lost accepted")

	item, err := c.deps.LostRepo.Get(ctx, uint(request.GetId()))
	if err != nil {
		c.log.WithError(err).Debug(ctx, "get lost error")
		return nil, errors.New("name not exist")
	}

	return &bchmv1.GetLostResponse{
		Item: c.conver.FromModelLostToProtoLost(item),
	}, err
}

func (c *bchmInfoController) CreateLost(ctx context.Context, request *bchmv1.CreateLostRequest) (*bchmv1.CreateLostResponse, error) {
	c.log.Debug(ctx, "create lost accepted")

	item := c.conver.FromProtoCreateLostToModelLost(request)
	err := c.deps.LostRepo.Create(ctx, item)
	if err != nil {
		c.log.WithError(err).Debug(ctx, "create lost item failed")
		return nil, errors.New("create lost item failed")
	}

	return &bchmv1.CreateLostResponse{
		Item: c.conver.FromModelLostToProtoLost(item),
	}, err
}

func (c *bchmInfoController) UpdateLost(ctx context.Context, request *bchmv1.UpdateLostRequest) (*bchmv1.UpdateLostResponse, error) {
	c.log.Debug(ctx, "update lost accepted")
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

	return &bchmv1.UpdateLostResponse{
		Item: c.conver.FromModelLostToProtoLost(item),
	}, err
}

func (c *bchmInfoController) DeleteLost(ctx context.Context, request *bchmv1.DeleteLostRequest) (*empty.Empty, error) {
	c.log.Debug(ctx, "delete lost accepted")

	err := c.deps.LostRepo.Delete(ctx, uint(request.GetId()))
	if err != nil {
		c.log.WithError(err).Debug(ctx, "delete lost item failed")
		return nil, errors.New("delete lost item failed")
	}

	return &empty.Empty{}, nil
}
