// Package eppoclient provides client for eppo.cloud.
// Check InitClient to get started.
package eppoclient

import "net/http"

// InitClient is required to start polling of experiments configurations and create
// an instance of EppoClient, which could be used to get assignments information.
func InitClient(config Config) (*EppoClient, error) {
	err := config.validate()
	if err != nil {
		return nil, err
	}
	sdkParams := SDKParams{sdkKey: config.SdkKey, sdkName: "go", sdkVersion: __version__}
	applicationLogger := config.ApplicationLogger

	var httpClientInstance *http.Client
	if config.HttpClient != nil {
		httpClientInstance = config.HttpClient
	} else {
		httpClientInstance = &http.Client{Timeout: REQUEST_TIMEOUT_SECONDS}
	}
	httpClient := newHttpClient(config.BaseUrl, httpClientInstance, sdkParams)
	configStore := newConfigurationStore()
	requestor := newConfigurationRequestor(*httpClient, configStore, applicationLogger)

	poller := newPoller(config.PollerInterval, requestor.FetchAndStoreConfigurations, applicationLogger)
	client := newEppoClient(
		configStore,
		requestor,
		poller,
		config.AssignmentLogger,
		config.AssignmentLoggerContext,
		applicationLogger,
	)

	client.poller.Start()

	return client, nil
}
