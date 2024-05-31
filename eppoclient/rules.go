package eppoclient

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

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
		panic(fmt.Sprintf("unknown condition operator: %s", condition.Operator))
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
	v := reflect.ValueOf(attributeValue)

	if v.Kind() != reflect.String {
		attributeValue = fmt.Sprintf("%v", attributeValue)
	}

	for _, value := range conditionValue {
		if strings.EqualFold(value, attributeValue.(string)) {
			return true
		}
	}

	return false
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
