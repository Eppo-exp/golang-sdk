package eppoclient

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testExp = ExperimentConfiguration{
	SubjectShards:   1000,
	PercentExposure: 1,
	Enabled:         true,
	Variations:      []Variation{},
	Name:            "randomization_algo",
}

const TEST_MAX_SIZE = 10

func Test_GetConfiguration_unknownKey(t *testing.T) {
	var store = NewConfigurationStore(TEST_MAX_SIZE)
	store.SetConfigurations(Dictionary{"randomization_algo": testExp})
	_, err := store.GetConfiguration("unknown_exp")

	assert.Error(t, err)
}

func Test_GetConfiguration_knownKey(t *testing.T) {
	var store = NewConfigurationStore(TEST_MAX_SIZE)
	store.SetConfigurations(Dictionary{"randomization_algo": testExp})
	result, _ := store.GetConfiguration("randomization_algo")

	expected := "randomization_algo"

	assert.Equal(t, expected, result.Name)
}

func Test_GetConfiguration_evictsOldEntriesWhenMaxSizeExceeded(t *testing.T) {
	var store = NewConfigurationStore(TEST_MAX_SIZE)
	store.SetConfigurations(Dictionary{"item_to_be_evicted": testExp})
	result, _ := store.GetConfiguration("item_to_be_evicted")

	expected := "randomization_algo"
	assert.Equal(t, expected, result.Name)

	for i := 0; i < TEST_MAX_SIZE; i++ {
		dictKey := fmt.Sprintf("test-entry-%v", i)
		store.SetConfigurations(Dictionary{dictKey: testExp})
	}

	result, err := store.GetConfiguration("item_to_be_evicted")
	assert.Error(t, err)

	result, _ = store.GetConfiguration(fmt.Sprintf("test-entry-%v", TEST_MAX_SIZE-1))
	assert.Equal(t, expected, result.Name)
}
