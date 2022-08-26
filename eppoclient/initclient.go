// Package eppoclient provides client for eppo.cloud.
// Check InitClient to get started.
package eppoclient

import "net/http"

var __version__ = "1.0.0"

// InitClient is required to start polling of experiments configurations and create
// an instance of EppoClient, which could be used to get assignments information.
func InitClient(config Config) *EppoClient {
	config.validate()
	sdkParams := SDKParams{apiKey: config.ApiKey, sdkName: "go", sdkVersion: __version__}

	httpClient := NewHttpClient(config.BaseUrl, &http.Client{Timeout: REQUEST_TIMEOUT_SECONDS}, sdkParams)
	configStore := NewConfigurationStore(MAX_CACHE_ENTRIES)
	requestor := NewExperimentConfigurationRequestor(*httpClient, *configStore)
	assignmentLogger := NewAssignmentLogger()

	client := NewEppoClient(requestor, assignmentLogger)

	client.poller.Start()

	return client
}
