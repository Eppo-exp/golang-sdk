package eppoclient

import (
	"testing"
)

func Test_getShard(t *testing.T) {
	input := "test-string"
	var expected int64 = 0
	result := getShard(input, 5)

	if result != expected {
		t.Errorf("\"getShard('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
	} else {
		t.Logf("\"getShard('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
	}

	expected = 1
	result = getShard(input, 2)

	if result != expected {
		t.Errorf("\"getShard('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
	} else {
		t.Logf("\"getShard('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
	}
}

func Test_isShardInRange_Fail(t *testing.T) {
	input := 5
	inputRange := ShardRange{Start: 1, End: 5}
	expected := false
	result := isShardInRange(input, inputRange)

	if result != expected {
		t.Errorf("\"isShardInRange(%v,  %+v)\" FAILED, expected -> %t, got -> %t", input, inputRange, expected, result)
	} else {
		t.Logf("\"isShardInRange(%v,  %+v)\" SUCCEDED, expected -> %t, got -> %t", input, inputRange, expected, result)
	}
}

func Test_isShardInRange_Success(t *testing.T) {
	input := 3
	inputRange := ShardRange{Start: 1, End: 7}
	expected := true
	result := isShardInRange(input, inputRange)

	if result != expected {
		t.Errorf("\"isShardInRange(%v,  %+v)\" FAILED, expected -> %t, got -> %t", input, inputRange, expected, result)
	} else {
		t.Logf("\"isShardInRange(%v,  %+v)\" SUCCEDED, expected -> %t, got -> %t", input, inputRange, expected, result)
	}
}
