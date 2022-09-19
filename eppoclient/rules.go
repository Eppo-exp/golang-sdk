package eppoclient

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
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

	return rule{}, errors.New("No matching rule")
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
