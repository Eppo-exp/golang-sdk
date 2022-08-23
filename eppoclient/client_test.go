package eppoclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_AssignBlankExperiment(t *testing.T) {
	var mockConfigRequestor = new(MockConfigRequestor)
	var mockLogger = new(MockLogger)
	client := NewEppoClient(mockConfigRequestor, mockLogger)

	assert.Panics(t, func() { client.GetAssignment("subject-1", "", Dictionary{}) })
}

func Test_AssignBlankSubject(t *testing.T) {
	var mockConfigRequestor = new(MockConfigRequestor)
	var mockLogger = new(MockLogger)
	client := NewEppoClient(mockConfigRequestor, mockLogger)

	assert.Panics(t, func() { client.GetAssignment("", "experiment-1", Dictionary{}) })
}

func Test_SubjectNotInSample(t *testing.T) {
	var mockLogger = new(MockLogger)
	var mockConfigRequestor = new(MockConfigRequestor)
	overrides := make(Dictionary)
	var mockVariations = []Variation{
		{Name: "control", ShardRange: ShardRange{Start: 0, End: 10000}},
	}
	mockResult := ExperimentConfiguration{
		Name:            "recommendation_algo",
		PercentExposure: 0,
		Enabled:         true,
		SubjectShards:   1000,
		Overrides:       overrides,
		Variations:      mockVariations,
	}

	mockConfigRequestor.Mock.On("GetConfiguration", mock.Anything).Return(mockResult, nil)

	client := NewEppoClient(mockConfigRequestor, mockLogger)

	assignment, _ := client.GetAssignment("user-1", "experiment-key-1", Dictionary{})

	assert.Equal(t, "", assignment)
}

func Test_LogAssignment(t *testing.T) {
	var mockLogger = new(MockLogger)
	mockLogger.Mock.On("LogAssignment", mock.AnythingOfType("string")).Return()

	var mockConfigRequestor = new(MockConfigRequestor)
	overrides := make(Dictionary)

	var mockVariations = []Variation{
		{Name: "control", ShardRange: ShardRange{Start: 0, End: 10000}},
	}
	mockResult := ExperimentConfiguration{
		Name:            "recommendation_algo",
		PercentExposure: 100,
		Enabled:         true,
		SubjectShards:   1000,
		Overrides:       overrides,
		Variations:      mockVariations,
	}
	mockConfigRequestor.Mock.On("GetConfiguration", "experiment-key-1").Return(mockResult, nil)

	client := NewEppoClient(mockConfigRequestor, mockLogger)

	assignment, err := client.GetAssignment("user-1", "experiment-key-1", Dictionary{})
	expected := "control"

	if err != nil {
		t.Errorf("\"EppoClient.GetAssignment()\" FAILED, expected -> %v, got -> %v", expected, assignment)
	}

	assert.Equal(t, expected, assignment)
	mockLogger.AssertNumberOfCalls(t, "LogAssignment", 1)
}
