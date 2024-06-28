package eppoclient

import (
	"encoding/json"

	"github.com/Eppo-exp/golang-sdk/v4/eppoclient/applicationlogger"
)

const UFC_ENDPOINT = "/flag-config/v1/config"

type iConfigRequestor interface {
	GetConfiguration(key string) (flagConfiguration, error)
	FetchAndStoreConfigurations()
}

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

func (ecr *configurationRequestor) GetConfiguration(experimentKey string) (flagConfiguration, error) {
	if ecr.httpClient.isUnauthorized {
		// should we panic here or return an error?
		panic("Unauthorized: please check your SDK key")
	}

	result, err := ecr.configStore.GetConfiguration(experimentKey)

	return result, err
}

func (ecr *configurationRequestor) FetchAndStoreConfigurations() {
	result, err := ecr.httpClient.get(UFC_ENDPOINT)
	if err != nil {
		maskedErr := maskSensitiveInfo(err.Error())
		ecr.applicationLogger.Error("Failed to fetch UFC response", maskedErr)
		return
	}

	var wrapper ufcResponse
	err = json.Unmarshal([]byte(result), &wrapper)
	if err != nil {
		ecr.applicationLogger.Error("Failed to unmarshal UFC response JSON", result)
		ecr.applicationLogger.Error(err)
		return
	}

	// Now wrapper.Flags contains all configurations mapped by their keys
	// Pass this map directly to SetConfigurations
	err = ecr.configStore.SetConfigurations(wrapper.Flags)
	if err != nil {
		ecr.applicationLogger.Error("Failed to set configurations in configuration store", err)
	}
}
