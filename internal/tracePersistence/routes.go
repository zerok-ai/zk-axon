package tracePersistence

import (
	"axon/internal/tracePersistence/handler"
	"github.com/kataras/iris/v12/core/router"
)

func Initialize(app router.Party, tph handler.TracePersistenceHandler) {

	ruleEngineAPI := app.Party("/u/trace")
	{
		ruleEngineAPI.Get("/incident", tph.GetIncidents)
		ruleEngineAPI.Get("/", tph.GetTraces)
		ruleEngineAPI.Get("/metadata", tph.GetTracesMetadata)
		ruleEngineAPI.Get("/raw-data", tph.GetTracesRawData)
		ruleEngineAPI.Get("/metadata/map", tph.GetMetadataMapData)
	}
}
