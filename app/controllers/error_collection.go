package controllers

import (
	"context"

	"github.com/airdb/xadmin-api/app/data"
	errorCollectionv1 "github.com/airdb/xadmin-api/genproto/error_collection/v1"
	"github.com/go-masonry/mortar/interfaces/log"
	"go.uber.org/fx"
)

// ErrorCollectionServiceController responsible for the business logic of our ErrorCollectionService
type ErrorCollectionServiceController interface {
	errorCollectionv1.ErrorCollectionServiceServer
}

type errorCollectionInfoControllerDeps struct {
	fx.In

	DB data.LostRepo
	// LostRepo data.LostRepo
	Logger log.Logger
}

type errorCollectionInfoController struct {
	errorCollectionv1.UnimplementedErrorCollectionServiceServer

	deps errorCollectionInfoControllerDeps
	log  log.Fields
}

// CreateErrorCollectionServiceController is a constructor for Fx
func CreateErrorCollectionServiceController(deps errorCollectionInfoControllerDeps) ErrorCollectionServiceController {
	return &errorCollectionInfoController{
		deps: deps,
		log:  deps.Logger.WithField("controller", "errorCollection"),
	}
}

func (c *errorCollectionInfoController) Collect(ctx context.Context, request *errorCollectionv1.CreateErrorCollectionRequest) (*errorCollectionv1.CreateErrorCollectionResponse, error) {
	c.log.Debug(ctx, "collect request accepted")
	return &errorCollectionv1.CreateErrorCollectionResponse{
		Id: 1,
	}, nil
	// return nil, status.Errorf(codes.Unimplemented, "method Collect not implemented")
}

func (c *errorCollectionInfoController) CreateErrorCollection(ctx context.Context, request *errorCollectionv1.CreateErrorCollectionRequest) (*errorCollectionv1.CreateErrorCollectionResponse, error) {
	c.log.Debug(ctx, "collect request accepted")
	return &errorCollectionv1.CreateErrorCollectionResponse{
		Id: 1,
	}, nil
	// return nil, status.Errorf(codes.Unimplemented, "method Collect not implemented")
}
