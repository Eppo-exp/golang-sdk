package eppoclient

import (
	"encoding/json"

	"github.com/Eppo-exp/golang-sdk/v4/eppoclient/applicationlogger"
)

const CONFIG_ENDPOINT = "/flag-config/v1/config"
const BANDIT_ENDPOINT = "/flag-config/v1/bandits"

type configurationRequestor struct {
	httpClient        httpClient
	configStore       *configurationStore
	applicationLogger applicationlogger.Logger
}

func newConfigurationRequestor(httpClient httpClient, configStore *configurationStore, applicationLogger applicationlogger.Logger) *configurationRequestor {
	return &configurationRequestor{
		httpClient:        httpClient,
		configStore:       configStore,
		applicationLogger: applicationLogger,
	}
}

func (cr *configurationRequestor) IsAuthorized() bool {
	return !cr.httpClient.isUnauthorized
}

func (cr *configurationRequestor) FetchAndStoreConfigurations() {
	configuration, err := cr.fetchConfiguration()
	if err != nil {
		cr.applicationLogger.Error("Failed to fetch UFC response", err)
		return
	}

	cr.configStore.setConfiguration(configuration)
}

func (cr *configurationRequestor) fetchConfiguration() (configuration, error) {
	var config configuration
	var err error

	config.flags, err = cr.fetchConfig()
	if err != nil {
		return configuration{}, err
	}

	if config.flags.Bandits != nil {
		config.bandits, err = cr.fetchBandits()
		if err != nil {
			return configuration{}, err
		}
	}

	return config, nil
}

func (cr *configurationRequestor) fetchConfig() (configResponse, error) {
	result, err := cr.httpClient.get(CONFIG_ENDPOINT)
	if err != nil {
		cr.applicationLogger.Error("Failed to fetch config response", err)
		return configResponse{}, err
	}

	var response configResponse
	err = json.Unmarshal(result, &response)
	if err != nil {
		cr.applicationLogger.Error("Failed to unmarshal config response JSON", result)
		cr.applicationLogger.Error(err)
		return configResponse{}, err
	}

	return response, nil
}

func (cr *configurationRequestor) fetchBandits() (banditResponse, error) {
	result, err := cr.httpClient.get(BANDIT_ENDPOINT)
	if err != nil {
		cr.applicationLogger.Error("Failed to fetch bandit response", err)
		return banditResponse{}, err
	}

	var response banditResponse
	err = json.Unmarshal(result, &response)
	if err != nil {
		cr.applicationLogger.Error("Failed to unmarshal bandit response JSON", result)
		cr.applicationLogger.Error(err)
		return banditResponse{}, err
	}

	return response, nil
}
