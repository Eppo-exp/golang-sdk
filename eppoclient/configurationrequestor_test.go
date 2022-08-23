package eppoclient

import "testing"

// todo remove

func Test_FetchAndStoreConfigurations(t *testing.T) {
	var sdkParams = SDKParams{apiKey: "tgcwcyYqosYfRpA5V3khTnsH8o2MlauhxSTyst6mDUM", sdkName: "", sdkVersion: ""}

	var httpClient = NewHttpClient("http://localhost:4000/api", sdkParams)
	var configStore = NewConfigurationStore(20)
	var requestor = NewExperimentConfigurationRequestor(*httpClient, *configStore)

	requestor.FetchAndStoreConfigurations()
}
