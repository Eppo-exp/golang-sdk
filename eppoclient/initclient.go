package eppoclient

import "net/http"

var __version__ = "1.0.0"

func InitClient(config Config) *EppoClient {
	config.validate()
	sdkParams := SDKParams{apiKey: config.apiKey, sdkName: "go", sdkVersion: __version__}

	httpClient := NewHttpClient(config.baseUrl, &http.Client{Timeout: REQUEST_TIMEOUT_SECONDS}, sdkParams)
	configStore := NewConfigurationStore(MAX_CACHE_ENTRIES)
	requestor := NewExperimentConfigurationRequestor(*httpClient, *configStore)
	assignmentLogger := NewAssignmentLogger()

	client := NewEppoClient(requestor, assignmentLogger)

	client.poller.Start()

	return client
}
