package mortar

import (
	"context"
	"encoding/json"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-masonry/mortar/constructors/partial"
	serverInt "github.com/go-masonry/mortar/interfaces/http/server"
	"github.com/go-masonry/mortar/interfaces/log"
	chiadapter "github.com/serverless-plus/tencent-serverless-go/chi"
	"github.com/serverless-plus/tencent-serverless-go/events"
	"github.com/serverless-plus/tencent-serverless-go/faas"
	"go.uber.org/fx"
)

func BuildServerlessFxOption() fx.Option {
	return fx.Options(
		fx.Provide(NewServerless),
		// fx.Provide(fx.Annotated{
		// 	Group: partial.FxGroupBuilderCallbacks,
		// 	Target: func(s *Serverless) partial.BuilderCallback {
		// 		return s.BuilderCallback()
		// 	},
		// }),
		fx.Provide(fx.Annotated{
			Group: partial.FxGroupExternalBuilderCallbacks,
			Target: func(s *Serverless) partial.RESTBuilderCallback {
				return s.ExternalBuilderCallback()
			},
		}),
		// fx.Provide(fx.Annotated{
		// 	Group: partial.FxGroupInternalBuilderCallbacks,
		// 	Target: func(s *Serverless) partial.RESTBuilderCallback {
		// 		return s.InternalBuilderCallback()
		// 	},
		// }),
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

func (s *Serverless) BuilderCallback() partial.BuilderCallback {
	return func(builder serverInt.GRPCWebServiceBuilder) serverInt.GRPCWebServiceBuilder {
		ln, err := s.listen("./output/grpc.socket")
		if err != nil {
			panic(err)
		}
		return builder.SetCustomListener(ln)
	}
}

func (s *Serverless) ExternalBuilderCallback() partial.RESTBuilderCallback {
	return func(builder serverInt.RESTBuilder) serverInt.RESTBuilder {
		// ln, err := s.listen("./output/external.socket")
		// if err != nil {
		// 	panic(err)
		// }
		// return builder.SetCustomListener(ln)
		return builder.SetCustomServer(&http.Server{
			Handler: s.externalMux,
		})
	}
}

func (s *Serverless) InternalBuilderCallback() partial.RESTBuilderCallback {
	return func(builder serverInt.RESTBuilder) serverInt.RESTBuilder {
		ln, err := s.listen("./output/internal.socket")
		if err != nil {
			panic(err)
		}
		return builder.SetCustomListener(ln)
	}
}

func (s *Serverless) listen(name string) (net.Listener, error) {
	ln, err := net.Listen("unix", name)
	if err != nil {
		return nil, err
	}
	// if ln, ok := ln.(*net.UnixListener); ok {
	// 	fp, err := ln.File()
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if err = fp.Sync(); err != nil {
	// 		return nil, err
	// 	}
	// }

	return ln, nil
}
