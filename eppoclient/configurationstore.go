package eppoclient

import (
	"errors"
)

type configurationStore struct {
	flags   map[string]flagConfiguration
}

func newConfigurationStore() *configurationStore {
	return &configurationStore{}
}

func (cs *configurationStore) getFlagConfiguration(key string) (flag flagConfiguration, err error) {
	flag, ok := cs.flags[key]
	if !ok {
		return flag, errors.New("flag configuration not found in configuration store")
	}

	return flag, nil
}

func (cs *configurationStore) setFlagsConfiguration(configs map[string]flagConfiguration) error {
	cs.flags = configs
	return nil
}
