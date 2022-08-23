package eppoclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getShard(t *testing.T) {
	input := "test-string"
	var expected int64 = 0
	result := getShard(input, 5)

	assert.Equal(t, expected, result)

	expected = 1
	result = getShard(input, 2)

	assert.Equal(t, expected, result)
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
