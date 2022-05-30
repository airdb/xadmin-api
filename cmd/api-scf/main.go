package main

import (
	"io"
	"os"

	"github.com/airdb/xadmin-api/app/mortar"
	"github.com/go-masonry/mortar/providers"
	"go.uber.org/fx"
)

const (
	cfgLocal = "./config/config_local.yml"
	jwtKey   = "./config/token_jwt_key.pem"
)

func main() {
	additionalFiles := []string{}
	if _, err := os.Stat(cfgLocal); err == nil {
		additionalFiles = append(additionalFiles, cfgLocal)
	}

	app := createApplication("./config/config.yml", additionalFiles)
	app.Run()
}

func createApplication(configFilePath string, additionalFiles []string) *fx.App {
	return fx.New(
		// fx.NopLogger, // remove fx debug
		mortar.ViperFxOption(configFilePath, additionalFiles...), // Configuration map
		mortar.LoggerFxOption(),                                  // Logger
		mortar.TracerFxOption(),                                  // Jaeger tracing
		mortar.PrometheusFxOption(),                              // Prometheus
		mortar.HttpClientFxOptions(),
		mortar.HttpServerFxOptions(),
		mortar.InternalHttpHandlersFxOptions(),
		// Tutorial service dependencies
		mortar.ServicesAPIsAndOtherDependenciesFxOption(), // register tutorial APIs
		// This one invokes all the above
		mortar.BuildAuthFxInvoke(),
		mortar.BuildServerlessFxOption(),          // serverless plugin
		providers.BuildMortarWebServiceFxOption(), // http server invoker
		mortar.InvokeWebServerlessFxOption(),      // serverless invoker
	)
}

func loadFIleFromEnv(e, f string) error {
	if _, err := os.Stat(f); err != nil && os.IsNotExist(err) {
		if s, ok := os.LookupEnv(e); ok && len(s) > 0 {
			f, err := os.OpenFile(f, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
			if err != nil {
				return err
			}
			n, err := f.WriteString(s)
			if err == nil && n < len(s) {
				return io.ErrShortWrite
			}
			if err := f.Close(); err != nil {
				return err
			}
		}
	}

	return nil
}
