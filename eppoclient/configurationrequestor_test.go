package eppoclient

import "testing"

// todo remove

func Test_FetchAndStoreConfigurations(t *testing.T) {
	var requestor = ExperimentConfigurationRequestor{}

	var httpClient = HttpClient{}
	var sdkParams = SDKParams{apiKey: "tgcwcyYqosYfRpA5V3khTnsH8o2MlauhxSTyst6mDUM", sdkName: "", sdkVersion: ""}

	httpClient.New("http://localhost:4000/api", sdkParams)

	var configStore = ConfigurationStore{}
	configStore.New(20)

	requestor.New(httpClient, configStore)

	requestor.FetchAndStoreConfigurations()
}
