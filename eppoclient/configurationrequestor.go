package eppoclient

import (
	"encoding/json"
	"fmt"
)

const UFC_ENDPOINT = "/flag-config/v1/config"

type iConfigRequestor interface {
	GetConfiguration(key string) (flagConfiguration, error)
	FetchAndStoreConfigurations()
}

type configurationRequestor struct {
	httpClient                                    HttpClientInterface
	configStore                                   *configurationStore
	storedUFCResponseETag                         string
	deserializeCount                              int
	skipDeserializeAndUpdateFlagConfigIfUnchanged bool
}

func newConfigurationRequestor(httpClient HttpClientInterface, configStore *configurationStore, skipDeserializeAndUpdateFlagConfigIfUnchanged bool) *configurationRequestor {
	return &configurationRequestor{
		httpClient:       httpClient,
		configStore:      configStore,
		deserializeCount: 0,
		skipDeserializeAndUpdateFlagConfigIfUnchanged: skipDeserializeAndUpdateFlagConfigIfUnchanged,
	}
}

func (ecr *configurationRequestor) GetConfiguration(experimentKey string) (flagConfiguration, error) {
	result, err := ecr.configStore.GetConfiguration(experimentKey)

	return result, err
}

func (ecr *configurationRequestor) FetchAndStoreConfigurations() {
	httpResponse, err := ecr.httpClient.get(UFC_ENDPOINT)
	if err != nil {
		fmt.Println("Failed to fetch UFC response", err)
		return
	}

	if ecr.skipDeserializeAndUpdateFlagConfigIfUnchanged {
		// Compare the current hash with the last saved hash
		if httpResponse.ETag == ecr.storedUFCResponseETag {
			fmt.Println("[EppoSDK] Response has not changed, skipping deserialization and cache update.")
			return
		}

		// Update the stored hash
		ecr.storedUFCResponseETag = httpResponse.ETag
	}

	var wrapper ufcResponse
	err = json.Unmarshal([]byte(httpResponse.Body), &wrapper)
	if err != nil {
		fmt.Println("Failed to unmarshal UFC response JSON", httpResponse.Body)
		fmt.Println(err)
		return
	}

	// Now wrapper.Flags contains all configurations mapped by their keys
	// Pass this map directly to SetConfigurations
	err = ecr.configStore.SetConfigurations(wrapper.Flags)
	if err != nil {
		fmt.Println("Failed to set configurations in configuration store", err)
	}
}
