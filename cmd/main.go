package main

import (
	"axon/internal/config"
	"axon/internal/prometheus"

	promHandler "axon/internal/prometheus/handler"
	promRepository "axon/internal/prometheus/repository"
	promService "axon/internal/prometheus/service"
	"axon/internal/scenarioDataPersistence"
	"fmt"
	"github.com/prometheus/client_golang/api"
	"os"

	scenarioHandler "axon/internal/scenarioDataPersistence/handler"
	scenarioRepository "axon/internal/scenarioDataPersistence/repository"
	scenarioService "axon/internal/scenarioDataPersistence/service"
	"github.com/kataras/iris/v12"
	zkConfig "github.com/zerok-ai/zk-utils-go/config"
	zkHttpConfig "github.com/zerok-ai/zk-utils-go/http/config"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	zkPostgres "github.com/zerok-ai/zk-utils-go/storage/sqlDB/postgres"
)

var LogTag = "main"

func main() {
	var cfg config.AppConfigs
	if err := zkConfig.ProcessArgs[config.AppConfigs](&cfg); err != nil {
		panic(err)
	}

	zkLogger.Info(LogTag, "")
	zkLogger.Info(LogTag, "********* Initializing Application *********")
	zkHttpConfig.Init(cfg.Http.Debug)
	zkLogger.Init(cfg.LogsConfig)
	zkPostgresRepo, err := zkPostgres.NewZkPostgresRepo(cfg.Postgres)
	if err != nil {
		return
	}

	zkLogger.Debug(LogTag, "Parsed Configuration", cfg)

	tpr := scenarioRepository.NewTracePersistenceRepo(zkPostgresRepo)
	tps := scenarioService.NewScenarioPersistenceService(tpr)
	tph := scenarioHandler.NewTracePersistenceHandler(tps, cfg)

	client, err := api.NewClient(api.Config{
		Address: "http://localhost:9090",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating Prometheus client: %v\n", err)
		os.Exit(1)
	}

	promRepo := promRepository.NewPromQLRepo(client)
	promSvc := promService.NewPrometheusService(promRepo)
	promH := promHandler.NewPrometheusHandler(promSvc, cfg)

	app := newApp()
	v1 := app.Party("/v1")
	scenarioDataPersistence.Initialize(v1, tph)
	prometheus.Initialize(v1, promH)

	configurator := iris.WithConfiguration(iris.Configuration{
		DisablePathCorrection: true,
		LogLevel:              cfg.LogsConfig.Level,
	})
	if err = app.Listen(":"+cfg.Server.Port, configurator); err != nil {
		panic(err)
	}
}

func newApp() *iris.Application {
	app := iris.Default()

	crs := func(ctx iris.Context) {
		ctx.Header("Access-Control-Allow-Credentials", "true")

		if ctx.Method() == iris.MethodOptions {
			//ctx.Header("Access-Control-Methods",
			//	"POST, PUT, PATCH, DELETE")
			// Removed this, will test it soon

			ctx.Header("Access-Control-Allow-Headers",
				"Access-Control-Allow-Origin,Content-Type")

			ctx.Header("Access-Control-Max-Age",
				"86400")

			ctx.StatusCode(iris.StatusNoContent)
			return
		}

		ctx.Next()
	}

	app.UseRouter(crs)
	app.AllowMethods(iris.MethodOptions)

	app.Get("/healthz", func(ctx iris.Context) {
		ctx.WriteString("pong")
	}).Describe("healthcheck")

	return app
}
