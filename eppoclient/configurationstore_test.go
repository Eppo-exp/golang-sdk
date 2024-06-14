package eppoclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetConfiguration_unknownKey(t *testing.T) {
	var store = newConfigurationStore()
	store.setFlagsConfiguration(map[string]flagConfiguration{})

	result, err := store.getFlagConfiguration("unknown_exp")

	assert.Error(t, err)
	assert.Equal(t, flagConfiguration{}, result)
}

func Test_GetConfiguration_knownKey(t *testing.T) {
	config := map[string]flagConfiguration{
		"experiment-key-1": flagConfiguration{
			Key:           "experiment-key-1",
			Enabled:       false,
			VariationType: stringVariation,
		},
	}

	var store = newConfigurationStore()
	store.setFlagsConfiguration(config)

	result, err := store.getFlagConfiguration("experiment-key-1")

	expected := "experiment-key-1"

	assert.NoError(t, err)
	assert.Equal(t, expected, result.Key)
}
