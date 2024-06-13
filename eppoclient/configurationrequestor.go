package eppoclient

import (
	"encoding/json"
	"fmt"
)

const UFC_ENDPOINT = "/flag-config/v1/config"

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

	// Now wrapper.Flags contains all configurations mapped by their keys
	// Pass this map directly to SetConfigurations
	err = ecr.configStore.setFlagsConfiguration(wrapper.Flags)
	if err != nil {
		fmt.Println("Failed to set configurations in configuration store", err)
	}
}
