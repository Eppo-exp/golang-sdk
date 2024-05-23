package eppoclient

import (
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

func Test_GetConfiguration_unknownKey(t *testing.T) {
	var store = newConfigurationStore()
	err := store.SetConfigurations(map[string]experimentConfiguration{
		"randomization_algo": testExp,
	})

	assert.NoError(t, err)
	result, err := store.GetConfiguration("unknown_exp")

	assert.Error(t, err)
	assert.Equal(t, experimentConfiguration{}, result)
}

func Test_GetConfiguration_knownKey(t *testing.T) {
	var store = newConfigurationStore()
	err := store.SetConfigurations(map[string]experimentConfiguration{
		"randomization_algo": testExp,
	})
	assert.NoError(t, err)
	result, err := store.GetConfiguration("randomization_algo")

	expected := "randomization_algo"

	assert.NoError(t, err)
	assert.Equal(t, expected, result.Name)
}
