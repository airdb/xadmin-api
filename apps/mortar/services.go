package mortar

import (
	"context"

	"github.com/airdb/xadmin-api/apps/bchm"
	"github.com/airdb/xadmin-api/apps/data"
	"github.com/airdb/xadmin-api/apps/docs"
	"github.com/airdb/xadmin-api/apps/passport"
	"github.com/airdb/xadmin-api/apps/teamwork"
	"github.com/airdb/xadmin-api/apps/uam"
	bchmv1 "github.com/airdb/xadmin-api/genproto/bchm/v1"
	docsv1 "github.com/airdb/xadmin-api/genproto/docs/v1"
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
	Docs     docsv1.ServiceServer
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
		docsv1.RegisterServiceServer(srv, deps.Docs)
		passportv1.RegisterPassportServiceServer(srv, deps.Passport)
		bchmv1.RegisterServiceServer(srv, deps.Bchm)
		teamworkv1.RegisterTeamworkServiceServer(srv, deps.Teamwork)
		uamv1.RegisterServiceServer(srv, deps.Uam)
		// Any additional gRPC Implementations should be called here
	}
}

func grpcGatewayHandlers() []serverInt.GRPCGatewayGeneratedHandlers {
	dialOpts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	return []serverInt.GRPCGatewayGeneratedHandlers{
		// Register docs REST API
		func(mux *runtime.ServeMux, endpoint string) error {
			return docsv1.RegisterServiceHandlerFromEndpoint(context.Background(), mux, endpoint, dialOpts)
		},
		// Register passport REST API
		func(mux *runtime.ServeMux, endpoint string) error {
			return passportv1.RegisterPassportServiceHandlerFromEndpoint(context.Background(), mux, endpoint, dialOpts)
		},
		// Register Bchm REST API
		func(mux *runtime.ServeMux, endpoint string) error {
			return bchmv1.RegisterServiceHandlerFromEndpoint(context.Background(), mux, endpoint, dialOpts)
		},
		// Register Teamwork REST API
		func(mux *runtime.ServeMux, endpoint string) error {
			return teamworkv1.RegisterTeamworkServiceHandlerFromEndpoint(context.Background(), mux, endpoint, dialOpts)
		},
		// Register Uam REST API
		func(mux *runtime.ServeMux, endpoint string) error {
			return uamv1.RegisterServiceHandlerFromEndpoint(context.Background(), mux, endpoint, dialOpts)
		},
		// Any additional gRPC gateway registrations should be called here
	}
}

func servicesDependencies() fx.Option {
	return fx.Provide(
		// Docs
		docs.CreatePassportServiceService,
		docs.CreateController,
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
