package eppoclient

var __version__ = "1.0.0"

func InitClient(config Config) EppoClient {
	config.validate()
	sdkParams := SDKParams{apiKey: config.apiKey, sdkName: "go", sdkVersion: __version__}

	httpClient := HttpClient{}
	httpClient.New(config.baseUrl, sdkParams)

	configStore := ConfigurationStore{}
	configStore.New(MAX_CACHE_ENTRIES)

	requestor := NewExperimentConfigurationRequestor()
	requestor.New(httpClient, configStore)

	var assignmentLogger = NewAssignmentLogger()

	var client EppoClient
	client.New(&requestor, assignmentLogger)

	return client
}
