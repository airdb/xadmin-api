package main

import (
	"github.com/airdb/xadmin-api/app/mortar"
	"github.com/go-masonry/mortar/providers"
	"go.uber.org/fx"
)

func main() {
	app := createApplication("./config/config.yml", []string{"./config/config_local.yml"})
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
