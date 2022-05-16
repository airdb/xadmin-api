package mortar

import (
	"context"
	"net/http"

	"github.com/go-masonry/mortar/providers/groups"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
	"google.golang.org/grpc/metadata"
)

func GRPCGatewayMetadataFxOption() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Group:  groups.GRPCGatewayMuxOptions,
			Target: MetadataCarrierOption,
		})
}

// MetadataCarrierOption
// Refer: https://grpc-ecosystem.github.io/grpc-gateway/docs/operations/annotated_context/#get-http-path-pattern
func MetadataCarrierOption() runtime.ServeMuxOption {
	return runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
		md := make(map[string]string)
		if method, ok := runtime.RPCMethod(ctx); ok {
			md["grpc-method"] = method
		}
		if pattern, ok := runtime.HTTPPathPattern(ctx); ok {
			md["http-path-pattern"] = pattern
		}
		return metadata.New(md)
	})
}
