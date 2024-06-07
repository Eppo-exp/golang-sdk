package eppoclient

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"

	semver "github.com/Masterminds/semver/v3"
)

func (rule rule) matches(subjectAttributes SubjectAttributes) bool {
	for _, condition := range rule.Conditions {
		if !condition.matches(subjectAttributes) {
			return false
		}
	}

	return true
}

func (condition condition) matches(subjectAttributes SubjectAttributes) bool {
	subjectValue, exists := subjectAttributes[condition.Attribute]
	if condition.Operator == "IS_NULL" {
		isNull := !exists || subjectValue == nil
		expectedNull, ok := condition.Value.(bool)
		if !ok {
			return false
		}

		return isNull == expectedNull
	}

	if !exists {
		return false
	}

	switch condition.Operator {
	case "MATCHES":
		v := reflect.ValueOf(subjectValue)
		if v.Kind() != reflect.String {
			subjectValue = strconv.Itoa(subjectValue.(int))
		}
		r, _ := regexp.MatchString(condition.Value.(string), subjectValue.(string))
		return r
	case "ONE_OF":
		return isOneOf(subjectValue, convertToStringArray(condition.Value))
	case "NOT_ONE_OF":
		return !isOneOf(subjectValue, convertToStringArray(condition.Value))
	case "GTE", "GT", "LTE", "LT":
		// Attempt to coerce both values to float64 and compare them.
		subjectValueNumeric, isNumericSubjectErr := toFloat64(subjectValue)
		conditionValueNumeric, isNumericConditionErr := toFloat64(condition.Value)
		if isNumericSubjectErr == nil && isNumericConditionErr == nil {
			return evaluateNumericCondition(subjectValueNumeric, conditionValueNumeric, condition)
		}

		// Attempt to compare using semantic versioning if both values are strings.
		subjectValueStr, isStringSubject := subjectValue.(string)
		conditionValueStr, isStringCondition := condition.Value.(string)
		if isStringSubject && isStringCondition {
			// Attempt to parse both values as semantic versions.
			subjectSemVer, errSubject := semver.NewVersion(subjectValueStr)
			conditionSemVer, errCondition := semver.NewVersion(conditionValueStr)

			// If parsing succeeds, evaluate the semver condition.
			if errSubject == nil && errCondition == nil {
				return evaluateSemVerCondition(subjectSemVer, conditionSemVer, condition)
			}
		}

		// Fallback logic if neither numeric nor semver comparison is applicable.
		return false
	default:
		fmt.Printf("unknown condition operator: %s", condition.Operator)
		return false
	}
}

func convertToStringArray(conditionValue interface{}) []string {
	if reflect.TypeOf(conditionValue).Elem().Kind() == reflect.String {
		return conditionValue.([]string)
	}
	conditionValueStrings := make([]string, len(conditionValue.([]interface{})))
	for i, v := range conditionValue.([]interface{}) {
		conditionValueStrings[i] = v.(string)
	}
	return conditionValueStrings
}

func isOneOf(attributeValue interface{}, conditionValue []string) bool {
	for _, value := range conditionValue {
		if isOne(attributeValue, value) {
			return true
		}
	}

	return false
}

// Return true if `attributeValue` is the same as `s` under eppo
// evaluation rules.
func isOne(attributeValue interface{}, s string) bool {
	switch attributeValue.(type) {
	case string:
		return attributeValue == s
	case float32:
		value, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return false
		}
		return attributeValue == value
	case float64:
		value, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return false
		}
		return attributeValue == value
	case int, int8, int16, int32, int64:
		value, err := strconv.ParseInt(s, 0, 64)
		if err != nil {
			return false
		}
		return promoteInt(attributeValue) == value
	case uint, uint8, uint16, uint32, uint64:
		value, err := strconv.ParseUint(s, 0, 64)
		if err != nil {
			return false
		}
		return promoteUint(attributeValue) == value
	case bool:
		value, err := strconv.ParseBool(s)
		if err != nil {
			return false
		}
		return attributeValue == value
	default:
		attributeValue = fmt.Sprintf("%v", attributeValue)
		return attributeValue == s
	}
}

func promoteInt(i interface{}) int64 {
	switch i := i.(type) {
	case int:
		return int64(i)
	case int8:
		return int64(i)
	case int16:
		return int64(i)
	case int32:
		return int64(i)
	case int64:
		return i
	}
	panic(fmt.Errorf("unexpected type passed to promoteInt: %T", i))
}

func promoteUint(i interface{}) uint64 {
	switch i := i.(type) {
	case uint:
		return uint64(i)
	case uint8:
		return uint64(i)
	case uint16:
		return uint64(i)
	case uint32:
		return uint64(i)
	case uint64:
		return i
	}
	panic(fmt.Errorf("unexpected type passed to promoteUint: %T", i))
}

func evaluateSemVerCondition(subjectValue *semver.Version, conditionValue *semver.Version, condition condition) bool {
	comp := subjectValue.Compare(conditionValue)
	switch condition.Operator {
	case "GT":
		return comp > 0
	case "GTE":
		return comp >= 0
	case "LT":
		return comp < 0
	case "LTE":
		return comp <= 0
	default:
		panic("Incorrect condition operator")
	}
}

func evaluateNumericCondition(subjectValue float64, conditionValue float64, condition condition) bool {
	switch condition.Operator {
	case "GT":
		return subjectValue > conditionValue
	case "GTE":
		return subjectValue >= conditionValue
	case "LT":
		return subjectValue < conditionValue
	case "LTE":
		return subjectValue <= conditionValue
	default:
		panic("Incorrect condition operator")
	}
}

// toFloat64 attempts to convert an interface{} value to a float64.
// It supports inputs of type float64 or string (which can be parsed as float64).
// Returns a float64 and nil error on success, or 0 and an error on failure.
func toFloat64(val interface{}) (float64, error) {
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
