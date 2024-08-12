package eppoclient

import (
	"sync/atomic"
)

type configuration struct {
	flags   configResponse
	bandits banditResponse
	// flag key -> variation value -> banditVariation.
	//
	// This is cached from `flags` field for easier access in
	// evaluation.
	banditFlagAssociations map[string]map[string]banditVariation
}

func (c *configuration) refreshBanditFlagAssociations() {
	associations := make(map[string]map[string]banditVariation)

	for _, banditVariations := range c.flags.Bandits {
		for _, bandit := range banditVariations {
			byVariation, ok := associations[bandit.FlagKey]
			if !ok {
				byVariation = make(map[string]banditVariation)
				associations[bandit.FlagKey] = byVariation
			}
			byVariation[bandit.VariationValue] = bandit
		}
	}

	c.banditFlagAssociations = associations
}

// Return `true` if `flagKey` has associated bandits.
func (c configuration) isBanditFlag(flagKey string) bool {
	_, ok := c.banditFlagAssociations[flagKey]
	return ok
}

func (c configuration) getBanditVariant(flagKey, variation string) (result banditVariation, ok bool) {
	byVariation, ok := c.banditFlagAssociations[flagKey]
	if !ok {
		return result, false
	}
	result, ok = byVariation[variation]
	return result, ok
}

func (c configuration) getFlagConfiguration(key string) (flagConfiguration, error) {
	flag, ok := c.flags.Flags[key]
	if !ok {
		return flag, ErrFlagConfigurationNotFound
	}

	return flag, nil
}

func (c configuration) getBanditConfiguration(key string) (banditConfiguration, error) {
	bandit, ok := c.bandits.Bandits[key]
	if !ok {
		return bandit, ErrBanditConfigurationNotFound
	}

	return bandit, nil
}

// `configurationStore` is a thread-safe in-memory storage. It stores
// the currently active configuration and provides access to multiple
// readers (e.g., flag/bandit evaluation) and writers (e.g.,
// configuration requestor).
type configurationStore struct {
	configuration atomic.Pointer[configuration]
}

func newConfigurationStore(configuration configuration) *configurationStore {
	store := &configurationStore{}
	store.setConfiguration(configuration)
	return store
}

// Returns a snapshot of the currently active configuration.
func (cs *configurationStore) getConfiguration() configuration {
	return *cs.configuration.Load()
}

func (cs *configurationStore) setConfiguration(configuration configuration) {
	configuration.refreshBanditFlagAssociations()
	cs.configuration.Store(&configuration)
}
