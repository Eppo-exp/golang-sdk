package eppoclient

import (
	"errors"
)

type configurationStore struct {
	flags   map[string]flagConfiguration
	bandits map[string]banditConfiguration
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

func (cs *configurationStore) setFlagsConfiguration(configs map[string]flagConfiguration) {
	cs.flags = configs
}

func (cs *configurationStore) getBanditConfiguration(key string) (bandit banditConfiguration, err error) {
	bandit, ok := cs.bandits[key]
	if !ok {
		return bandit, errors.New("bandit configuration not found in configuration store")
	}

	return bandit, nil
}

func (cs *configurationStore) setBanditsConfiguration(configs map[string]banditConfiguration) {
	cs.bandits = configs
}
