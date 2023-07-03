package scenarioDataPersistence

import (
	"axon/internal/scenarioDataPersistence/handler"
	"github.com/kataras/iris/v12/core/router"
)

func Initialize(app router.Party, tph handler.TracePersistenceHandler) {

	ruleEngineAPI := app.Party("/c/trace")
	{
		ruleEngineAPI.Get("/incident", tph.GetIncidents)
		ruleEngineAPI.Get("/", tph.GetTraces)
		ruleEngineAPI.Get("/metadata", tph.GetSpan)
		ruleEngineAPI.Get("/raw-data", tph.GetSpanRawData)
		ruleEngineAPI.Get("/metadata/map", tph.GetMetadataMapData)
	}
}
