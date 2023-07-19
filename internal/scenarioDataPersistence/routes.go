package scenarioDataPersistence

import (
	"axon/internal/scenarioDataPersistence/handler"
	"github.com/kataras/iris/v12/core/router"
)

func Initialize(app router.Party, tph handler.TracePersistenceHandler) {

	ruleEngineAPI := app.Party("/c/issue")
	{
		ruleEngineAPI.Get("/", tph.GetIssuesListWithDetailsHandler)
		ruleEngineAPI.Get("/{issueHash}", tph.GetIssueDetailsHandler)
		ruleEngineAPI.Get("/{issueHash}/incident", tph.GetIncidentListHandler)
		ruleEngineAPI.Get("/{issueHash}/incident/{incidentId}", tph.GetIncidentDetailsHandler)
		ruleEngineAPI.Get("/{issueHash}/incident/{incidentId}/span/{spanId}", tph.GetSpanRawDataHandler)
	}
}
