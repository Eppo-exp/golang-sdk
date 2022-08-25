package eppoclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getShard(t *testing.T) {
	var tests = []struct {
		a    string
		b    int64
		want int64
	}{
		{"test-string", 5, 0},
		{"test-string", 2, 1},
	}

	for _, tt := range tests {
		result := getShard(tt.a, tt.b)

		assert.Equal(t, tt.want, result)
	}
}

func Test_isShardInRange_Fail(t *testing.T) {
	input := 5
	inputRange := ShardRange{Start: 1, End: 5}
	expected := false
	result := isShardInRange(input, inputRange)

	assert.Equal(t, expected, result)
}

func Test_isShardInRange_Success(t *testing.T) {
	input := 3
	inputRange := ShardRange{Start: 1, End: 7}
	expected := true
	result := isShardInRange(input, inputRange)

	assert.Equal(t, expected, result)
}
