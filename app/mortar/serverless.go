package mortar

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-masonry/mortar/interfaces/log"
	chiadapter "github.com/serverless-plus/tencent-serverless-go/chi"
	"github.com/serverless-plus/tencent-serverless-go/events"
	"github.com/serverless-plus/tencent-serverless-go/faas"
	"go.uber.org/fx"
)

func BuildServerlessFxOption() fx.Option {
	return fx.Options(
		fx.Provide(NewServerless),
	)
}

type serverlessDeps struct {
	fx.In

	LifeCycle  fx.Lifecycle
	Logger     log.Logger
	Serverless *Serverless
}

// InvokeWebServerlessFxOption creates the entire dependency graph
// and registers all provided fx.LifeCycle hooks
func InvokeWebServerlessFxOption() fx.Option {
	return fx.Invoke(func(deps serverlessDeps) error {
		deps.Serverless.initAdapter()

		deps.LifeCycle.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				go faas.Start(deps.Serverless.scfHandler)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return nil
			},
		})

		return nil
	})
}

type Serverless struct {
	externalMux *http.ServeMux
	chiFaas     *chiadapter.ChiFaas
}

func NewServerless() *Serverless {
	return &Serverless{
		externalMux: http.NewServeMux(),
	}
}

func (s *Serverless) initAdapter() {
	r := chi.NewRouter()
	r.Mount("/", s.externalMux)
	r.Route("/debug", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			buf, _ := json.Marshal(map[string]interface{}{
				"message": "Hello Serverless Chi",
				"query":   r.URL.Query().Get("q"),
			})
			w.Write(buf)
		})
	})

	s.chiFaas = chiadapter.New(r)
}

// Handler serverless faas handler
func (s *Serverless) scfHandler(ctx context.Context, req events.APIGatewayRequest) (events.APIGatewayResponse, error) {
	return s.chiFaas.ProxyWithContext(ctx, req)
}
