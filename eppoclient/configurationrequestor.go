package eppoclient

import (
	"encoding/json"
	"fmt"
)

const RAC_ENDPOINT = "/randomized_assignment/config"

type IConfigRequestor interface {
	GetConfiguration(key string) (ExperimentConfiguration, error)
	FetchAndStoreConfigurations()
}

type ExperimentConfigurationRequestor struct {
	httpClient  HttpClient
	configStore ConfigurationStore
}

func NewExperimentConfigurationRequestor(httpClient HttpClient, configStore ConfigurationStore) *ExperimentConfigurationRequestor {
	return &ExperimentConfigurationRequestor{
		httpClient:  httpClient,
		configStore: configStore,
	}
}

func (ecr *ExperimentConfigurationRequestor) GetConfiguration(experimentKey string) (ExperimentConfiguration, error) {
	if ecr.httpClient.isUnauthorized {
		// should we panic here or return an error?
		panic("Unauthorized: please check your API key")
	}

	result, err := ecr.configStore.GetConfiguration(experimentKey)

	return result, err
}

func (ecr *ExperimentConfigurationRequestor) FetchAndStoreConfigurations() {
	var responseBody map[string]json.RawMessage

	configs := Dictionary{}
	result := ecr.httpClient.Get(RAC_ENDPOINT)

	err := json.Unmarshal([]byte(result), &responseBody)

	if err != nil {
		fmt.Println("Failed to unmarshal RAC response json", result)
		fmt.Println(err)
	}

	err = json.Unmarshal(responseBody["experiments"], &configs)

	if err != nil {
		fmt.Println("Failed to unmarshal RAC response json", result)
		fmt.Println(err)
	}

	ecr.configStore.SetConfigurations(configs)
}
