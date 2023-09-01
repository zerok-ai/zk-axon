package main

import (
	"axon/internal/config"
	"axon/internal/scenarioDataPersistence"
	"axon/internal/scenarioDataPersistence/handler"
	"axon/internal/scenarioDataPersistence/repository"
	"axon/internal/scenarioDataPersistence/service"
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

	tpr := repository.NewTracePersistenceRepo(zkPostgresRepo)
	tps := service.NewScenarioPersistenceService(tpr)
	tph := handler.NewTracePersistenceHandler(tps, cfg)

	app := newApp(tph)

	configurator := iris.WithConfiguration(iris.Configuration{
		DisablePathCorrection: true,
		//R: This should come from config.
		//R: Resolved
		LogLevel: cfg.LogsConfig.Level,
	})
	//R: Catch and log the error in the below line.
	//R: Resolved
	if err = app.Listen(":"+cfg.Server.Port, configurator); err != nil {
		panic(err)
	}
}

func newApp(tph handler.TracePersistenceHandler) *iris.Application {
	app := iris.Default()

	crs := func(ctx iris.Context) {
		ctx.Header("Access-Control-Allow-Credentials", "true")

		if ctx.Method() == iris.MethodOptions {
			//R: I don't see Get in the below list. What is the list for?
			//R: We only have all Get Apis, why are we allowing all these other methods?
			//ctx.Header("Access-Control-Methods",
			//	"POST, PUT, PATCH, DELETE")
			// Removed this, will test it soon

			ctx.Header("Access-Control-Allow-Headers",
				"Access-Control-Allow-Origin,Content-Type")

			ctx.Header("Access-Control-Max-Age",
				"86400")

			//R: What does this mean? Why are we setting status code here?
			// This is what GPT says:
			//In CORS, a successful preflight request should typically result in a response with a 204 status code to indicate that the request is allowed.
			//The response body is empty (no content)
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

	v1 := app.Party("/v1")
	scenarioDataPersistence.Initialize(v1, tph)

	return app
}
