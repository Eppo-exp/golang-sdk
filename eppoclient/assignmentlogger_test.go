package eppoclient

import (
	"encoding/json"
	"reflect"
	"testing"
)

// TestAssignmentEventSerialization tests the JSON serialization and deserialization of AssignmentEvent
func TestAssignmentEventSerialization(t *testing.T) {
	// Create a test case with each type
	testCases := []AssignmentEvent{
		{
			Experiment:  "testExperiment",
			FeatureFlag: "testFeatureFlag",
			Allocation:  "testAllocation",
			Variation:   Bool(true),
			Subject:     "testSubject",
			Timestamp:   "testTimestamp",
			// SubjectAttributes: dictionary{"testKey": String("testValue")},
		},
		{
			Experiment:  "testExperiment",
			FeatureFlag: "testFeatureFlag",
			Allocation:  "testAllocation",
			Variation:   Numeric(123.45),
			Subject:     "testSubject",
			Timestamp:   "testTimestamp",
			//SubjectAttributes: dictionary{"testKey": String("testValue")},
		},
		{
			Experiment:  "testExperiment",
			FeatureFlag: "testFeatureFlag",
			Allocation:  "testAllocation",
			Variation:   String("testVariation"),
			Subject:     "testSubject",
			Timestamp:   "testTimestamp",
			//SubjectAttributes: dictionary{"testKey": String("testValue")},
		},
		{
			Experiment:  "testExperiment",
			FeatureFlag: "testFeatureFlag",
			Allocation:  "testAllocation",
			Variation:   String("{\"foo\":\"bar\",\"car\":\"far\"}"),
			Subject:     "testSubject",
			Timestamp:   "testTimestamp",
			//SubjectAttributes: dictionary{"testKey": String("testValue")},
		},
	}

	for _, original := range testCases {
		// Marshal to JSON
		marshaled, err := json.Marshal(original)
		if err != nil {
			t.Errorf("Failed to marshal: %v", err)
		}

		// Unmarshal from JSON
		var unmarshaled AssignmentEvent
		err = json.Unmarshal(marshaled, &unmarshaled)
		if err != nil {
			t.Errorf("Failed to unmarshal: %v", err)
		}

		// Compare the original and unmarshaled
		if !reflect.DeepEqual(original, unmarshaled) {
			t.Errorf("Original and unmarshaled Value are not equal. Original: %+v, Unmarshaled: %+v", original, unmarshaled)
		}
	}
}
