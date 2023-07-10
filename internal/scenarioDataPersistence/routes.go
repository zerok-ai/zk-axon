package scenarioDataPersistence

import (
	"axon/internal/scenarioDataPersistence/handler"
	"github.com/kataras/iris/v12/core/router"
)

func Initialize(app router.Party, tph handler.TracePersistenceHandler) {

	ruleEngineAPI := app.Party("/c/issue")
	{
		ruleEngineAPI.Get("/", tph.GetIssuesListWithDetailsHandler)
		ruleEngineAPI.Get("/{issueId}", tph.GetIssueDetailsHandler)
		ruleEngineAPI.Get("/{issueId}/incident", tph.GetIncidentListHandler)
		ruleEngineAPI.Get("/{issueId}/incident/{incidentId}", tph.GetIncidentDetailsHandler)
		ruleEngineAPI.Get("/{issueId}/incident/{incidentId}/span/{spanId}", tph.GetSpanRawDataHandler)
	}
}
