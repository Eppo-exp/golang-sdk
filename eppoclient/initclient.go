package eppoclient

var __version__ = "1.0.0"

func InitClient(config Config) EppoClient {
	config.validate()
	sdkParams := SDKParams{apiKey: config.apiKey, sdkName: "go", sdkVersion: __version__}

	httpClient := NewHttpClient(config.baseUrl, sdkParams)
	configStore := NewConfigurationStore(MAX_CACHE_ENTRIES)
	requestor := NewExperimentConfigurationRequestor(*httpClient, *configStore)
	assignmentLogger := NewAssignmentLogger()

	var client EppoClient
	client.New(requestor, assignmentLogger)

	return client
}
