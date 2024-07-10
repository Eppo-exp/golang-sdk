package eppoclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockHttpClient struct {
	mock.Mock
}

func (m *mockHttpClient) get(url string) (HttpResponse, error) {
	args := m.Called(url)
	return args.Get(0).(HttpResponse), args.Error(1)
}

func Test_configurationRequestor_requestBandits(t *testing.T) {
	flags := readJsonFile[configResponse]("test-data/ufc/bandit-flags-v1.json")
	bandits := readJsonFile[banditResponse]("test-data/ufc/bandit-models-v1.json")
	server := newTestServer(flags, bandits)

	sdkParams := SDKParams{sdkKey: "blah", sdkName: "go", sdkVersion: __version__}
	httpClient := newHttpClient(server.URL, &http.Client{Timeout: REQUEST_TIMEOUT_SECONDS}, sdkParams)
	configurationStore := newConfigurationStore(configuration{})
	configurationRequestor := newConfigurationRequestor(httpClient, configurationStore, false)

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
	configurationRequestor := newConfigurationRequestor(httpClient, configurationStore, false)

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

func Test_FetchAndStoreConfigurations_SkipDeserializeIfUnchanged(t *testing.T) {
	mockHttpClient := new(mockHttpClient)
	mockConfigStore := newConfigurationStore(configuration{})
	configRequestor := newConfigurationRequestor(mockHttpClient, mockConfigStore, true) // true to skip deserialize

	// Mock the HTTP client to return a fixed response
	mockResponse1 := `
	{
		"createdAt": "2024-04-17T19:40:53.716Z",
		"environment": {
			"name": "Test"
		},
		"flags": {
			"empty_flag": {
				"key": "empty_flag",
				"enabled": true,
				"variationType": "STRING",
				"variations": {},
				"allocations": [],
				"totalShards": 10000
			}
		}
	}
	`
	mockCall := mockHttpClient.On("get", CONFIG_ENDPOINT).Return(HttpResponse{Body: mockResponse1, ETag: "tag_1"}, nil)

	// First fetch and store configurations
	configRequestor.FetchAndStoreConfigurations()

	// Fetch and store configurations again
	configRequestor.FetchAndStoreConfigurations()

	// Assert that configuration was fetched two times but deserialize was only called once
	mockHttpClient.AssertNumberOfCalls(t, "get", 2)
	assert.Equal(t, 1, configRequestor.deserializeCount)

	// Assert that configuration was stored as desired
	flag, err := mockConfigStore.getConfiguration().getFlagConfiguration("empty_flag")
	assert.Nil(t, err)
	assert.Equal(t, "empty_flag", flag.Key)
	mockCall.Unset()

	// change the remote config
	mockResponse2 := `
	{
		"createdAt": "2024-04-17T19:40:53.716Z",
		"environment": {
			"name": "Test"
		},
		"flags": {
			"empty_flag_2": {
				"key": "empty_flag_2",
				"enabled": true,
				"variationType": "STRING",
				"variations": {},
				"allocations": [],
				"totalShards": 10000
			}
		}
	}
	`
	mockCall = mockHttpClient.On("get", CONFIG_ENDPOINT).Return(HttpResponse{Body: mockResponse2, ETag: "tag_2"}, nil)

	// fetch and store again
	configRequestor.FetchAndStoreConfigurations()

	// assert that another fetch was called and deserialize was called
	mockHttpClient.AssertNumberOfCalls(t, "get", 3)
	assert.Equal(t, 2, configRequestor.deserializeCount)

	// assert that the new config is stored
	flag, err = mockConfigStore.getConfiguration().getFlagConfiguration("empty_flag")
	assert.NotNil(t, err)
	flag, err = mockConfigStore.getConfiguration().getFlagConfiguration("empty_flag_2")
	assert.Nil(t, err)
	assert.Equal(t, "empty_flag_2", flag.Key)
	mockCall.Unset()

	// change remote config back to original
	mockCall = mockHttpClient.On("get", CONFIG_ENDPOINT).Return(HttpResponse{Body: mockResponse1, ETag: "tag_1"}, nil)

	// fetch and store again
	configRequestor.FetchAndStoreConfigurations()

	// assert that another fetch was called and deserialize was called
	mockHttpClient.AssertNumberOfCalls(t, "get", 4)
	assert.Equal(t, 3, configRequestor.deserializeCount)

	flag, err = mockConfigStore.getConfiguration().getFlagConfiguration("empty_flag")
	assert.Nil(t, err)
	assert.Equal(t, "empty_flag", flag.Key)
	flag, err = mockConfigStore.getConfiguration().getFlagConfiguration("empty_flag_2")
	assert.NotNil(t, err)
	mockCall.Unset()
}
