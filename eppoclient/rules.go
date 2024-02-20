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
	subjectValue := subjectAttributes[condition.Attribute]

	if subjectValue != nil {
		if condition.Operator == "MATCHES" {
			v := reflect.ValueOf(subjectValue)
			if v.Kind() != reflect.String {
				subjectValue = strconv.Itoa(subjectValue.(int))
			}
			r, _ := regexp.MatchString(condition.Value.(string), subjectValue.(string))
			return r
		} else if condition.Operator == "ONE_OF" {
			return isOneOf(subjectValue, convertToStringArray(condition.Value))
		} else if condition.Operator == "NOT_ONE_OF" {
			return isNotOneOf(subjectValue, convertToStringArray(condition.Value))
		} else {
			// If the condition value is a string, we try to convert it to a semver.
			// If it's not a semver, we fall back to numeric comparison.
			subjectValueStr, subjectValueStrOk := subjectValue.(string)
			conditionValueStr, conditionValueStrOk := condition.Value.(string)

			if subjectValueStrOk && conditionValueStrOk {
				subjectValueSemVer, subjectValueSemVerError := semver.NewVersion(subjectValueStr)
				conditionValueSemVer, conditionValueSemVerError := semver.NewVersion(conditionValueStr)
				if subjectValueSemVerError == nil && conditionValueSemVerError == nil {
					return evaluateSemVerCondition(subjectValueSemVer, conditionValueSemVer, condition)
				}
			}

			return evaluateNumericCondition(subjectValue, condition)
		}
	}
	return false
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

func evaluateNumericCondition(subjectValue interface{}, condition condition) bool {
	v := reflect.ValueOf(subjectValue)

	if v.Kind() == reflect.String {
		return false
	}

	if v.Kind() == reflect.Int {
		subjectValue = float64(subjectValue.(int))
	}

	switch condition.Operator {
	case "GT":
		return subjectValue.(float64) > condition.Value.(float64)
	case "GTE":
		return subjectValue.(float64) >= condition.Value.(float64)
	case "LT":
		return subjectValue.(float64) < condition.Value.(float64)
	case "LTE":
		return subjectValue.(float64) <= condition.Value.(float64)
	default:
		panic("Incorrect condition operator")
	}
}
