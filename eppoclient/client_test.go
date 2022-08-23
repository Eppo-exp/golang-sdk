package eppoclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AssignBlankExperiment(t *testing.T) {
	var mockConfigRequestor = &MockConfigRequestor{}

	client := NewEppoClient(mockConfigRequestor, mockLogger)

	assert.Panics(t, func() { client.GetAssignment("subject-1", "", Dictionary{}) })
}

func Test_AssignBlankSubject(t *testing.T) {
	var mockConfigRequestor = &MockConfigRequestor{}

	client := NewEppoClient(mockConfigRequestor, mockLogger)

	assert.Panics(t, func() { client.GetAssignment("", "experiment-1", Dictionary{}) })
}

func Test_SubjectNotInSample(t *testing.T) {
	var mockConfigRequestor = &MockConfigRequestor{}

	client := NewEppoClient(mockConfigRequestor, mockLogger)

	assignment, _ := client.GetAssignment("user-1", "experiment-key-1", Dictionary{})

	assert.Equal(t, "", assignment)
}

func Test_LogAssignment(t *testing.T) {
	var mockConfigRequestor = &MockConfigRequestor100PercentExposure{}

	client := NewEppoClient(mockConfigRequestor, mockLogger)

	assignment, err := client.GetAssignment("user-1", "experiment-key-1", Dictionary{})
	expected := "control"

	if err != nil {
		t.Errorf("\"EppoClient.GetAssignment()\" FAILED, expected -> %v, got -> %v", expected, assignment)
	}

	assert.Equal(t, expected, assignment)
}
