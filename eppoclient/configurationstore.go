package eppoclient

import (
	"errors"
)

type configurationStore struct {
	configs map[string]experimentConfiguration
}

type Variation struct {
	Name       string     `json:"name"`
	Value      Value      `json:"typedValue"`
	ShardRange shardRange `json:"shardRange"`
}

type Allocation struct {
	PercentExposure float32     `json:"percentExposure"`
	Variations      []Variation `json:"variations"`
}

type experimentConfiguration struct {
	Name          string                `json:"name"`
	Enabled       bool                  `json:"enabled"`
	SubjectShards int                   `json:"subjectShards"`
	Rules         []rule                `json:"rules"`
	Overrides     map[string]Value      `json:"typedOverrides"`
	Allocations   map[string]Allocation `json:"allocations"`
}

func newConfigurationStore() *configurationStore {
	return &configurationStore{
		configs: make(map[string]experimentConfiguration),
	}
}

func (cs *configurationStore) GetConfiguration(key string) (expConfig experimentConfiguration, err error) {
	expConfig, ok := cs.configs[key]
	if !ok {
		return expConfig, errors.New("configuration not found in configuration store")
	}

	return expConfig, nil
}

func (cs *configurationStore) SetConfigurations(configs map[string]experimentConfiguration) error {
	cs.configs = configs
	return nil
}
