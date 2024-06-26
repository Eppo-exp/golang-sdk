package eppoclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_AssignBlankExperiment(t *testing.T) {
	var mockConfigRequestor = new(mockConfigRequestor)
	var poller = newPoller(10, mockConfigRequestor.FetchAndStoreConfigurations)
	var mockLogger = new(mockLogger)
	client := newEppoClient(mockConfigRequestor, poller, mockLogger)

	_, err := client.GetStringAssignment("", "subject-1", Attributes{}, "")
	assert.Error(t, err)
}

func Test_AssignBlankSubject(t *testing.T) {
	var mockConfigRequestor = new(mockConfigRequestor)
	var poller = newPoller(10, mockConfigRequestor.FetchAndStoreConfigurations)
	var mockLogger = new(mockLogger)
	client := newEppoClient(mockConfigRequestor, poller, mockLogger)

	_, err := client.GetStringAssignment("experiment-1", "", Attributes{}, "")
	assert.Error(t, err)
}
func Test_LogAssignment(t *testing.T) {
	var mockLogger = new(mockLogger)
	mockLogger.Mock.On("LogAssignment", mock.Anything).Return()

	var mockConfigRequestor = new(mockConfigRequestor)
	var poller = newPoller(10, mockConfigRequestor.FetchAndStoreConfigurations)

	config := map[string]flagConfiguration{
		"experiment-key-1": flagConfiguration{
			Key:           "experiment-key-1",
			Enabled:       true,
			TotalShards:   10000,
			VariationType: stringVariation,
			Variations: map[string]variation{
				"control": variation{
					Key:   "control",
					Value: "control",
				},
			},
			Allocations: []allocation{
				{
					Key: "allocation-key",
					Splits: []split{
						{
							VariationKey: "control",
							Shards: []shard{
								{
									Salt: "",
									Ranges: []shardRange{
										{
											Start: 0,
											End:   10000,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	mockConfigRequestor.Mock.On("GetConfiguration", "experiment-key-1").Return(config["experiment-key-1"], nil)

	client := newEppoClient(mockConfigRequestor, poller, mockLogger)

	assignment, err := client.GetStringAssignment("experiment-key-1", "user-1", Attributes{}, "")
	expected := "control"

	assert.Nil(t, err)
	assert.Equal(t, expected, assignment)
	mockLogger.AssertNumberOfCalls(t, "LogAssignment", 1)
}

func Test_GetStringAssignmentHandlesLoggingPanic(t *testing.T) {
	var mockLogger = new(mockLogger)
	mockLogger.Mock.On("LogAssignment", mock.Anything).Panic("logging panic")

	var mockConfigRequestor = new(mockConfigRequestor)
	var poller = newPoller(10, mockConfigRequestor.FetchAndStoreConfigurations)

	config := map[string]flagConfiguration{
		"experiment-key-1": flagConfiguration{
			Key:           "experiment-key-1",
			Enabled:       true,
			TotalShards:   10000,
			VariationType: stringVariation,
			Variations: map[string]variation{
				"control": variation{
					Key:   "control",
					Value: "control",
				},
			},
			Allocations: []allocation{
				{
					Key: "allocation-key",
					Splits: []split{
						{
							VariationKey: "control",
							Shards: []shard{
								{
									Salt: "",
									Ranges: []shardRange{
										{
											Start: 0,
											End:   10000,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	mockConfigRequestor.Mock.On("GetConfiguration", "experiment-key-1").Return(config["experiment-key-1"], nil)

	client := newEppoClient(mockConfigRequestor, poller, mockLogger)

	assignment, err := client.GetStringAssignment("experiment-key-1", "user-1", Attributes{}, "")
	expected := "control"

	assert.Nil(t, err)
	assert.Equal(t, expected, assignment)
}
