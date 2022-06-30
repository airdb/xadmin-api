package mortar

import (
	"context"

	"github.com/airdb/xadmin-api/apps/bchm"
	"github.com/airdb/xadmin-api/apps/data"
	"github.com/airdb/xadmin-api/apps/passport"
	"github.com/airdb/xadmin-api/apps/teamwork"
	"github.com/airdb/xadmin-api/apps/uam"
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
	Bchm     bchmv1.ServiceServer
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
		bchmv1.RegisterServiceServer(srv, deps.Bchm)
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
			return bchmv1.RegisterServiceHandlerFromEndpoint(context.Background(), mux, endpoint, []grpc.DialOption{grpc.WithInsecure()})
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
		passport.CreatePassportServiceService,
		passport.CreatePassportServiceController,
		data.NewPassportRepo,
		passport.CreatePassportServiceValidations,
		// Bchm dependents
		bchm.CreateServiceService,
		bchm.CreateController,
		data.NewFileRepo,
		data.NewCategoryRepo,
		data.NewLostRepo,
		data.NewLostStatRepo,
		bchm.CreateServiceValidations,
		// Teamwork dependents
		teamwork.CreateTeamworkServiceService,
		teamwork.CreateTeamworkServiceController,
		data.NewProjectRepo,
		data.NewIssueRepo,
		teamwork.CreateTeamworkServiceValidations,
		// Uam dependents
		uam.CreateUamServiceService,
		uam.CreateUamServiceController,
		data.NewUserRepo,
		uam.CreateUamServiceValidations,
		// data.NewTeamworkRepo,
	)
}
