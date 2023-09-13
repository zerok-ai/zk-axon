package main

import (
	"axon/internal/config"
	"axon/internal/prometheus"
	promHandler "axon/internal/prometheus/handler"
	"axon/internal/prometheus/model/dto"
	promRepository "axon/internal/prometheus/repository"
	promService "axon/internal/prometheus/service"
	"axon/internal/scenarioDataPersistence"
	scenarioHandler "axon/internal/scenarioDataPersistence/handler"
	scenarioRepository "axon/internal/scenarioDataPersistence/repository"
	scenarioService "axon/internal/scenarioDataPersistence/service"
	"github.com/kataras/iris/v12"
	"github.com/prometheus/client_golang/api"
	zkConfig "github.com/zerok-ai/zk-utils-go/config"
	zkHttpConfig "github.com/zerok-ai/zk-utils-go/http/config"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	store "github.com/zerok-ai/zk-utils-go/storage/redis"
	zkPostgres "github.com/zerok-ai/zk-utils-go/storage/sqlDB/postgres"
	"time"
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
	var metricServerDatasource promRepository.PromQLRepo
	RefreshInterval := 20 * time.Minute
	datasourceStore, err := store.GetVersionedStore[dto.Datasource](cfg.Redis, "integrations", RefreshInterval)
	if err != nil {
		zkLogger.Error(LogTag, "Error creating datasource store: %v\n", err)
		return nil, nil, nil
	}

	dataSourcesMap := datasourceStore.GetAllValues()
	for _, datasource := range dataSourcesMap {
		promClient, err := api.NewClient(api.Config{
			Address: datasource.Url,
		})
		if err != nil {
			zkLogger.Error(LogTag, "Error creating Prometheus client: %v\n", err)
			continue
		}
		dataSources[datasource.Id] = promRepository.NewPromQLRepo(promClient)

		if datasource.MetricServer {
			metricServerDatasource = dataSources[datasource.Id]
		}
	}

	promSvc := promService.NewPrometheusService(metricServerDatasource, dataSources)
	promH := promHandler.NewPrometheusHandler(promSvc, tps, cfg)

	return promH, promSvc, metricServerDatasource
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
