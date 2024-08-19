package eppoclient

import (
	"math"
	"testing"
)

func TestToFloat64(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		expected  float64
		expectErr bool
	}{
		{"Float64Input", 123.456, 123.456, false},
		{"Float32Input", float32(123.456), 123.456, false},
		{"Int8Input", int8(123), 123.0, false},
		{"Int16Input", int16(123), 123.0, false},
		{"Int32Input", int16(123), 123.0, false},
		{"Int64Input", int16(123), 123.0, false},
		{"UInt8Input", uint8(123), 123.0, false},
		{"UInt16Input", uint16(123), 123.0, false},
		{"UInt32Input", uint16(123), 123.0, false},
		{"UInt64Input", uint16(123), 123.0, false},
		{"StringIntInputValid", "789", 789.0, false},
		{"StringFloatInputValid", "789.012", 789.012, false},
		{"StringNegativeInputValid", "-789.012", -789.012, false},
		{"StringInputInvalid", "abc", 0, true},
		{"SemVerInputInvalid", "1.2.3", 0, true},
		{"BoolInput", true, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := toFloat64(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("ToFloat64(%v) error = %v, expectErr %v", tt.input, err, tt.expectErr)
				return
			}
			closeEnough := math.Abs(result-tt.expected)/tt.expected < 0.00001
			if !tt.expectErr && !closeEnough {
				t.Errorf("ToFloat64(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
