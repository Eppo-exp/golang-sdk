package eppoclient

import (
	"encoding/json"
	"errors"
	"fmt"
)

const RAC_ENDPOINT = "/randomized_assignment/v3/config"

type racResponse struct {
	Flags map[string]experimentConfiguration `json:"flags"`
}

type iConfigRequestor interface {
	GetConfiguration(key string) (experimentConfiguration, error)
	FetchAndStoreConfigurations()
}

type experimentConfigurationRequestor struct {
	httpClient  httpClient
	configStore *configurationStore
}

func newExperimentConfigurationRequestor(httpClient httpClient, configStore *configurationStore) *experimentConfigurationRequestor {
	return &experimentConfigurationRequestor{
		httpClient:  httpClient,
		configStore: configStore,
	}
}

func (ecr *experimentConfigurationRequestor) GetConfiguration(experimentKey string) (experimentConfiguration, error) {
	if ecr.httpClient.isUnauthorized {
		return experimentConfiguration{}, errors.New("Unauthorized: please check your API key")
	}

	result, err := ecr.configStore.GetConfiguration(experimentKey)

	return result, err
}

func (ecr *experimentConfigurationRequestor) FetchAndStoreConfigurations() {
	result, err := ecr.httpClient.get(RAC_ENDPOINT)
	if err != nil {
		fmt.Println("Failed to fetch RAC response", err)
		return
	}
	var wrapper racResponse

	// Unmarshal JSON data directly into the wrapper struct
	err = json.Unmarshal([]byte(result), &wrapper)
	if err != nil {
		fmt.Println("Failed to unmarshal RAC response JSON", result)
		fmt.Println(err)
		return
	}

	// Now wrapper.Flags contains all configurations mapped by their keys
	// Pass this map directly to SetConfigurations
	err = ecr.configStore.SetConfigurations(wrapper.Flags)
	if err != nil {
		fmt.Println("Failed to set configurations in cache", err)
	}
}
