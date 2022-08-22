package eppoclient

import (
	"fmt"
	"testing"
)

var testExp = ExperimentConfiguration{
	SubjectShards:   1000,
	PercentExposure: 1,
	Enabled:         true,
	Variations:      []Variation{},
	Name:            "randomization_algo",
}

const TEST_MAX_SIZE = 10

var store = ConfigurationStore{}

func init() {
	store.New(TEST_MAX_SIZE)
}

func Test_GetConfiguration_unknownKey(t *testing.T) {
	store.SetConfigurations(Dictionary{"randomization_algo": testExp})
	result, err := store.GetConfiguration("unknown_exp")

	if err == nil {
		t.Errorf("\"store.GetConfiguration()\" FAILED, expected -> %v, got -> %v", "", result)
	}
}

func Test_GetConfiguration_knownKey(t *testing.T) {
	store.SetConfigurations(Dictionary{"randomization_algo": testExp})
	result, err := store.GetConfiguration("randomization_algo")

	if err != nil {
		t.Errorf("\"store.GetConfiguration()\" FAILED, expected -> %v, got -> %v", "", result)
	}

	expected := "randomization_algo"

	if result.Name != expected {
		t.Errorf("\"store.GetConfiguration()\" FAILED, expected -> %v, got -> %v", expected, result.Name)
	}
}

func Test_GetConfiguration_evictsOldEntriesWhenMaxSizeExceeded(t *testing.T) {
	store.SetConfigurations(Dictionary{"item_to_be_evicted": testExp})
	result, err := store.GetConfiguration("item_to_be_evicted")

	expected := "randomization_algo"

	if err != nil || result.Name != expected {
		t.Errorf("\"store.GetConfiguration()\" FAILED, expected -> %v, got -> %v", expected, result.Name)
	}

	for i := 0; i < TEST_MAX_SIZE; i++ {
		dictKey := fmt.Sprintf("test-entry-%v", i)
		store.SetConfigurations(Dictionary{dictKey: testExp})
	}

	result, err = store.GetConfiguration("item_to_be_evicted")

	if err == nil {
		t.Errorf("\"store.GetConfiguration()\" FAILED, expected -> %v, got -> %v", "", result)
	}

	result, err = store.GetConfiguration(fmt.Sprintf("test-entry-%v", TEST_MAX_SIZE-1))

	if err != nil || result.Name != expected {
		t.Errorf("\"store.GetConfiguration()\" FAILED, expected -> %v, got -> %v", expected, result.Name)
	}
}
