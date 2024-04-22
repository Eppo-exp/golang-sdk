package eppoclient

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	store, err := newConfigurationStore(TEST_MAX_SIZE)
	require.NoError(t, err, "Failed to create configuration store")

	err = store.SetConfigurations(map[string]experimentConfiguration{
		"randomization_algo": testExp,
	})

	assert.NoError(t, err)
	result, err := store.GetConfiguration("unknown_exp")

	assert.Error(t, err)
	assert.Equal(t, experimentConfiguration{}, result)
}

func Test_GetConfiguration_knownKey(t *testing.T) {
	store, err := newConfigurationStore(TEST_MAX_SIZE)
	require.NoError(t, err, "Failed to create configuration store")

	err = store.SetConfigurations(map[string]experimentConfiguration{
		"randomization_algo": testExp,
	})
	assert.NoError(t, err)
	result, err := store.GetConfiguration("randomization_algo")

	expected := "randomization_algo"

	assert.NoError(t, err)
	assert.Equal(t, expected, result.Name)
}

func Test_GetConfiguration_evictsOldEntriesWhenMaxSizeExceeded(t *testing.T) {
	store, err := newConfigurationStore(TEST_MAX_SIZE)
	require.NoError(t, err, "Failed to create configuration store")

	err = store.SetConfigurations(map[string]experimentConfiguration{
		"item_to_be_evicted": testExp,
	})
	assert.NoError(t, err)
	result, err := store.GetConfiguration("item_to_be_evicted")

	expected := "randomization_algo"
	assert.NoError(t, err)
	assert.Equal(t, expected, result.Name)

	for i := 0; i < TEST_MAX_SIZE; i++ {
		dictKey := fmt.Sprintf("test-entry-%v", i)
		err := store.SetConfigurations(map[string]experimentConfiguration{
			dictKey: testExp,
		})
		assert.NoError(t, err)
	}

	result, err = store.GetConfiguration("item_to_be_evicted")
	assert.Error(t, err)

	result, err = store.GetConfiguration(fmt.Sprintf("test-entry-%v", TEST_MAX_SIZE-1))
	assert.NoError(t, err)
	assert.Equal(t, expected, result.Name)
}
