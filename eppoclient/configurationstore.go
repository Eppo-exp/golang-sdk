package eppoclient

import (
	"errors"
)

type configurationStore struct {
	configs map[string]flagConfiguration
}

func newConfigurationStore() *configurationStore {
	return &configurationStore{
		configs: make(map[string]flagConfiguration),
	}
}

func (cs *configurationStore) GetConfiguration(key string) (flag flagConfiguration, err error) {
	flag, ok := cs.configs[key]
	if !ok {
		return flag, errors.New("configuration not found in configuration store")
	}

	return flag, nil
}

func (cs *configurationStore) SetConfigurations(configs map[string]flagConfiguration) error {
	cs.configs = configs
	return nil
}
