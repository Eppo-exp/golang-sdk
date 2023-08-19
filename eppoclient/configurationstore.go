package eppoclient

import (
	"encoding/json"
	"errors"
	"log"

	lru "github.com/hashicorp/golang-lru"
)

type configurationStore struct {
	cache *lru.Cache
}

type Variation struct {
	Name       string     `json:"name"`
	Value      Value      `json:"value"`
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
	Overrides     dictionary            `json:"typedOverrides"`
	Allocations   map[string]Allocation `json:"allocations"`
}

func newConfigurationStore(maxEntries int) *configurationStore {
	var configStore = &configurationStore{}

	lruCache, err := lru.New(maxEntries)
	configStore.cache = lruCache

	if err != nil {
		panic(err)
	}

	return configStore
}

func (cs *configurationStore) GetConfiguration(key string) (expConfig experimentConfiguration, err error) {
	value, _ := cs.cache.Get(key)

	if value == nil {
		err = errors.New("not found")
		return
	}

	jsonString, err := json.Marshal(value)

	if err != nil {
		log.Fatalln("Incorrect json")
	}
	ec := experimentConfiguration{}
	json.Unmarshal(jsonString, &ec)

	return ec, nil
}

func (cs *configurationStore) SetConfigurations(configs dictionary) {
	for key, element := range configs {
		cs.cache.Add(key, element)
	}
}
