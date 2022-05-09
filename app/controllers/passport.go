package controllers

import (
	"context"
	"errors"

	"github.com/airdb/xadmin-api/app/data"
	passportv1 "github.com/airdb/xadmin-api/genproto/passport/v1"
	"github.com/airdb/xadmin-api/pkg/authkit"
	"github.com/casdoor/casdoor-go-sdk/auth"
	"github.com/go-masonry/mortar/constructors"
	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/fx"
)

// PassportServiceController responsible for the business logic of our PassportService
type PassportServiceController interface {
	passportv1.PassportServiceServer
}

type passportInfoControllerDeps struct {
	fx.In

	DB     data.PassportRepo
	Logger log.Logger
}

type passportInfoController struct {
	passportv1.UnimplementedPassportServiceServer

	deps   passportInfoControllerDeps
	log    log.Fields
	conver *passportConvert
}

// CreatePassportServiceController is a constructor for Fx
func CreatePassportServiceController(deps passportInfoControllerDeps) PassportServiceController {
	return &passportInfoController{
		deps:   deps,
		log:    deps.Logger.WithField("controller", "passport"),
		conver: &passportConvert{},
	}
}

func (c *passportInfoController) Preset(ctx context.Context, request *passportv1.PresetRequest) (*passportv1.PresetResponse, error) {
	c.log.Debug(ctx, "preset accepted")

	return &passportv1.PresetResponse{
		Url: auth.GetSigninUrl("http://localhost:5381/v1/passport/callback"),
	}, nil
}

func (c *passportInfoController) Login(ctx context.Context, request *passportv1.LoginRequest) (*passportv1.LoginResponse, error) {
	c.log.Debug(ctx, "login accepted")

	info, err := c.deps.DB.GetInfo(ctx, request.GetName())
	if err != nil {
		c.log.WithError(err).Debug(ctx, "get passport info")
		return nil, errors.New("name not exist")
	}

	if info.Password != request.Password {
		return nil, errors.New("password incorect")
	}

	return &passportv1.LoginResponse{}, err
}

func (c *passportInfoController) Callback(ctx context.Context, request *passportv1.CallbackRequest) (*passportv1.CallbackResponse, error) {
	c.log.Debug(ctx, "callback accepted")

	token, err := auth.GetOAuthToken(request.GetCode(), request.GetState())
	if err != nil {
		c.log.WithError(err).Debug(ctx, "get oauth token")
		return nil, errors.New("get oauth token error")
	}

	claims, err := auth.ParseJwtToken(token.AccessToken)
	if err != nil {
		c.log.WithError(err).Debug(ctx, "parse oauth token")
		return nil, errors.New("parse oauth token error")
	}

	user, err := auth.GetUser(claims.User.Name)
	if err != nil {
		c.log.WithError(err).Debug(ctx, "parse oauth token")
	}

	return &passportv1.CallbackResponse{
		Info:  c.conver.FromClaimsUserToProtoInfo(user),
		Token: token.AccessToken,
	}, nil
}

func (c *passportInfoController) Profile(ctx context.Context, request *passportv1.ProfileRequest) (*passportv1.ProfileResponse, error) {
	c.log.Debug(ctx, "profile accepted")

	user := authkit.FromContextUser(ctx)
	if user == nil {
		return nil, errors.New("can not find profile")
	}

	return &passportv1.ProfileResponse{
		Info: c.conver.FromClaimsUserToProtoInfo(user),
	}, nil
}

func (c *passportInfoController) Logout(ctx context.Context, request *passportv1.LogoutRequest) (*empty.Empty, error) {
	c.log.Debug(ctx, "logout accepted")

	extractor := constructors.DefaultJWTTokenExtractor()
	extractor.FromContext(ctx)
	return &empty.Empty{}, nil
}
