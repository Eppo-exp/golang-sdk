package eppoclient

import (
	"errors"

	lru "github.com/hashicorp/golang-lru"
)

type ConfigurationStore struct {
	cache lru.Cache
}

type Variation struct {
	Name       string `json:"name"`
	ShardRange ShardRange
}

type ExperimentConfiguration struct {
	Name            string      `json:"name"`
	PercentExposure float32     `json:"percentExposure"`
	Enabled         bool        `json:"enabled"`
	SubjectShards   int         `json:"subjectShards"`
	Variations      []Variation `json:"variations"`
	Rules           []Rule      `json:"rules"`
	Overrides       Dictionary  `json:"overrides"`
}

func NewConfigurationStore(maxEntries int) *ConfigurationStore {
	var configStore = &ConfigurationStore{}

	lruCache, err := lru.New(maxEntries)
	configStore.cache = *lruCache

	if err != nil {
		panic(err)
	}

	return configStore
}

func (cs *ConfigurationStore) GetConfiguration(key string) (expConfig ExperimentConfiguration, err error) {
	value, _ := cs.cache.Get(key)

	if value == nil {
		err = errors.New("not found")
		return
	}

	return value.(ExperimentConfiguration), nil
}

func (cs *ConfigurationStore) SetConfigurations(configs Dictionary) {
	for key, element := range configs {
		cs.cache.Add(key, element)
	}
}
