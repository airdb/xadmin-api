package mortar

import (
	"context"
	"os"

	"github.com/airdb/xadmin-api/app/common"
	"github.com/casdoor/casdoor-go-sdk/auth"
	"github.com/go-masonry/mortar/interfaces/cfg"
	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/go-masonry/mortar/interfaces/monitor"
	"github.com/go-masonry/mortar/providers/groups"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func BuildCasdoorFxInvoke() fx.Option {
	return fx.Invoke(initCasdoor)
}

func BuildCasdoorFxOption() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Group:  groups.UnaryServerInterceptors,
			Target: CasdoorInterceptor,
		})
}

func initCasdoor(cfg cfg.Config) {
	data := cfg.Get("xadmin.casdoor")
	var ac auth.AuthConfig
	mapstructure.Decode(data.Raw(), &ac)
	jwtKeyContent, err := os.ReadFile(ac.JwtPublicKey)
	if err != nil {
		panic(err)
	}
	auth.InitConfig(
		ac.Endpoint,
		ac.ClientId,
		ac.ClientSecret,
		string(jwtKeyContent),
		ac.OrganizationName,
		ac.ApplicationName)
}

type casdoorInterceptorDeps struct {
	fx.In

	Logger  log.Logger
	Metrics monitor.Metrics `optional:"true"`
}

func CasdoorInterceptor(deps casdoorInterceptorDeps) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			deps.Logger.Debug(ctx, "%v", md)
		}

		token := ""
		if tokens := md.Get("authorization"); len(tokens) > 0 {
			token = tokens[0]
		}

		claims, err := auth.ParseJwtToken(token)
		if err != nil {
			deps.Logger.WithError(err).Debug(ctx, "parse oauth token")
		}

		resp, err = handler(common.NewCasdoorContext(ctx, claims), req)
		return
	}
}
