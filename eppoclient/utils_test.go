package eppoclient

import (
	"testing"
	"time"
)

func TestToFloat64(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		expected  float64
		expectErr bool
	}{
		{"Float64Input", 123.456, 123.456, false},
		{"StringInputValid", "789.012", 789.012, false},
		{"StringInputInvalid", "abc", 0, true},
		{"SemVerInputInvalid", "1.2.3", 0, true},
		{"BoolInput", true, 0, true},
		{"IntInput", 123, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ToFloat64(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("ToFloat64(%v) error = %v, expectErr %v", tt.input, err, tt.expectErr)
				return
			}
			if !tt.expectErr && result != tt.expected {
				t.Errorf("ToFloat64(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTimeNow(t *testing.T) {
	result := TimeNow()
	_, err := time.Parse(time.RFC3339, result)
	if err != nil {
		t.Errorf("TimeNow() = %v, want %v", result, "")
	}
}
