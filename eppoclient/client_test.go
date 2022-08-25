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

	assignment, err := client.GetAssignment("user-1", "experiment-key-1", Dictionary{})

	assert.Equal(t, "", assignment)
	assert.NotNil(t, err)
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

	assert.Nil(t, err)
	assert.Equal(t, expected, assignment)
	mockLogger.AssertNumberOfCalls(t, "LogAssignment", 1)
}

func Test_GetAssignmentHandlesLoggingPanic(t *testing.T) {
	var mockLogger = new(MockLogger)
	mockLogger.Mock.On("LogAssignment", mock.AnythingOfType("string")).Panic("logging panic")

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

	assert.Nil(t, err)
	assert.Equal(t, expected, assignment)
}

func Test_AssignSubjectWithAttributesAndRules(t *testing.T) {
	var mockLogger = new(MockLogger)
	mockLogger.Mock.On("LogAssignment", mock.AnythingOfType("string")).Return()

	var matchesEmailCondition = Condition{operator: "MATCHES", value: ".*@eppo.com", attribute: "email"}
	var textRule = Rule{conditions: []Condition{matchesEmailCondition}}
	var mockConfigRequestor = new(MockConfigRequestor)
	var overrides = make(Dictionary)
	var mockVariations = []Variation{
		{Name: "control", ShardRange: ShardRange{Start: 0, End: 10000}},
	}
	var mockResult = ExperimentConfiguration{
		Name:            "recommendation_algo",
		PercentExposure: 100,
		Enabled:         true,
		SubjectShards:   1000,
		Overrides:       overrides,
		Variations:      mockVariations,
		Rules:           []Rule{textRule},
	}
	mockConfigRequestor.Mock.On("GetConfiguration", "experiment-key-1").Return(mockResult, nil)

	client := NewEppoClient(mockConfigRequestor, mockLogger)

	expected := ""
	assignment, _ := client.GetAssignment("user-1", "experiment-key-1", Dictionary{})
	assert.Equal(t, expected, assignment)

	assignment, _ = client.GetAssignment("user-1", "experiment-key-1", Dictionary{
		"email": "test@example.com",
	})
	assert.Equal(t, expected, assignment)

	expected = "control"
	assignment, _ = client.GetAssignment("user-1", "experiment-key-1", Dictionary{
		"email": "test@eppo.com",
	})
	assert.Equal(t, expected, assignment)
}

func Test_WithSubjectInOverrides(t *testing.T) {
	var mockLogger = new(MockLogger)
	mockLogger.Mock.On("LogAssignment", mock.AnythingOfType("string")).Return()

	var mockConfigRequestor = new(MockConfigRequestor)
	var mockVariations = []Variation{
		{Name: "control", ShardRange: ShardRange{Start: 0, End: 100}},
	}
	var overrides = make(Dictionary)
	overrides["d6d7705392bc7af633328bea8c4c6904"] = "override-variation"
	var mockResult = ExperimentConfiguration{
		Name:            "recommendation_algo",
		PercentExposure: 100,
		Enabled:         true,
		SubjectShards:   1000,
		Overrides:       overrides,
		Variations:      mockVariations,
		Rules:           []Rule{textRule},
	}

	mockConfigRequestor.Mock.On("GetConfiguration", "experiment-key-1").Return(mockResult, nil)

	client := NewEppoClient(mockConfigRequestor, mockLogger)

	expected := "override-variation"
	assignment, _ := client.GetAssignment("user-1", "experiment-key-1", Dictionary{})
	assert.Equal(t, expected, assignment)
}

func Test_WithSubjectInOverridesExpDisabled(t *testing.T) {
	var mockLogger = new(MockLogger)
	mockLogger.Mock.On("LogAssignment", mock.AnythingOfType("string")).Return()

	var mockConfigRequestor = new(MockConfigRequestor)
	var mockVariations = []Variation{
		{Name: "control", ShardRange: ShardRange{Start: 0, End: 100}},
	}
	var overrides = make(Dictionary)
	overrides["d6d7705392bc7af633328bea8c4c6904"] = "override-variation"
	var mockResult = ExperimentConfiguration{
		Name:            "recommendation_algo",
		PercentExposure: 100,
		Enabled:         false,
		SubjectShards:   1000,
		Overrides:       overrides,
		Variations:      mockVariations,
		Rules:           []Rule{textRule},
	}

	mockConfigRequestor.Mock.On("GetConfiguration", "experiment-key-1").Return(mockResult, nil)

	client := NewEppoClient(mockConfigRequestor, mockLogger)

	expected := "override-variation"
	assignment, err := client.GetAssignment("user-1", "experiment-key-1", Dictionary{})

	assert.Nil(t, err)
	assert.Equal(t, expected, assignment)
}

func Test_WithNullExpConfig(t *testing.T) {
	var mockLogger = new(MockLogger)
	mockLogger.Mock.On("LogAssignment", mock.AnythingOfType("string")).Return()

	var mockConfigRequestor = new(MockConfigRequestor)
	mockConfigRequestor.Mock.On("GetConfiguration", "experiment-key-1").Return(ExperimentConfiguration{}, nil)

	client := NewEppoClient(mockConfigRequestor, mockLogger)

	expected := ""
	assignment, err := client.GetAssignment("user-1", "experiment-key-1", Dictionary{})

	assert.NotNil(t, err)
	assert.Equal(t, expected, assignment)
}
