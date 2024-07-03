package eppoclient

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

const UFC_ENDPOINT = "/flag-config/v1/config"

type iConfigRequestor interface {
	GetConfiguration(key string) (flagConfiguration, error)
	FetchAndStoreConfigurations()
}

type configurationRequestor struct {
	httpClient            HttpClientInterface
	configStore           *configurationStore
	storedFlagConfigsHash string
}

func newConfigurationRequestor(httpClient HttpClientInterface, configStore *configurationStore) *configurationRequestor {
	return &configurationRequestor{
		httpClient:  httpClient,
		configStore: configStore,
	}
}

func (ecr *configurationRequestor) GetConfiguration(experimentKey string) (flagConfiguration, error) {
	result, err := ecr.configStore.GetConfiguration(experimentKey)

	return result, err
}

func (ecr *configurationRequestor) FetchAndStoreConfigurations() {
	result, err := ecr.httpClient.get(UFC_ENDPOINT)
	if err != nil {
		fmt.Println("Failed to fetch UFC response", err)
		return
	}

	// Calculate the hash of the current response
	hash := sha256.New()
	hash.Write([]byte(result))
	receivedFlagConfigsHash := hex.EncodeToString(hash.Sum(nil))

	// Compare the current hash with the last saved hash
	if receivedFlagConfigsHash == ecr.storedFlagConfigsHash {
		fmt.Println("[EppoSDK] Response has not changed, skipping deserialization and cache update.")
		return
	}

	ecr.storedFlagConfigsHash = receivedFlagConfigsHash

	var wrapper ufcResponse
	err = json.Unmarshal([]byte(result), &wrapper)
	if err != nil {
		fmt.Println("Failed to unmarshal UFC response JSON", result)
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
