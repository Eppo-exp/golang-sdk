package eppoclient

import (
	"encoding/json"
	"fmt"
)

const CONFIG_ENDPOINT = "/flag-config/v1/config"
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
		fmt.Println("Failed to fetch config response", err)
		return configResponse{}, err
	}

	var response configResponse
	err = json.Unmarshal(result, &response)
	if err != nil {
		fmt.Println("Failed to unmarshal config response JSON", result)
		fmt.Println(err)
		return configResponse{}, err
	}

	return response, nil
}

func (cr *configurationRequestor) fetchBandits() (banditResponse, error) {
	result, err := cr.httpClient.get(BANDIT_ENDPOINT)
	if err != nil {
		fmt.Println("Failed to fetch bandit response", err)
		return banditResponse{}, err
	}

	var response banditResponse
	err = json.Unmarshal(result, &response)
	if err != nil {
		fmt.Println("Failed to unmarshal bandit response JSON", result)
		fmt.Println(err)
		return banditResponse{}, err
	}

	return response, nil
}
