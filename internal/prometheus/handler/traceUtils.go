package handler

import (
	tracePersistence "axon/internal/scenarioDataPersistence/service"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	zkErrors "github.com/zerok-ai/zk-utils-go/zkerrors"
	"strings"
)

func getPodsAndNSListFromTrace(traceId string, tps tracePersistence.TracePersistenceService) (podsList []string, nsList []string, err *zkErrors.ZkError) {
	spansList, err := tps.GetIncidentDetailsService(traceId, "", 0, 50)
	if err != nil {
		zkLogger.Error(LogTag, "Error while collecting spanList: ", err)
		return nil, nil, err
	}
	podsMap := make(map[string]bool)
	spamItems := spansList.Spans
	for _, spanItems := range spamItems {
		if spanItems.Source != "" {
			podsMap[spanItems.Source] = true
		}
		if spanItems.Destination == "" {
			podsMap[spanItems.Destination] = true
		}
	}
	for podName := range podsMap {
		// split namespace and service name from pod name
		podNameParts := strings.Split(podName, "/")
		if len(podNameParts) != 2 {
			continue
		}
		podsList = append(podsList, podNameParts[1]+".*")
		nsList = append(nsList, podNameParts[0])
	}
	return podsList, nsList, nil
}
