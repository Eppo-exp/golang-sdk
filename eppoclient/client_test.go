package eppoclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_AssignBlankExperiment(t *testing.T) {
	var mockConfigRequestor = new(mockConfigRequestor)
	var mockLogger = new(mockLogger)
	client := newEppoClient(mockConfigRequestor, mockLogger)

	assert.Panics(t, func() {
		_, err := client.GetAssignment("subject-1", "", dictionary{})

		if err == nil {
			t.Error("Test failed")
		}
	})
}

func Test_AssignBlankSubject(t *testing.T) {
	var mockConfigRequestor = new(mockConfigRequestor)
	var mockLogger = new(mockLogger)
	client := newEppoClient(mockConfigRequestor, mockLogger)

	assert.Panics(t, func() {
		_, err := client.GetAssignment("", "experiment-1", dictionary{})

		if err == nil {
			t.Error("Test failed")
		}
	})
}

func Test_SubjectNotInSample(t *testing.T) {
	var mockLogger = new(mockLogger)
	var mockConfigRequestor = new(mockConfigRequestor)
	overrides := make(dictionary)
	var mockVariations = []Variation{
		{Name: "control", ShardRange: shardRange{Start: 0, End: 10000}},
	}
	mockResult := experimentConfiguration{
		Name:            "recommendation_algo",
		PercentExposure: 0,
		Enabled:         true,
		SubjectShards:   1000,
		Overrides:       overrides,
		Variations:      mockVariations,
	}

	mockConfigRequestor.Mock.On("GetConfiguration", mock.Anything).Return(mockResult, nil)

	client := newEppoClient(mockConfigRequestor, mockLogger)

	assignment, err := client.GetAssignment("user-1", "experiment-key-1", dictionary{})

	assert.Equal(t, "", assignment)
	assert.NotNil(t, err)
}

func Test_LogAssignment(t *testing.T) {
	var mockLogger = new(mockLogger)
	mockLogger.Mock.On("LogAssignment", mock.Anything).Return()

	var mockConfigRequestor = new(mockConfigRequestor)
	overrides := make(dictionary)

	var mockVariations = []Variation{
		{Name: "control", ShardRange: shardRange{Start: 0, End: 10000}},
	}
	mockResult := experimentConfiguration{
		Name:            "recommendation_algo",
		PercentExposure: 100,
		Enabled:         true,
		SubjectShards:   1000,
		Overrides:       overrides,
		Variations:      mockVariations,
	}
	mockConfigRequestor.Mock.On("GetConfiguration", "experiment-key-1").Return(mockResult, nil)

	client := newEppoClient(mockConfigRequestor, mockLogger)

	assignment, err := client.GetAssignment("user-1", "experiment-key-1", dictionary{})
	expected := "control"

	assert.Nil(t, err)
	assert.Equal(t, expected, assignment)
	mockLogger.AssertNumberOfCalls(t, "LogAssignment", 1)
}

func Test_GetAssignmentHandlesLoggingPanic(t *testing.T) {
	var mockLogger = new(mockLogger)
	mockLogger.Mock.On("LogAssignment", mock.Anything).Panic("logging panic")

	var mockConfigRequestor = new(mockConfigRequestor)
	overrides := make(dictionary)

	var mockVariations = []Variation{
		{Name: "control", ShardRange: shardRange{Start: 0, End: 10000}},
	}
	mockResult := experimentConfiguration{
		Name:            "recommendation_algo",
		PercentExposure: 100,
		Enabled:         true,
		SubjectShards:   1000,
		Overrides:       overrides,
		Variations:      mockVariations,
	}
	mockConfigRequestor.Mock.On("GetConfiguration", "experiment-key-1").Return(mockResult, nil)

	client := newEppoClient(mockConfigRequestor, mockLogger)

	assignment, err := client.GetAssignment("user-1", "experiment-key-1", dictionary{})
	expected := "control"

	assert.Nil(t, err)
	assert.Equal(t, expected, assignment)
}

