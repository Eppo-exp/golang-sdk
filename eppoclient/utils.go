package eppoclient

import (
	"errors"
	"fmt"
	"strconv"
)

type dictionary map[string]interface{}

type testData struct {
	Experiment             string                  `json:"experiment"`
	ValueType              string                  `json:"valueType"`
	PercentExposure        float32                 `json:"percentExposure"`
	Variations             []testDataVariations    `json:"variations"`
	Subjects               []string                `json:"subjects"`
	SubjectsWithAttributes []subjectWithAttributes `json:"subjectsWithAttributes"`
	ExpectedAssignments    []Value                 `json:"expectedAssignments"`
}

type subjectWithAttributes struct {
	SubjectKey        string     `json:"subjectKey"`
	SubjectAttributes dictionary `json:"subjectAttributes"`
}

type testDataVariations struct {
	Name       string             `json:"name"`
	ShardRange testDataShardRange `json:"shardRange"`
}

type testDataShardRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// ToFloat64 attempts to convert an interface{} value to a float64.
// It supports inputs of type float64 or string (which can be parsed as float64).
// Returns a float64 and nil error on success, or 0 and an error on failure.
func ToFloat64(val interface{}) (float64, error) {
	switch v := val.(type) {
	case float64:
		return v, nil
	case string:
		floatVal, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, fmt.Errorf("cannot convert string '%s' to float64: %w", v, err)
		}
		return floatVal, nil
	default:
		return 0, errors.New("value is neither a float64 nor a convertible string")
	}
}
