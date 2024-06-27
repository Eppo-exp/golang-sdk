package eppoclient

import (
	"encoding/json"
	"fmt"
)

const UFC_ENDPOINT = "/flag-config/v1/config"
const BANDIT_ENDPOINT = "/flag-config/v1/bandits"

type configurationRequestor struct {
	httpClient  httpClient
	configStore *configurationStore
}

func newConfigurationRequestor(httpClient httpClient, configStore *configurationStore) *configurationRequestor {
	return &configurationRequestor{
		httpClient:  httpClient,
		configStore: configStore,
	}
}

func (cr *configurationRequestor) IsAuthorized() bool {
	return !cr.httpClient.isUnauthorized
}

func (cr *configurationRequestor) FetchAndStoreConfigurations() {
	configuration, err := cr.fetchConfiguration()
	if err != nil {
		return
	}

	cr.configStore.setConfiguration(configuration)
}

func (cr *configurationRequestor) fetchConfiguration() (configuration, error) {
	var config configuration
	var err error

	config.ufc, err = cr.fetchUfc()
	if err != nil {
		return configuration{}, err
	}

	if config.ufc.Bandits != nil {
		config.bandits, err = cr.fetchBandits()
		if err != nil {
			return configuration{}, err
		}
	}

	return config, nil
}

func (cr *configurationRequestor) fetchUfc() (ufcResponse, error) {
	result, err := cr.httpClient.get(UFC_ENDPOINT)
	if err != nil {
		fmt.Println("Failed to fetch UFC response", err)
		return ufcResponse{}, err
	}

	var ufc ufcResponse
	err = json.Unmarshal(result, &ufc)
	if err != nil {
		fmt.Println("Failed to unmarshal UFC response JSON", result)
		fmt.Println(err)
		return ufcResponse{}, err
	}

	return ufc, nil
}

func (cr *configurationRequestor) fetchBandits() (banditResponse, error) {
	result, err := cr.httpClient.get(BANDIT_ENDPOINT)
	if err != nil {
		fmt.Println("Failed to fetch bandit response", err)
		return banditResponse{}, err
	}

	var bandits banditResponse
	err = json.Unmarshal(result, &bandits)
	if err != nil {
		fmt.Println("Failed to unmarshal bandit response JSON", result)
		fmt.Println(err)
		return banditResponse{}, err
	}

	return bandits, nil
}
