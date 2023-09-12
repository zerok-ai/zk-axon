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

	app := newApp()
	v1 := app.Party("/v1")

	tracePersistenceHandler, tracePersistenceService, _ := getTracePersistenceHandler(cfg)
	scenarioDataPersistence.Initialize(v1, tracePersistenceHandler)

	promQueryHandler, _, _ := getPrometheusHandler(cfg, tracePersistenceService)
	prometheus.Initialize(v1, promQueryHandler)

	configurator := iris.WithConfiguration(iris.Configuration{
		DisablePathCorrection: true,
		LogLevel:              cfg.LogsConfig.Level,
	})
	if err := app.Listen(":"+cfg.Server.Port, configurator); err != nil {
		panic(err)
	}
}

func getPrometheusHandler(cfg config.AppConfigs, tps scenarioService.TracePersistenceService) (promHandler.PrometheusHandler, promService.PrometheusService, promRepository.PromQLRepo) {
	dataSources := make(map[string]promRepository.PromQLRepo)
	promAddr := cfg.Prometheus.Protocol + "://" + cfg.Prometheus.Host + ":" + cfg.Prometheus.Port
	client, err := api.NewClient(api.Config{
		Address: promAddr,
	})
	if err != nil {
		errorStr, err := fmt.Fprintf(os.Stderr, "Error creating Prometheus client: %v\n", err)
		if err != nil {
			zkLogger.Error(LogTag, err)
			return nil, nil, nil
		}
		zkLogger.Error(LogTag, errorStr)
		os.Exit(1)
	}

	promRepo := promRepository.NewPromQLRepo(client)
	promSvc := promService.NewPrometheusService(promRepo, dataSources)
	promH := promHandler.NewPrometheusHandler(promSvc, tps, cfg)

	return promH, promSvc, promRepo
}

func getTracePersistenceHandler(cfg config.AppConfigs) (scenarioHandler.TracePersistenceHandler, scenarioService.TracePersistenceService, scenarioRepository.TracePersistenceRepo) {
	zkPostgresRepo, err := zkPostgres.NewZkPostgresRepo(cfg.Postgres)
	if err != nil {
		return nil, nil, nil
	}

	zkLogger.Debug(LogTag, "Parsed Configuration", cfg)

	tpr := scenarioRepository.NewTracePersistenceRepo(zkPostgresRepo)
	tps := scenarioService.NewScenarioPersistenceService(tpr)
	tph := scenarioHandler.NewTracePersistenceHandler(tps, cfg)
	return tph, tps, tpr
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
