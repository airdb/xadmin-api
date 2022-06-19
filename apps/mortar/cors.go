package mortar

import (
	"net/http"

	serverInt "github.com/go-masonry/mortar/interfaces/http/server"
	"github.com/go-masonry/mortar/providers/groups"
	"github.com/rs/cors"
	"go.uber.org/fx"
)

func GRPCGatewayCorsFxOption() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Group:  groups.ExternalHTTPInterceptors,
			Target: CorsOption,
		})
}

// CorsOption
func CorsOption() serverInt.GRPCGatewayInterceptor {
	cors.AllowAll()
	c := cors.New(cors.Options{
		AllowOriginRequestFunc: func(r *http.Request, origin string) bool {
			return true
		},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           600,
	})
	return func(handler http.Handler) http.Handler {
		return c.Handler(handler)
	}
}
