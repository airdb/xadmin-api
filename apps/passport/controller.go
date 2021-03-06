package passport

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/airdb/xadmin-api/apps/data"
	passportv1 "github.com/airdb/xadmin-api/genproto/passport/v1"
	"github.com/airdb/xadmin-api/pkg/authkit"
	"github.com/casdoor/casdoor-go-sdk/auth"
	"github.com/go-masonry/mortar/constructors"
	"github.com/go-masonry/mortar/interfaces/cfg"
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

	Config cfg.Config
	DB     data.PassportRepo
	Logger log.Logger
}

type passportController struct {
	passportv1.UnimplementedPassportServiceServer

	deps   passportInfoControllerDeps
	log    log.Fields
	conver *passportConvert
	domain string
}

// CreatePassportServiceController is a constructor for Fx
func CreatePassportServiceController(deps passportInfoControllerDeps) PassportServiceController {
	return &passportController{
		deps:   deps,
		log:    deps.Logger.WithField("controller", "passport"),
		conver: &passportConvert{},
		domain: strings.Trim(deps.Config.Get("xadmin.domain").String(), "/"),
	}
}

func (c *passportController) Preset(ctx context.Context, request *passportv1.PresetRequest) (*passportv1.PresetResponse, error) {
	c.log.Debug(ctx, "preset accepted")
	domain := request.GetRedirectUri()
	if len(domain) == 0 {
		domain = fmt.Sprintf("%s/v1/passport/callback", c.domain)
	}

	return &passportv1.PresetResponse{
		Url: auth.GetSigninUrl(domain),
	}, nil
}

func (c *passportController) Login(ctx context.Context, request *passportv1.LoginRequest) (*passportv1.LoginResponse, error) {
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

func (c *passportController) Callback(ctx context.Context, request *passportv1.CallbackRequest) (*passportv1.CallbackResponse, error) {
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
		Token: fmt.Sprintf("%s %s", token.TokenType, token.AccessToken),
	}, nil
}

func (c *passportController) Profile(ctx context.Context, request *passportv1.ProfileRequest) (*passportv1.ProfileResponse, error) {
	c.log.Debug(ctx, "profile accepted")

	user := authkit.FromContextUser(ctx)
	if user == nil {
		return nil, errors.New("can not find profile")
	}

	return &passportv1.ProfileResponse{
		Info: c.conver.FromClaimsUserToProtoInfo(user),
	}, nil
}

func (c *passportController) Logout(ctx context.Context, request *passportv1.LogoutRequest) (*empty.Empty, error) {
	c.log.Debug(ctx, "logout accepted")

	extractor := constructors.DefaultJWTTokenExtractor()
	extractor.FromContext(ctx)
	return &empty.Empty{}, nil
}