func Test_AssignSubjectWithAttributesAndRules(t *testing.T) {
	var mockLogger = new(mockLogger)
	mockLogger.Mock.On("LogAssignment", mock.Anything).Return()

	var matchesEmailCondition = condition{Operator: "MATCHES", Value: ".*@eppo.com", Attribute: "email"}
	var textRule = rule{Conditions: []condition{matchesEmailCondition}}
	var mockConfigRequestor = new(mockConfigRequestor)
	var overrides = make(dictionary)
	var mockVariations = []Variation{
		{Name: "control", ShardRange: shardRange{Start: 0, End: 10000}},
	}
	var mockResult = experimentConfiguration{
		Name:            "recommendation_algo",
		PercentExposure: 100,
		Enabled:         true,
		SubjectShards:   1000,
		Overrides:       overrides,
		Variations:      mockVariations,
		Rules:           []rule{textRule},
	}
	mockConfigRequestor.Mock.On("GetConfiguration", "experiment-key-1").Return(mockResult, nil)

	tests := []struct {
		a    string
		b    string
		c    dictionary
		want string
	}{
		{"user-1", "experiment-key-1", dictionary{}, ""},
		{"user-1", "experiment-key-1", dictionary{
			"email": "test@example.com",
		}, ""},
		{"user-1", "experiment-key-1", dictionary{
			"email": "test@eppo.com",
		}, "control"},
	}

	client := newEppoClient(mockConfigRequestor, mockLogger)

	for _, tt := range tests {
		assignment, _ := client.GetAssignment(tt.a, tt.b, tt.c)

		assert.Equal(t, tt.want, assignment)
	}
}

func Test_WithSubjectInOverrides(t *testing.T) {
	var mockLogger = new(mockLogger)
	mockLogger.Mock.On("LogAssignment", mock.Anything).Return()

	var mockConfigRequestor = new(mockConfigRequestor)
	var mockVariations = []Variation{
		{Name: "control", ShardRange: shardRange{Start: 0, End: 100}},
	}
	var overrides = make(dictionary)
	overrides["d6d7705392bc7af633328bea8c4c6904"] = "override-variation"
	var mockResult = experimentConfiguration{
		Name:            "recommendation_algo",
		PercentExposure: 100,
		Enabled:         true,
		SubjectShards:   1000,
		Overrides:       overrides,
		Variations:      mockVariations,
		Rules:           []rule{textRule},
	}

	mockConfigRequestor.Mock.On("GetConfiguration", "experiment-key-1").Return(mockResult, nil)

	client := newEppoClient(mockConfigRequestor, mockLogger)

	expected := "override-variation"
	assignment, _ := client.GetAssignment("user-1", "experiment-key-1", dictionary{})
	assert.Equal(t, expected, assignment)
}

func Test_WithSubjectInOverridesExpDisabled(t *testing.T) {
	var mockLogger = new(mockLogger)
	mockLogger.Mock.On("LogAssignment", mock.Anything).Return()

	var mockConfigRequestor = new(mockConfigRequestor)
	var mockVariations = []Variation{
		{Name: "control", ShardRange: shardRange{Start: 0, End: 100}},
	}
	var overrides = make(dictionary)
	overrides["d6d7705392bc7af633328bea8c4c6904"] = "override-variation"
	var mockResult = experimentConfiguration{
		Name:            "recommendation_algo",
		PercentExposure: 100,
		Enabled:         false,
		SubjectShards:   1000,
		Overrides:       overrides,
		Variations:      mockVariations,
		Rules:           []rule{textRule},
	}

	mockConfigRequestor.Mock.On("GetConfiguration", "experiment-key-1").Return(mockResult, nil)

	client := newEppoClient(mockConfigRequestor, mockLogger)

	expected := "override-variation"
	assignment, err := client.GetAssignment("user-1", "experiment-key-1", dictionary{})

	assert.Nil(t, err)
	assert.Equal(t, expected, assignment)
}

func Test_WithNullExpConfig(t *testing.T) {
	var mockLogger = new(mockLogger)
	mockLogger.Mock.On("LogAssignment", mock.Anything).Return()

	var mockConfigRequestor = new(mockConfigRequestor)
	mockConfigRequestor.Mock.On("GetConfiguration", "experiment-key-1").Return(experimentConfiguration{}, nil)

	client := newEppoClient(mockConfigRequestor, mockLogger)

	expected := ""
	assignment, err := client.GetAssignment("user-1", "experiment-key-1", dictionary{})

	assert.NotNil(t, err)
	assert.Equal(t, expected, assignment)
}
