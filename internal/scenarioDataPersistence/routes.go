package scenarioDataPersistence

import (
	"axon/internal/scenarioDataPersistence/handler"
	"axon/utils"
	"github.com/kataras/iris/v12/core/router"
)

func Initialize(app router.Party, tph handler.TracePersistenceHandler) {

	ruleEngineAPI := app.Party("/c/axon")
	{
		ruleEngineAPI.Get("/issue", tph.GetIssuesListWithDetailsHandler)
		ruleEngineAPI.Get("/issue/{issueHash}", tph.GetIssueDetailsHandler)
		ruleEngineAPI.Get("/issue/{issueHash}/incident", tph.GetIncidentListHandler)
		ruleEngineAPI.Get("/issue/{issueHash}/incident/{incidentId}", tph.GetIncidentDetailsHandler)
		ruleEngineAPI.Get("/issue/{issueHash}/incident/{incidentId}/span/{spanId}", tph.GetSpanRawDataHandler)

		ruleEngineAPI.Get("/scenario", tph.GetScenarioDetailsHandler)
		ruleEngineAPI.Get("/scenario/{"+utils.ScenarioId+"}/incident", tph.GetIncidentListForScenarioId)
	}
}
