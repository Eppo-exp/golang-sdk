package eppoclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_maskSensitiveInfo(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "mask apiKey",
			input:    "https://example.com?apiKey=123456&anotherParam=foo",
			expected: "https://example.com?apiKey=XXXXXX&anotherParam=foo",
		},
		{
			name:     "mask sdkKey",
			input:    "https://example.com?sdkKey=abcdef&anotherParam=foo",
			expected: "https://example.com?sdkKey=XXXXXX&anotherParam=foo",
		},
		{
			name:     "no sensitive info",
			input:    "https://example.com?param=value&anotherParam=foo",
			expected: "https://example.com?param=value&anotherParam=foo",
		},
		{
			name:     "mask apiKey and sdkKey",
			input:    "https://example.com?apiKey=123456&sdkKey=abcdef&anotherParam=foo",
			expected: "https://example.com?apiKey=XXXXXX&sdkKey=XXXXXX&anotherParam=foo",
		},
		{
			name:     "mask apiKey and sdkKey out of order",
			input:    "https://example.com?anotherParam=foo&apiKey=123456&sdkKey=abcdef",
			expected: "https://example.com?anotherParam=foo&apiKey=XXXXXX&sdkKey=XXXXXX",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maskSensitiveInfo(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
