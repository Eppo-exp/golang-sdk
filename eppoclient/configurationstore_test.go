package eppoclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetConfiguration_unknownKey(t *testing.T) {
	var store = newConfigurationStore()
	err := store.SetConfigurations(map[string]flagConfiguration{})

	assert.NoError(t, err)
	result, err := store.GetConfiguration("unknown_exp")

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
	err := store.SetConfigurations(config)
	assert.NoError(t, err)
	result, err := store.GetConfiguration("experiment-key-1")

	expected := "experiment-key-1"

	assert.NoError(t, err)
	assert.Equal(t, expected, result.Key)
}
