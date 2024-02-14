package eppoclient

import (
	"errors"
	lru "github.com/hashicorp/golang-lru/v2"
)

type configurationStore struct {
	cache *lru.Cache[string, experimentConfiguration]
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

func newConfigurationStore(maxEntries int) *configurationStore {
	var configStore = &configurationStore{}

	lruCache, err := lru.New[string, experimentConfiguration](maxEntries)
	configStore.cache = lruCache

	if err != nil {
		panic(err)
	}

	return configStore
}

func (cs *configurationStore) GetConfiguration(key string) (expConfig experimentConfiguration, err error) {
	// Attempt to get the value from the cache
	expConfig, ok := cs.cache.Get(key)
	if !ok {
		return expConfig, errors.New("configuration not found in cache")
	}

	return expConfig, nil
}

func (cs *configurationStore) SetConfigurations(configs map[string]experimentConfiguration) error {
	for key, config := range configs {
		cs.cache.Add(key, config)
	}
	return nil
}
