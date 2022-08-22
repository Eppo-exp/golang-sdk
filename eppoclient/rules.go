package eppoclient

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Condition struct {
	attribute string
	value     interface{}
	operator  string `validator:"regexp=^(MATCHES|GTE|GT|LTE|LT|ONE_OF|NOT_ONE_OF)$"`
}

type Rule struct {
	conditions []Condition
}

func matchesAnyRule(subjectAttributes Dictionary, rules []Rule) bool {
	for _, rule := range rules {
		if matchesRule(subjectAttributes, rule) {
			return true
		}
	}

	return false
}

func matchesRule(subjectAttributes Dictionary, rule Rule) bool {
	for _, condition := range rule.conditions {
		if !evaluateCondition(subjectAttributes, condition) {
			return false
		}
	}

	return true
}

func evaluateCondition(subjectAttributes Dictionary, condition Condition) bool {
	subjectValue := subjectAttributes[condition.attribute]

	if subjectValue != nil {
		if condition.operator == "MATCHES" {
			v := reflect.ValueOf(subjectValue)
			if v.Kind() != reflect.String {
				subjectValue = strconv.Itoa(subjectValue.(int))
			}
			r, _ := regexp.MatchString(condition.value.(string), subjectValue.(string))
			return r
		} else if condition.operator == "ONE_OF" {
			return isOneOf(subjectValue, condition.value.([]string))
		} else if condition.operator == "NOT_ONE_OF" {
			return isNotOneOf(subjectValue, condition.value.([]string))
		} else {
			return evaluateNumericCondition(subjectValue, condition)
		}
	}
	return false
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

func evaluateNumericCondition(subjectValue interface{}, condition Condition) bool {
	v := reflect.ValueOf(subjectValue)

	if v.Kind() == reflect.String {
		return false
	}

	if v.Kind() == reflect.Int {
		subjectValue = float64(subjectValue.(int))
	}

	switch condition.operator {
	case "GT":
		return subjectValue.(float64) > condition.value.(float64)
	case "GTE":
		return subjectValue.(float64) >= condition.value.(float64)
	case "LT":
		return subjectValue.(float64) < condition.value.(float64)
	case "LTE":
		return subjectValue.(float64) <= condition.value.(float64)
	default:
		panic("Incorrect condition operator")
	}
}
