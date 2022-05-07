package mortar

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/airdb/xadmin-api/app/common"
	"github.com/casdoor/casdoor-go-sdk/auth"
	"github.com/go-masonry/mortar/auth/jwt"
	jwtInt "github.com/go-masonry/mortar/interfaces/auth/jwt"
	"github.com/go-masonry/mortar/interfaces/cfg"
	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/go-masonry/mortar/providers/groups"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func BuildAuthFxInvoke() fx.Option {
	return fx.Invoke(func(interceptor *AuthInterceptor) {
		interceptor.log.Debug(context.Background(), "auth")
	})
}

func BuildAuthFxOptions() fx.Option {
	return fx.Options(
		fx.Provide(
			NweAuthInterceptor,
			func(interceptor *AuthInterceptor) jwtInt.TokenExtractor {
				return interceptor.TokenExtractor()
			},
			fx.Annotated{
				Group: groups.UnaryServerInterceptors,
				Target: func(interceptor *AuthInterceptor) grpc.UnaryServerInterceptor {
					return interceptor.Unary()
				},
			},
		),
	)
}

const (
	authorizationHeader                = "authorization"
	grpcGatewayAuthorizationWithPrefix = runtime.MetadataPrefix + "authorization"
)

type authInterceptorDeps struct {
	fx.In

	Config cfg.Config
	Logger log.Logger
}

type AuthInterceptor struct {
	deps authInterceptorDeps
	log  log.Fields
}

func NweAuthInterceptor(deps authInterceptorDeps) *AuthInterceptor {
	interceptor := &AuthInterceptor{
		deps: deps,
		log:  deps.Logger.WithField("mortar", "auth"),
	}
	interceptor.initCasdoor()

	return interceptor
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		token, err := interceptor.contextExtractor(ctx)
		if err != nil {
			interceptor.log.WithError(err).Debug(ctx, "token not exist")
			return handler(ctx, req)
		}

		claims, err := auth.ParseJwtToken(token)
		if err != nil {
			interceptor.log.WithError(err).Debug(ctx, "parse oauth token")
		}

		user, err := auth.GetUser(claims.User.Name)
		if err != nil {
			interceptor.log.WithError(err).Debug(ctx, "get oauth user error")
		}

		return handler(common.NewCasdoorContext(ctx, user), req)
	}
}

func (interceptor *AuthInterceptor) TokenExtractor() jwtInt.TokenExtractor {
	return jwt.Builder().
		SetContextExtractor(interceptor.contextExtractor).
		Build()
}

// Handles use cases where 'authorization' header value is
// 		bearer <token>
//		basic <token>
func (interceptor *AuthInterceptor) contextExtractor(ctx context.Context) (string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		var headerValue string = strings.Join(md.Get(authorizationHeader), " ")
		if !(len(headerValue) > 0) {
			headerValue = strings.Join(md.Get(grpcGatewayAuthorizationWithPrefix), " ")
		}
		if len(headerValue) > 0 {
			rawTokenWithBearer := strings.Split(headerValue, " ")
			if len(rawTokenWithBearer) == 2 {
				return rawTokenWithBearer[1], nil
			}
			return "", fmt.Errorf(
				"%s/%s header value [%s] is of a wrong format",
				authorizationHeader,
				grpcGatewayAuthorizationWithPrefix,
				headerValue,
			)
		}
		return "", fmt.Errorf("context missing %s/%s header",
			authorizationHeader,
			grpcGatewayAuthorizationWithPrefix,
		)
	}
	return "", fmt.Errorf("context missing gRPC incoming key")
}

func (interceptor *AuthInterceptor) initCasdoor() {
	data := interceptor.deps.Config.Get("xadmin.casdoor")
	var authCfg auth.AuthConfig
	mapstructure.Decode(data.Raw(), &authCfg)
	jwtKeyContent, err := os.ReadFile(authCfg.JwtPublicKey)
	if err != nil {
		panic(err)
	}
	auth.InitConfig(
		authCfg.Endpoint,
		authCfg.ClientId,
		authCfg.ClientSecret,
		string(jwtKeyContent),
		authCfg.OrganizationName,
		authCfg.ApplicationName,
	)
}
