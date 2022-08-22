package eppoclient

import (
	"encoding/json"
	"fmt"
)

const RAC_ENDPOINT = "/randomized_assignment/config"

type IConfigRequestor interface {
	New(httpClient HttpClient, configStore ConfigurationStore)
	GetConfiguration(key string) (ExperimentConfiguration, error)
	FetchAndStoreConfigurations()
}

type ExperimentConfigurationRequestor struct {
	httpClient  HttpClient
	configStore ConfigurationStore
}

func NewExperimentConfigurationRequestor() ExperimentConfigurationRequestor {
	return ExperimentConfigurationRequestor{}
}

func (ecr *ExperimentConfigurationRequestor) New(httpClient HttpClient, configStore ConfigurationStore) {
	ecr.httpClient = httpClient
	ecr.configStore = configStore
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
		fmt.Println(err)
	}

	err = json.Unmarshal(responseBody["experiments"], &configs)

	if err != nil {
		fmt.Println(err)
	}

	ecr.configStore.SetConfigurations(configs)
}
