package eppoclient

import (
	"sync/atomic"
)

// `configurationStore` is a thread-safe in-memory storage. It stores
// the currently active configuration and provides access to multiple
// readers (e.g., flag/bandit evaluation) and writers (e.g.,
// configuration requestor).
type configurationStore struct {
	configuration atomic.Pointer[configuration]

	// `initializedCh` is closed when we receive a proper
	// configuration.
	initializedCh chan struct{}
	// `isInitialized` is used to protect `initializedCh`, so we
	// don’t double-close it (which is an error in Go).
	isInitialized atomic.Bool
}

func newConfigurationStore() *configurationStore {
	return &configurationStore{
		initializedCh: make(chan struct{}),
	}
}

func newConfigurationStoreWithConfig(configuration configuration) *configurationStore {
	store := newConfigurationStore()
	store.setConfiguration(configuration)
	return store
}

// Returns a snapshot of the currently active configuration.
func (cs *configurationStore) getConfiguration() configuration {
	if config := cs.configuration.Load(); config != nil {
		return *config
	} else {
		return configuration{}
	}
}

func (cs *configurationStore) setConfiguration(configuration configuration) {
	configuration.precompute()
	cs.configuration.Store(&configuration)
	cs.setInitialized()
}

// Set `initialized` flag to `true` notifying anyone waiting on it.
func (cs *configurationStore) setInitialized() {
	if cs.isInitialized.CompareAndSwap(false, true) {
		// Channels can only be closed once, so we protect the
		// call to `close` with a CAS.
		close(cs.initializedCh)
	}
}

// Returns a channel that gets closed after configuration store is
// successfully initialized.
func (cs *configurationStore) Initialized() <-chan struct{} {
	return cs.initializedCh
}
