package eppoclient

import (
	"encoding/json"
	"fmt"
)

const RAC_ENDPOINT = "/randomized_assignment/v2/config"

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
		// should we panic here or return an error?
		panic("Unauthorized: please check your API key")
	}

	result, err := ecr.configStore.GetConfiguration(experimentKey)

	return result, err
}

func (ecr *experimentConfigurationRequestor) FetchAndStoreConfigurations() {
	var responseBody map[string]json.RawMessage

	configs := dictionary{}
	result := ecr.httpClient.get(RAC_ENDPOINT)

	err := json.Unmarshal([]byte(result), &responseBody)

	if err != nil {
		fmt.Println("Failed to unmarshal RAC response json", result)
		fmt.Println(err)
	}

	err = json.Unmarshal(responseBody["flags"], &configs)

	if err != nil {
		fmt.Println("Failed to unmarshal RAC response json in experiments section", result)
		fmt.Println(err)
	}

	ecr.configStore.SetConfigurations(configs)
}
