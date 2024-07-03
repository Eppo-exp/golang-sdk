package eppoclient

import (
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

const flagConfig_1 = `
{
	"createdAt": "2024-04-17T19:40:53.716Z",
	"environment": {
		"name": "Test"
	},
	"flags": {
		"empty_flag_1": {
			"key": "empty_flag_1",
			"enabled": true,
			"variationType": "STRING",
			"variations": {},
			"allocations": [],
			"totalShards": 10000
		}
	}
}
`

const flagConfig_2 = `
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

func Test_FetchAndStoreConfigurations_SkipDeserializeAndUpdateFlagConfigIfUnchanged(t *testing.T) {
	mockHttpClient := new(mockHttpClient)
	mockConfigStore := newConfigurationStore()
	configRequestor := newConfigurationRequestor(mockHttpClient, mockConfigStore, true)

	// Mock the HTTP client to return a fixed response
	mockResponse1 := HttpResponse{
		Body: flagConfig_1,
		ETag: "tag_1",
	}
	mockCall := mockHttpClient.On("get", UFC_ENDPOINT).Return(mockResponse1, nil)

	// First fetch and store configurations
	configRequestor.FetchAndStoreConfigurations()

	// Fetch and store configurations again
	configRequestor.FetchAndStoreConfigurations()

	// Assert that the configurations were only deserialized once
	assert.Equal(t, 1, len(mockConfigStore.configs))
	assert.Equal(t, "empty_flag_1", mockConfigStore.configs["empty_flag_1"].Key)
	mockHttpClient.AssertNumberOfCalls(t, "get", 2)
	mockCall.Unset()

	// change the remote config
	mockResponse2 := HttpResponse{
		Body: flagConfig_2,
		ETag: "tag_2",
	}
	mockCall = mockHttpClient.On("get", UFC_ENDPOINT).Return(mockResponse2, nil)

	// fetch and store again
	configRequestor.FetchAndStoreConfigurations()

	// assert that the new config is stored
	assert.Equal(t, 1, len(mockConfigStore.configs))
	assert.Equal(t, "empty_flag_2", mockConfigStore.configs["empty_flag_2"].Key)
	mockHttpClient.AssertNumberOfCalls(t, "get", 3)
	mockCall.Unset()

	// change remote config back to original
	mockCall = mockHttpClient.On("get", UFC_ENDPOINT).Return(mockResponse1, nil)

	// fetch and store again
	configRequestor.FetchAndStoreConfigurations()

	assert.Equal(t, 1, len(mockConfigStore.configs))
	assert.Equal(t, "empty_flag_1", mockConfigStore.configs["empty_flag_1"].Key)
	mockHttpClient.AssertNumberOfCalls(t, "get", 4)
	mockCall.Unset()
}

func Test_FetchAndStoreConfigurations_AlwaysDeserializeAndUpdateFlagConfig(t *testing.T) {
	mockHttpClient := new(mockHttpClient)
	mockConfigStore := newConfigurationStore()
	configRequestor := newConfigurationRequestor(mockHttpClient, mockConfigStore, false)

	// Mock the HTTP client to return a fixed response
	mockResponse1 := HttpResponse{
		Body: flagConfig_1,
		ETag: "tag_1",
	}
	mockCall := mockHttpClient.On("get", UFC_ENDPOINT).Return(mockResponse1, nil)

	// First fetch and store configurations
	configRequestor.FetchAndStoreConfigurations()

	// Fetch and store configurations again
	configRequestor.FetchAndStoreConfigurations()

	// Assert that the configurations were only deserialized once
	assert.Equal(t, 1, len(mockConfigStore.configs))
	assert.Equal(t, "empty_flag_1", mockConfigStore.configs["empty_flag_1"].Key)
	mockHttpClient.AssertNumberOfCalls(t, "get", 2)
	mockCall.Unset()

	// change the remote config
	mockResponse2 := HttpResponse{
		Body: flagConfig_2,
		ETag: "tag_2",
	}
	mockCall = mockHttpClient.On("get", UFC_ENDPOINT).Return(mockResponse2, nil)

	// fetch and store again
	configRequestor.FetchAndStoreConfigurations()

	// assert that the new config is stored
	assert.Equal(t, 1, len(mockConfigStore.configs))
	assert.Equal(t, "empty_flag_2", mockConfigStore.configs["empty_flag_2"].Key)
	mockHttpClient.AssertNumberOfCalls(t, "get", 3)
	mockCall.Unset()
}
