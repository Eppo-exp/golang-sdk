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

func (ecr *configurationRequestor) IsAuthorized() bool {
	return !ecr.httpClient.isUnauthorized
}

func (ecr *configurationRequestor) FetchAndStoreConfigurations() {
	result, err := ecr.httpClient.get(UFC_ENDPOINT)
	if err != nil {
		fmt.Println("Failed to fetch UFC response", err)
		return
	}

	var wrapper ufcResponse
	err = json.Unmarshal(result, &wrapper)
	if err != nil {
		fmt.Println("Failed to unmarshal UFC response JSON", result)
		fmt.Println(err)
		return
	}

	ecr.configStore.setFlagsConfiguration(wrapper.Flags)

	if wrapper.Bandits != nil {
		ecr.fetchAndStoreBandits()
	}
}

func (ecr *configurationRequestor) fetchAndStoreBandits() {
	result, err := ecr.httpClient.get(BANDIT_ENDPOINT)
	if err != nil {
		fmt.Println("Failed to fetch bandit response", err)
		return
	}

	var bandits banditResponse
	err = json.Unmarshal(result, &bandits)
	if err != nil {
		fmt.Println("Failed to unmarshal bandit response JSON", result)
		fmt.Println(err)
		return
	}

	ecr.configStore.setBanditsConfiguration(bandits.Bandits)
}
