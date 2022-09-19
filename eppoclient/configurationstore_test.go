package eppoclient

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testAllocationMap = make(map[string]Allocation)

var testExp = experimentConfiguration{
	SubjectShards: 1000,
	Enabled:       true,
	Allocations:   testAllocationMap,
	Rules:         []rule{},
	Name:          "randomization_algo",
}

const TEST_MAX_SIZE = 10

func Test_GetConfiguration_unknownKey(t *testing.T) {
	var store = newConfigurationStore(TEST_MAX_SIZE)
	store.SetConfigurations(dictionary{"randomization_algo": testExp})
	_, err := store.GetConfiguration("unknown_exp")

	assert.Error(t, err)
}

func Test_GetConfiguration_knownKey(t *testing.T) {
	var store = newConfigurationStore(TEST_MAX_SIZE)
	store.SetConfigurations(dictionary{"randomization_algo": testExp})
	result, _ := store.GetConfiguration("randomization_algo")

	expected := "randomization_algo"

	assert.Equal(t, expected, result.Name)
}

func Test_GetConfiguration_evictsOldEntriesWhenMaxSizeExceeded(t *testing.T) {
	var store = newConfigurationStore(TEST_MAX_SIZE)
	store.SetConfigurations(dictionary{"item_to_be_evicted": testExp})
	result, _ := store.GetConfiguration("item_to_be_evicted")

	expected := "randomization_algo"
	assert.Equal(t, expected, result.Name)

	for i := 0; i < TEST_MAX_SIZE; i++ {
		dictKey := fmt.Sprintf("test-entry-%v", i)
		store.SetConfigurations(dictionary{dictKey: testExp})
	}

	result, err := store.GetConfiguration("item_to_be_evicted")
	assert.Error(t, err)

	result, _ = store.GetConfiguration(fmt.Sprintf("test-entry-%v", TEST_MAX_SIZE-1))
	assert.Equal(t, expected, result.Name)
}
