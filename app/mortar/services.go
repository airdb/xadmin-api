package mortar

import (
	"context"

	"github.com/airdb/xadmin-api/app/controllers"
	"github.com/airdb/xadmin-api/app/data"
	"github.com/airdb/xadmin-api/app/services"
	"github.com/airdb/xadmin-api/app/validations"
	bchmv1 "github.com/airdb/xadmin-api/genproto/bchm/v1"
	passportv1 "github.com/airdb/xadmin-api/genproto/passport/v1"
	teamworkv1 "github.com/airdb/xadmin-api/genproto/teamwork/v1"
	uamv1 "github.com/airdb/xadmin-api/genproto/uam/v1"
	serverInt "github.com/go-masonry/mortar/interfaces/http/server"
	"github.com/go-masonry/mortar/providers/groups"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type servicesDeps struct {
	fx.In

	// API Implementations
	Passport passportv1.PassportServiceServer
	Bchm     bchmv1.BchmServiceServer
	Teamwork teamworkv1.TeamworkServiceServer
	Uam      uamv1.ServiceServer
}

func ServicesAPIsAndOtherDependenciesFxOption() fx.Option {
	return fx.Options(
		// GRPC Service APIs registration
		fx.Provide(fx.Annotated{
			Group:  groups.GRPCServerAPIs,
			Target: grpcServiceAPIs,
		}),
		// GRPC Gateway Generated Handlers registration
		fx.Provide(fx.Annotated{
			Group:  groups.GRPCGatewayGeneratedHandlers + ",flatten", // "flatten" does this [][]serverInt.GRPCGatewayGeneratedHandlers -> []serverInt.GRPCGatewayGeneratedHandlers
			Target: grpcGatewayHandlers,
		}),
		// Migrate
		data.MigratorFxOption(),
		// All other tutorial dependencies
		servicesDependencies(),
	)
}

func grpcServiceAPIs(deps servicesDeps) serverInt.GRPCServerAPI {
	return func(srv *grpc.Server) {
		passportv1.RegisterPassportServiceServer(srv, deps.Passport)
		bchmv1.RegisterBchmServiceServer(srv, deps.Bchm)
		teamworkv1.RegisterTeamworkServiceServer(srv, deps.Teamwork)
		uamv1.RegisterServiceServer(srv, deps.Uam)
		// Any additional gRPC Implementations should be called here
	}
}

func grpcGatewayHandlers() []serverInt.GRPCGatewayGeneratedHandlers {
	return []serverInt.GRPCGatewayGeneratedHandlers{
		// Register passport REST API
		func(mux *runtime.ServeMux, endpoint string) error {
			return passportv1.RegisterPassportServiceHandlerFromEndpoint(context.Background(), mux, endpoint, []grpc.DialOption{grpc.WithInsecure()})
		},
		// Register Bchm REST API
		func(mux *runtime.ServeMux, endpoint string) error {
			return bchmv1.RegisterBchmServiceHandlerFromEndpoint(context.Background(), mux, endpoint, []grpc.DialOption{grpc.WithInsecure()})
		},
		// Register Teamwork REST API
		func(mux *runtime.ServeMux, endpoint string) error {
			return teamworkv1.RegisterTeamworkServiceHandlerFromEndpoint(context.Background(), mux, endpoint, []grpc.DialOption{grpc.WithInsecure()})
		},
		// Register Uam REST API
		func(mux *runtime.ServeMux, endpoint string) error {
			return uamv1.RegisterServiceHandlerFromEndpoint(context.Background(), mux, endpoint, []grpc.DialOption{grpc.WithInsecure()})
		},
		// Any additional gRPC gateway registrations should be called here
	}
}

func servicesDependencies() fx.Option {
	return fx.Provide(
		// Passport dependents
		services.CreatePassportServiceService,
		controllers.CreatePassportServiceController,
		data.NewPassportRepo,
		validations.CreatePassportServiceValidations,
		// Bchm dependents
		services.CreateBchmServiceService,
		controllers.CreateBchmServiceController,
		data.NewLostRepo,
		validations.CreateBchmServiceValidations,
		// Teamwork dependents
		services.CreateTeamworkServiceService,
		controllers.CreateTeamworkServiceController,
		data.NewProjectRepo,
		data.NewIssueRepo,
		validations.CreateTeamworkServiceValidations,
		// Uam dependents
		services.CreateUamServiceService,
		controllers.CreateUamServiceController,
		data.NewUserRepo,
		validations.CreateUamServiceValidations,
		// data.NewTeamworkRepo,
	)
}
