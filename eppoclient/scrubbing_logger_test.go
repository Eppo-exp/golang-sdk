package eppoclient

import (
	"errors"
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

func TestScrubbingLogger_scrub(t *testing.T) {
	logger := NewScrubbingLogger(nil)

	t.Run("scrubs string arguments", func(t *testing.T) {
		args := []interface{}{"https://example.com?apiKey=secret123"}
		scrubbed := logger.scrub(args...)
		assert.Equal(t, "https://example.com?apiKey=XXXXXX", scrubbed[0])
	})

	t.Run("scrubs error arguments", func(t *testing.T) {
		err := errors.New("Get \"https://example.com?apiKey=secret123&sdkName=go\": connection refused")
		args := []interface{}{err}
		scrubbed := logger.scrub(args...)
		scrubbedErr, ok := scrubbed[0].(error)
		assert.True(t, ok, "expected scrubbed value to be an error")
		assert.Equal(t, "Get \"https://example.com?apiKey=XXXXXX&sdkName=go\": connection refused", scrubbedErr.Error())
	})

	t.Run("scrubs mixed arguments", func(t *testing.T) {
		err := errors.New("failed: apiKey=secret456")
		args := []interface{}{"message with sdkKey=abc123", err, 42}
		scrubbed := logger.scrub(args...)

		assert.Equal(t, "message with sdkKey=XXXXXX", scrubbed[0])

		scrubbedErr, ok := scrubbed[1].(error)
		assert.True(t, ok)
		assert.Equal(t, "failed: apiKey=XXXXXX", scrubbedErr.Error())

		assert.Equal(t, 42, scrubbed[2]) // non-string/error passes through
	})

	t.Run("passes through non-sensitive data unchanged", func(t *testing.T) {
		args := []interface{}{"normal message", errors.New("normal error"), 123, true}
		scrubbed := logger.scrub(args...)

		assert.Equal(t, "normal message", scrubbed[0])
		assert.Equal(t, "normal error", scrubbed[1].(error).Error())
		assert.Equal(t, 123, scrubbed[2])
		assert.Equal(t, true, scrubbed[3])
	})
}
