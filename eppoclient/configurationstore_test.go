package eppoclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetConfiguration_unknownKey(t *testing.T) {
	store := newConfigurationStore()

	config := store.getConfiguration()
	result, err := config.getFlagConfiguration("unknown_exp")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func Test_GetConfiguration_knownKey(t *testing.T) {
	flags := configResponse{
		Flags: map[string]*flagConfiguration{
			"experiment-key-1": &flagConfiguration{
				Key:           "experiment-key-1",
				Enabled:       false,
				VariationType: stringVariation,
			},
		},
	}
	store := newConfigurationStoreWithConfig(configuration{flags: flags})

	config := store.getConfiguration()
	result, err := config.getFlagConfiguration("experiment-key-1")

	expected := "experiment-key-1"

	assert.NoError(t, err)
	assert.Equal(t, expected, result.Key)
}
