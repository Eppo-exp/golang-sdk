package eppoclient

import (
	"testing"

	semver "github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
)

func TestConditionPrecompute(t *testing.T) {
	tests := []struct {
		name                   string
		condition              condition
		expectedNumVal         float64
		expectedNumValValid    bool
		expectedSemVerVal      *semver.Version
		expectedSemVerValValid bool
	}{
		{
			name: "valid numeric value",
			condition: condition{
				Value: 42.0,
			},
			expectedNumVal:         42.0,
			expectedNumValValid:    true,
			expectedSemVerValValid: false,
		},
		{
			name: "valid semver value",
			condition: condition{
				Value: "1.2.3",
			},
			expectedNumValValid:    false,
			expectedSemVerVal:      semver.MustParse("1.2.3"),
			expectedSemVerValValid: true,
		},
		{
			name: "invalid value",
			condition: condition{
				Value: "not a number or semver",
			},
			expectedNumValValid:    false,
			expectedSemVerValValid: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.condition.Precompute()
			assert.Equal(t, tc.expectedNumVal, tc.condition.NumericValue)
			assert.Equal(t, tc.expectedNumValValid, tc.condition.NumericValueValid)
			assert.Equal(t, tc.expectedSemVerVal, tc.condition.SemVerValue)
			assert.Equal(t, tc.expectedSemVerValValid, tc.condition.SemVerValueValid)
		})
	}
}
