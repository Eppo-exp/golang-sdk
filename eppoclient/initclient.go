// Package eppoclient provides client for eppo.cloud.
// Check InitClient to get started.
package eppoclient

import "net/http"

var __version__ = "4.0.1"

// InitClient is required to start polling of experiments configurations and create
// an instance of EppoClient, which could be used to get assignments information.
func InitClient(config Config) (*EppoClient, error) {
	err := config.validate()
	if err != nil {
		return nil, err
	}
	sdkParams := SDKParams{sdkKey: config.SdkKey, sdkName: "go", sdkVersion: __version__}

	httpClient := newHttpClient(config.BaseUrl, &http.Client{Timeout: REQUEST_TIMEOUT_SECONDS}, sdkParams)
	configStore := newConfigurationStore()
	requestor := newConfigurationRequestor(*httpClient, configStore)
	assignmentLogger := config.AssignmentLogger

	poller := newPoller(config.PollerInterval, requestor.FetchAndStoreConfigurations)

	client := newEppoClient(requestor, poller, assignmentLogger)

	client.poller.Start()

	return client, nil
}
