package eppoclient

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"

	semver "github.com/Masterminds/semver/v3"
)

func (rule rule) matches(subjectAttributes Attributes) bool {
	for _, condition := range rule.Conditions {
		if !condition.matches(subjectAttributes) {
			return false
		}
	}

	return true
}

func (condition condition) matches(subjectAttributes Attributes) bool {
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
		return matches(subjectValue, condition.Value.(string))
	case "NOT_MATCHES":
		return !matches(subjectValue, condition.Value.(string))
	case "ONE_OF":
		return isOneOf(subjectValue, convertToStringArray(condition.Value))
	case "NOT_ONE_OF":
		return !isOneOf(subjectValue, convertToStringArray(condition.Value))
	case "GTE", "GT", "LTE", "LT":
		// Attempt to coerce the subject value to float64 and compare it
		// against the condition value.
		subjectValueNumeric, isNumericSubjectErr := toFloat64(subjectValue)
		if isNumericSubjectErr == nil && condition.NumericValueValid {
			result, err := evaluateNumericCondition(subjectValueNumeric, condition.NumericValue, condition)
			if err != nil {
				return false
			}
			return result
		}

		// Attempt to compare using semantic versioning if the subject value is a string.
		// and the condition value is a valid semantic version.
		subjectValueStr, isStringSubject := subjectValue.(string)
		if isStringSubject && condition.SemVerValueValid {
			// Attempt to parse the subject value as a semantic version.
			subjectSemVer, errSubject := semver.NewVersion(subjectValueStr)

			// If parsing succeeds, evaluate the semver condition.
			if errSubject == nil {
				result, err := evaluateSemVerCondition(subjectSemVer, condition.SemVerValue, condition)
				if err != nil {
					return false
				}
				return result
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

func matches(subjectValue interface{}, conditionValue string) bool {
	var v string
	switch subjectValue := subjectValue.(type) {
	case string:
		v = subjectValue
	case int:
		v = strconv.Itoa(subjectValue)
	case bool:
		if subjectValue {
			v = "true"
		} else {
			v = "false"
		}
	default:
		return false
	}

	r, _ := regexp.MatchString(conditionValue, v)
	return r
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
		promotedAttributeValue, err := promoteInt(attributeValue)
		if err != nil {
			return false
		}
		return promotedAttributeValue == value
	case uint, uint8, uint16, uint32, uint64:
		value, err := strconv.ParseUint(s, 0, 64)
		if err != nil {
			return false
		}
		promotedAttributeValue, err := promoteUint(attributeValue)
		if err != nil {
			return false
		}
		return promotedAttributeValue == value
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

func evaluateSemVerCondition(subjectValue *semver.Version, conditionValue *semver.Version, condition condition) (bool, error) {
	comp := subjectValue.Compare(conditionValue)
	switch condition.Operator {
	case "GT":
		return comp > 0, nil
	case "GTE":
		return comp >= 0, nil
	case "LT":
		return comp < 0, nil
	case "LTE":
		return comp <= 0, nil
	default:
		return false, fmt.Errorf("incorrect condition operator: %s", condition.Operator)
	}
}

func evaluateNumericCondition(subjectValue float64, conditionValue float64, condition condition) (bool, error) {
	switch condition.Operator {
	case "GT":
		return subjectValue > conditionValue, nil
	case "GTE":
		return subjectValue >= conditionValue, nil
	case "LT":
		return subjectValue < conditionValue, nil
	case "LTE":
		return subjectValue <= conditionValue, nil
	default:
		return false, fmt.Errorf("incorrect condition operator: %s", condition.Operator)
	}
}

// toFloat64 attempts to convert an interface{} value to a float64.
// It supports inputs of type float64 or string (which can be parsed as float64).
// Returns a float64 and nil error on success, or 0 and an error on failure.
func toFloat64(val interface{}) (float64, error) {
	switch v := val.(type) {
	case float32, float64:
		return promoteFloat(v)
	case int, int8, int16, int32, int64:
		promotedInt, err := promoteInt(v)
		if err != nil {
			return 0, err
		}
		return float64(promotedInt), nil
	case uint, uint8, uint16, uint32, uint64:
		promotedUint, err := promoteUint(v)
		if err != nil {
			return 0, err
		}
		return float64(promotedUint), nil
	case string:
		floatVal, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, fmt.Errorf("cannot convert string '%s' to float64: %w", v, err)
		}
		return floatVal, nil
	default:
		return 0, errors.New("value is neither a number nor a convertible string")
	}
}

func promoteInt(i interface{}) (int64, error) {
	switch i := i.(type) {
	case int:
		return int64(i), nil
	case int8:
		return int64(i), nil
	case int16:
		return int64(i), nil
	case int32:
		return int64(i), nil
	case int64:
		return i, nil
	}
	return 0, fmt.Errorf("unexpected type passed to promoteInt: %T", i)
}

func promoteUint(i interface{}) (uint64, error) {
	switch i := i.(type) {
	case uint:
		return uint64(i), nil
	case uint8:
		return uint64(i), nil
	case uint16:
		return uint64(i), nil
	case uint32:
		return uint64(i), nil
	case uint64:
		return i, nil
	}
	return 0, fmt.Errorf("unexpected type passed to promoteUint: %T", i)
}

func promoteFloat(f interface{}) (float64, error) {
	switch f := f.(type) {
	case float32:
		return float64(f), nil
	case float64:
		return f, nil
	}
	return 0, fmt.Errorf("unexpected type passed to promoteFloat: %T", f)
}
