package mortar

import (
	"context"

	"github.com/airdb/xadmin-api/app/controllers"
	"github.com/airdb/xadmin-api/app/data"
	"github.com/airdb/xadmin-api/app/services"
	"github.com/airdb/xadmin-api/app/validations"
	apiv1 "github.com/airdb/xadmin-api/genproto/v1"
	serverInt "github.com/go-masonry/mortar/interfaces/http/server"
	"github.com/go-masonry/mortar/providers/groups"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type passportServiceDeps struct {
	fx.In

	// API Implementations
	Passport apiv1.PassportServiceServer
}

func PassportServiceAPIsAndOtherDependenciesFxOption() fx.Option {
	return fx.Options(
		// GRPC Service APIs registration
		fx.Provide(fx.Annotated{
			Group:  groups.GRPCServerAPIs,
			Target: passportGRPCServiceAPIs,
		}),
		// GRPC Gateway Generated Handlers registration
		fx.Provide(fx.Annotated{
			Group:  groups.GRPCGatewayGeneratedHandlers + ",flatten", // "flatten" does this [][]serverInt.GRPCGatewayGeneratedHandlers -> []serverInt.GRPCGatewayGeneratedHandlers
			Target: passportGRPCGatewayHandlers,
		}),
		// All other tutorial dependencies
		passportDependencies(),
	)
}

func passportGRPCServiceAPIs(deps passportServiceDeps) serverInt.GRPCServerAPI {
	return func(srv *grpc.Server) {
		apiv1.RegisterPassportServiceServer(srv, deps.Passport)
		// Any additional gRPC Implementations should be called here
	}
}

func passportGRPCGatewayHandlers() []serverInt.GRPCGatewayGeneratedHandlers {
	return []serverInt.GRPCGatewayGeneratedHandlers{
		// Register passport REST API
		func(mux *runtime.ServeMux, endpoint string) error {
			return apiv1.RegisterPassportServiceHandlerFromEndpoint(context.Background(), mux, endpoint, []grpc.DialOption{grpc.WithInsecure()})
		},
		// Any additional gRPC gateway registrations should be called here
	}
}

func passportDependencies() fx.Option {
	return fx.Provide(
		services.CreatePassportServiceService,
		controllers.CreatePassportServiceController,
		data.CreatePassportInfoDB,
		validations.CreatePassportServiceValidations,
	)
}
