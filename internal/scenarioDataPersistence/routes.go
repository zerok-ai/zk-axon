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
		//Done
		ruleEngineAPI.Get("/issue/{"+utils.IssueHash+"}", tph.GetIssueDetailsHandler)
		ruleEngineAPI.Get("/issue/{"+utils.IssueHash+"}/incident", tph.GetIncidentListHandler)
		ruleEngineAPI.Get("/issue/{"+utils.IssueHash+"}/incident/{"+utils.IncidentId+"}", tph.GetIncidentDetailsHandler)
		ruleEngineAPI.Get("/issue/{"+utils.IssueHash+"}/incident/{"+utils.IncidentId+"}/span/{"+utils.SpanId+"}", tph.GetSpanRawDataHandler)
		ruleEngineAPI.Post("/issue/error", tph.GetErrorDataHandler)

		ruleEngineAPI.Get("/scenario", tph.GetScenarioDetailsHandler)
		ruleEngineAPI.Get("/scenario/{"+utils.ScenarioId+"}/incident", tph.GetIncidentListForScenarioId)
	}
}
