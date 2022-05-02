package controllers

import (
	"context"
	"errors"

	"github.com/airdb/xadmin-api/app/common"
	"github.com/airdb/xadmin-api/app/data"
	apiv1 "github.com/airdb/xadmin-api/genproto/v1"
	"github.com/casdoor/casdoor-go-sdk/auth"
	"github.com/go-masonry/mortar/constructors"
	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/fx"
)

// PassportServiceController responsible for the business logic of our PassportService
type PassportServiceController interface {
	apiv1.PassportServiceServer
}

type passportInfoControllerDeps struct {
	fx.In

	DB     data.PassportInfoDB
	Logger log.Logger
}

type passportInfoController struct {
	apiv1.UnimplementedPassportServiceServer // if keep this one added even when you change your interface this code will compile

	deps passportInfoControllerDeps
}

// CreatePassportServiceController is a constructor for Fx
func CreatePassportServiceController(deps passportInfoControllerDeps) PassportServiceController {
	return &passportInfoController{
		deps: deps,
	}
}

func (w *passportInfoController) LoginUrl(ctx context.Context, request *apiv1.LoginUrlRequest) (*apiv1.LoginUrlResponse, error) {
	w.deps.Logger.Debug(ctx, "login url accepted")
	return &apiv1.LoginUrlResponse{
		Url: auth.GetSigninUrl("http://localhost:5381/v1/passport/login:callback"),
	}, nil
}

func (w *passportInfoController) Login(ctx context.Context, request *apiv1.LoginRequest) (*apiv1.LoginResponse, error) {
	w.deps.Logger.Debug(ctx, "login accepted")

	info, err := w.deps.DB.GetPassportInfo(ctx, request.GetName())
	if err != nil {
		w.deps.Logger.WithError(err).Debug(ctx, "get passport info")
		return nil, errors.New("name not exist")
	}

	if info.Password != request.Password {
		return nil, errors.New("password incorect")
	}

	return &apiv1.LoginResponse{}, err
}

func (w *passportInfoController) LoginCallback(ctx context.Context, request *apiv1.LoginCallbackRequest) (*apiv1.LoginCallbackResponse, error) {
	w.deps.Logger.Debug(ctx, "login accepted")

	token, err := auth.GetOAuthToken(request.GetCode(), request.GetState())
	if err != nil {
		w.deps.Logger.WithError(err).Debug(ctx, "get oauth token")
		return nil, errors.New("get oauth token error")
	}

	claims, err := auth.ParseJwtToken(token.AccessToken)
	if err != nil {
		w.deps.Logger.WithError(err).Debug(ctx, "parse oauth token")
		return nil, errors.New("parse oauth token error")
	}

	return &apiv1.LoginCallbackResponse{
		Info: &apiv1.PassportInfo{
			Name: claims.User.Name,
		},
		Token: token.AccessToken,
	}, nil
}

func (w *passportInfoController) Profile(ctx context.Context, request *apiv1.ProfileRequest) (*apiv1.ProfileResponse, error) {
	claims := common.FromCurrentCasdoorContext(ctx)
	if claims == nil {
		return nil, errors.New("can not find profile")
	}

	return &apiv1.ProfileResponse{
		Info: &apiv1.PassportInfo{
			Name: claims.User.Name,
		},
	}, nil
}

func (w *passportInfoController) Logout(ctx context.Context, request *apiv1.LogoutRequest) (*empty.Empty, error) {
	extractor := constructors.DefaultJWTTokenExtractor()
	extractor.FromContext(ctx)
	return &empty.Empty{}, nil
}
