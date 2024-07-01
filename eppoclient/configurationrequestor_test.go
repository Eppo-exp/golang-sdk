package eppoclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_configurationRequestor_requestBandits(t *testing.T) {
	flags := readJsonFile[configResponse]("test-data/ufc/bandit-flags-v1.json")
	bandits := readJsonFile[banditResponse]("test-data/ufc/bandit-models-v1.json")
	server := newTestServer(flags, bandits)

	sdkParams := SDKParams{sdkKey: "blah", sdkName: "go", sdkVersion: __version__}
	httpClient := newHttpClient(server.URL, &http.Client{Timeout: REQUEST_TIMEOUT_SECONDS}, sdkParams)
	configurationStore := newConfigurationStore(configuration{})
	configurationRequestor := newConfigurationRequestor(*httpClient, configurationStore)

	configurationRequestor.FetchAndStoreConfigurations()

	config := configurationStore.getConfiguration()

	assert.NotEmpty(t, config.bandits.Bandits)
}

func Test_configurationRequestor_shouldNotRequestBanditsIfNotPresentInFlags(t *testing.T) {
	// flags-v1.json does not have a flag.Bandits field, so we
	// don't need to fetch bandits.
	flags := readJsonFile[configResponse]("test-data/ufc/flags-v1.json")
	bandits := readJsonFile[banditResponse]("test-data/ufc/bandit-models-v1.json")
	server := newTestServer(flags, bandits)

	sdkParams := SDKParams{sdkKey: "blah", sdkName: "go", sdkVersion: __version__}
	httpClient := newHttpClient(server.URL, &http.Client{Timeout: REQUEST_TIMEOUT_SECONDS}, sdkParams)
	configurationStore := newConfigurationStore(configuration{})
	configurationRequestor := newConfigurationRequestor(*httpClient, configurationStore)

	configurationRequestor.FetchAndStoreConfigurations()

	config := configurationStore.getConfiguration()

	assert.Empty(t, config.bandits.Bandits)
}

func newTestServer(configResponse configResponse, banditsResponse banditResponse) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/flag-config/v1/config":
			err := json.NewEncoder(w).Encode(configResponse)
			if err != nil {
				fmt.Println("Error encoding test response")
			}
		case "/flag-config/v1/bandits":
			err := json.NewEncoder(w).Encode(banditsResponse)
			if err != nil {
				fmt.Println("Error encoding test response")
			}
		default:
			http.NotFoundHandler().ServeHTTP(w, r)
		}
	}))
}
