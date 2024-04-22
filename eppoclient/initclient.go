// Package eppoclient provides client for eppo.cloud.
// Check InitClient to get started.
package eppoclient

import (
	"fmt"
	"net/http"
)

var __version__ = "2.0.0"

// InitClient is required to start polling of experiments configurations and create
// an instance of EppoClient, which could be used to get assignments information.
func InitClient(config Config) *EppoClient {
	if err := config.validate(); err != nil {
		panic(err)
	}
	sdkParams := SDKParams{apiKey: config.ApiKey, sdkName: "go", sdkVersion: __version__}

	httpClient := newHttpClient(config.BaseUrl, &http.Client{Timeout: REQUEST_TIMEOUT_SECONDS}, sdkParams)
	configStore, err := newConfigurationStore(MAX_CACHE_ENTRIES)
	if err != nil {
		panic(err)
	}
	requestor := newExperimentConfigurationRequestor(*httpClient, configStore)
	assignmentLogger := config.AssignmentLogger

	client := newEppoClient(requestor, assignmentLogger)

	client.poller.Start()

	return client
}

type ClientStopFn func()

// NewClient creates a new EppoClient instance and starts polling for experiment configurations.
// To stop polling and release resources, call the returned ClientStopFn.
func NewClient(config Config) (*EppoClient, ClientStopFn, error) {
	if err := config.validate(); err != nil {
		return nil, nil, fmt.Errorf("validate config: %v", err)
	}

	sdkParams := SDKParams{apiKey: config.ApiKey, sdkName: "go", sdkVersion: __version__}

	httpClient := newHttpClient(config.BaseUrl, &http.Client{Timeout: REQUEST_TIMEOUT_SECONDS}, sdkParams)
	configStore, err := newConfigurationStore(MAX_CACHE_ENTRIES)
	if err != nil {
		return nil, nil, fmt.Errorf("create configuration store: %v", err)
	}

	requestor := newExperimentConfigurationRequestor(*httpClient, configStore)
	requestor.FetchAndStoreConfigurations()

	assignmentLogger := config.AssignmentLogger
	client := newEppoClient(requestor, assignmentLogger)
	client.poller.Start()

	return client, ClientStopFn(client.poller.Stop), nil
}
