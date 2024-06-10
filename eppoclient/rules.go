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

type condition struct {
	Attribute string      `json:"attribute"`
	Value     interface{} `json:"value"`
	Operator  string      `validator:"regexp=^(MATCHES|GTE|GT|LTE|LT|ONE_OF|NOT_ONE_OF)$" json:"operator"`
}

type rule struct {
	AllocationKey string      `json:"allocationKey"`
	Conditions    []condition `json:"conditions"`
}

func findMatchingRule(subjectAttributes dictionary, rules []rule) (rule, error) {
	for _, rule := range rules {
		if matchesRule(subjectAttributes, rule) {
			return rule, nil
		}
	}

	return rule{}, errors.New("no matching rule")
}

func matchesRule(subjectAttributes dictionary, rule rule) bool {
	for _, condition := range rule.Conditions {
		if !evaluateCondition(subjectAttributes, condition) {
			return false
		}
	}

	return true
}

func evaluateCondition(subjectAttributes dictionary, condition condition) bool {
	subjectValue, exists := subjectAttributes[condition.Attribute]
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
		return isNotOneOf(subjectValue, convertToStringArray(condition.Value))
	case "GTE", "GT", "LTE", "LT":
		// Attempt to coerce both values to float64 and compare them.
		subjectValueNumeric, isNumericSubjectErr := ToFloat64(subjectValue)
		conditionValueNumeric, isNumericConditionErr := ToFloat64(condition.Value)
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
	matches := getMatchingStringValues(attributeValue, conditionValue)
	return len(matches) > 0
}

func isNotOneOf(attributeValue interface{}, conditionValue []string) bool {
	matches := getMatchingStringValues(attributeValue, conditionValue)
	return len(matches) == 0
}

func getMatchingStringValues(attributeValue interface{}, conditionValue []string) []string {
	v := reflect.ValueOf(attributeValue)

	if v.Kind() != reflect.String {
		attributeValue = fmt.Sprintf("%v", attributeValue)
	}

	var result []string

	for _, value := range conditionValue {
		if strings.EqualFold(value, attributeValue.(string)) {
			result = append(result, value)
		}
	}

	return result
}

func evaluateSemVerCondition(subjectValue *semver.Version, conditionValue *semver.Version, condition condition) bool {
	switch condition.Operator {
	case "GT":
		return subjectValue.GreaterThan(conditionValue)
	case "GTE":
		return subjectValue.GreaterThan(conditionValue) || subjectValue.Equal(conditionValue)
	case "LT":
		return subjectValue.LessThan(conditionValue)
	case "LTE":
		return subjectValue.LessThan(conditionValue) || subjectValue.Equal(conditionValue)
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
