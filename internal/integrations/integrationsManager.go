package integrations

import (
	"axon/internal/config"
	"axon/internal/integrations/dto"
	tracePersistence "axon/internal/scenarioDataPersistence/service"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	store "github.com/zerok-ai/zk-utils-go/storage/redis"
	"time"
)

const (
	LogTag          = "integrations_manager"
	refreshInterval = 20 * time.Minute
)

type IntegrationsManager struct {
	tracePersistenceService *tracePersistence.TracePersistenceService
	integrationsStore       *store.VersionedStore[dto.Integration]
}

func NewIntegrationsManager(cfg *config.AppConfigs, tracePersistenceService *tracePersistence.TracePersistenceService) *IntegrationsManager {
	zkLogger.Info(LogTag, cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password)
	for name, db := range cfg.Redis.DBs {
		zkLogger.Error(LogTag, name, db)
	}
	integrationsStore, err := store.GetVersionedStore[dto.Integration](&cfg.Redis, "integrations", refreshInterval)
	if err != nil {
		zkLogger.Error(LogTag, "Error creating integrationsMap store: %v\n", err)
		return nil
	}
	return &IntegrationsManager{integrationsStore: integrationsStore, tracePersistenceService: tracePersistenceService}
}

func (im *IntegrationsManager) GetTracePersistenceService() tracePersistence.TracePersistenceService {
	return *im.tracePersistenceService
}

func (im *IntegrationsManager) GetAllIntegrations() map[string]*dto.Integration {
	return im.integrationsStore.GetAllValues()
}

func (im *IntegrationsManager) GetIntegrationById(id string) *dto.Integration {
	return im.integrationsStore.GetAllValues()[id]
}

func (im *IntegrationsManager) GetIntegrationsByType(integrationType dto.IntegrationType) map[string]*dto.Integration {
	var filteredIntegrations = make(map[string]*dto.Integration)
	for _, integration := range im.integrationsStore.GetAllValues() {
		if integration.Type == integrationType && integration.Disabled == false && integration.Deleted == false {
			filteredIntegrations[integration.Id] = integration
		}
	}
	return filteredIntegrations
}
