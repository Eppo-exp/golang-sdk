package eppoclient

import (
	"testing"
)

var greaterThanCondition = Condition{operator: "GT", value: 10.0, attribute: "age"}
var lessThanCondition = Condition{operator: "LT", value: 100.0, attribute: "age"}
var numericRule = Rule{conditions: []Condition{greaterThanCondition, lessThanCondition}}

var matchesEmailCondition = Condition{operator: "MATCHES", value: ".*@email.com", attribute: "email"}
var textRule = Rule{conditions: []Condition{matchesEmailCondition}}
var ruleWithEmptyConditions = Rule{conditions: []Condition{}}

func Test_matchesAnyRule_withEmptyRules(t *testing.T) {
	expected := false

	subjectAttributes := make(Dictionary)
	subjectAttributes["age"] = 20
	subjectAttributes["country"] = "US"

	result := matchesAnyRule(subjectAttributes, []Rule{})

	if expected != result {
		t.Errorf("\"matchesAnyRule()\" FAILED, expected -> %v, got -> %v", expected, result)
	} else {
		t.Logf("\"matchesAnyRule()\" SUCCEDED, expected -> %v, got -> %v", expected, result)
	}
}

func Test_matchesAnyRule_whenNoRulesMatch(t *testing.T) {
	expected := false

	subjectAttributes := make(Dictionary)
	subjectAttributes["age"] = 99
	subjectAttributes["country"] = "US"
	subjectAttributes["email"] = "test@example.com"

	result := matchesAnyRule(subjectAttributes, []Rule{textRule})

	if expected != result {
		t.Errorf("\"matchesAnyRule()\" FAILED, expected -> %v, got -> %v", expected, result)
	} else {
		t.Logf("\"matchesAnyRule()\" SUCCEDED, expected -> %v, got -> %v", expected, result)
	}
}

func Test_matchesAnyRule_Success(t *testing.T) {
	expected := true

	subjectAttributes := make(Dictionary)
	subjectAttributes["age"] = 99.0

	result := matchesAnyRule(subjectAttributes, []Rule{numericRule})

	if expected != result {
		t.Errorf("\"matchesAnyRule()\" FAILED, expected -> %v, got -> %v", expected, result)
	} else {
		t.Logf("\"matchesAnyRule()\" SUCCEDED, expected -> %v, got -> %v", expected, result)
	}
}

func Test_matchesAnyRule_NoAttributeForCondition(t *testing.T) {
	expected := false

	subjectAttributes := make(Dictionary)

	result := matchesAnyRule(subjectAttributes, []Rule{numericRule})

	if expected != result {
		t.Errorf("\"matchesAnyRule()\" FAILED, expected -> %v, got -> %v", expected, result)
	} else {
		t.Logf("\"matchesAnyRule()\" SUCCEDED, expected -> %v, got -> %v", expected, result)
	}
}

func Test_matchesAnyRule_NoConditionsForRule(t *testing.T) {
	expected := true

	subjectAttributes := make(Dictionary)

	result := matchesAnyRule(subjectAttributes, []Rule{ruleWithEmptyConditions})

	if expected != result {
		t.Errorf("\"matchesAnyRule()\" FAILED, expected -> %v, got -> %v", expected, result)
	} else {
		t.Logf("\"matchesAnyRule()\" SUCCEDED, expected -> %v, got -> %v", expected, result)
	}
}

func Test_matchesAnyRule_NumericOperatorWithString(t *testing.T) {
	expected := false

	subjectAttributes := make(Dictionary)
	subjectAttributes["age"] = "something"

	result := matchesAnyRule(subjectAttributes, []Rule{numericRule})

	if expected != result {
		t.Errorf("\"matchesAnyRule()\" FAILED, expected -> %v, got -> %v", expected, result)
	} else {
		t.Logf("\"matchesAnyRule()\" SUCCEDED, expected -> %v, got -> %v", expected, result)
	}
}

func Test_matchesAnyRule_NumericValueAndRegex(t *testing.T) {
	expected := true

	condition := Condition{operator: "MATCHES", value: "[0-9]+", attribute: "age"}
	rule := Rule{conditions: []Condition{condition}}

	subjectAttributes := make(Dictionary)
	subjectAttributes["age"] = 99

	result := matchesAnyRule(subjectAttributes, []Rule{rule})

	if expected != result {
		t.Errorf("\"matchesAnyRule()\" FAILED, expected -> %v, got -> %v", expected, result)
	} else {
		t.Logf("\"matchesAnyRule()\" SUCCEDED, expected -> %v, got -> %v", expected, result)
	}
}

func Test_matchesAnyRule_oneOfOperatorWithBoolean(t *testing.T) {
	oneOfRule := Rule{conditions: []Condition{{operator: "ONE_OF", value: []string{"true"}, attribute: "enabled"}}}
	notOneOfRule := Rule{conditions: []Condition{{operator: "NOT_ONE_OF", value: []string{"True"}, attribute: "enabled"}}}

	subjectAttributes := make(Dictionary)
	subjectAttributes["enabled"] = "true"

	expected := true
	result := matchesAnyRule(subjectAttributes, []Rule{oneOfRule})

	if result != expected {
		t.Errorf("\"getMatchingStringValues()\" FAILED, expected -> %v, got -> %v", expected, result)
	} else {
		t.Logf("\"getMatchingStringValues()\" SUCCEDED, expected -> %v, got -> %v", expected, result)
	}

	subjectAttributes["enabled"] = "false"

	expected = false
	result = matchesAnyRule(subjectAttributes, []Rule{oneOfRule})

	if result != expected {
		t.Errorf("\"getMatchingStringValues()\" FAILED, expected -> %v, got -> %v", expected, result)
	} else {
		t.Logf("\"getMatchingStringValues()\" SUCCEDED, expected -> %v, got -> %v", expected, result)
	}

	subjectAttributes["enabled"] = "true"

	expected = false
	result = matchesAnyRule(subjectAttributes, []Rule{notOneOfRule})

	if result != expected {
		t.Errorf("\"getMatchingStringValues()\" FAILED, expected -> %v, got -> %v", expected, result)
	} else {
		t.Logf("\"getMatchingStringValues()\" SUCCEDED, expected -> %v, got -> %v", expected, result)
	}

	subjectAttributes["enabled"] = "false"

	expected = true
	result = matchesAnyRule(subjectAttributes, []Rule{notOneOfRule})

	if result != expected {
		t.Errorf("\"getMatchingStringValues()\" FAILED, expected -> %v, got -> %v", expected, result)
	} else {
		t.Logf("\"getMatchingStringValues()\" SUCCEDED, expected -> %v, got -> %v", expected, result)
	}
}

func Test_matchesAnyRule_OneOfOperatorCaseInsensitive(t *testing.T) {
	oneOfRule := Rule{conditions: []Condition{{operator: "ONE_OF", value: []string{"1Ab", "Ron"}, attribute: "name"}}}
	expected := true

	subjectAttributes := make(Dictionary)
	subjectAttributes["name"] = "ron"

	result := matchesAnyRule(subjectAttributes, []Rule{oneOfRule})

	if result != expected {
		t.Errorf("\"getMatchingStringValues()\" FAILED, expected -> %v, got -> %v", expected, result)
	} else {
		t.Logf("\"getMatchingStringValues()\" SUCCEDED, expected -> %v, got -> %v", expected, result)
	}

	subjectAttributes["name"] = "1AB"
	result = matchesAnyRule(subjectAttributes, []Rule{oneOfRule})

	if result != expected {
		t.Errorf("\"getMatchingStringValues()\" FAILED, expected -> %v, got -> %v", expected, result)
	} else {
		t.Logf("\"getMatchingStringValues()\" SUCCEDED, expected -> %v, got -> %v", expected, result)
	}
}

func Test_matchesAnyRule_NotOneOfOperatorCaseInsensitive(t *testing.T) {
	notOneOfRule := Rule{conditions: []Condition{{operator: "NOT_ONE_OF", value: []string{"bbB", "1.1.ab"}, attribute: "name"}}}
	expected := false

	subjectAttributes := make(Dictionary)
	subjectAttributes["name"] = "BBB"

	result := matchesAnyRule(subjectAttributes, []Rule{notOneOfRule})

	if result != expected {
		t.Errorf("\"getMatchingStringValues()\" FAILED, expected -> %v, got -> %v", expected, result)
	} else {
		t.Logf("\"getMatchingStringValues()\" SUCCEDED, expected -> %v, got -> %v", expected, result)
	}

	subjectAttributes["name"] = "1.1.AB"
	result = matchesAnyRule(subjectAttributes, []Rule{notOneOfRule})

	if result != expected {
		t.Errorf("\"getMatchingStringValues()\" FAILED, expected -> %v, got -> %v", expected, result)
	} else {
		t.Logf("\"getMatchingStringValues()\" SUCCEDED, expected -> %v, got -> %v", expected, result)
	}
}

func Test_matchesAnyRule_OneOfOperatorWithString(t *testing.T) {
	oneOfRule := Rule{conditions: []Condition{{operator: "ONE_OF", value: []string{"john", "ron"}, attribute: "name"}}}
	notOneOfRule := Rule{conditions: []Condition{{operator: "NOT_ONE_OF", value: []string{"ron"}, attribute: "name"}}}

	subjectAttributes := make(Dictionary)
	subjectAttributes["name"] = "john"

	expected := true
	result := matchesAnyRule(subjectAttributes, []Rule{oneOfRule})

	if result != expected {
		t.Errorf("\"getMatchingStringValues()\" FAILED, expected -> %v, got -> %v", expected, result)
	} else {
		t.Logf("\"getMatchingStringValues()\" SUCCEDED, expected -> %v, got -> %v", expected, result)
	}

	subjectAttributes["name"] = "ron"
	result = matchesAnyRule(subjectAttributes, []Rule{oneOfRule})

	if result != expected {
		t.Errorf("\"getMatchingStringValues()\" FAILED, expected -> %v, got -> %v", expected, result)
	} else {
		t.Logf("\"getMatchingStringValues()\" SUCCEDED, expected -> %v, got -> %v", expected, result)
	}

	expected = false
	subjectAttributes["name"] = "sam"
	result = matchesAnyRule(subjectAttributes, []Rule{oneOfRule})

	if result != expected {
		t.Errorf("\"getMatchingStringValues()\" FAILED, expected -> %v, got -> %v", expected, result)
	} else {
		t.Logf("\"getMatchingStringValues()\" SUCCEDED, expected -> %v, got -> %v", expected, result)
	}

	expected = false
	subjectAttributes["name"] = "ron"
	result = matchesAnyRule(subjectAttributes, []Rule{notOneOfRule})

	if result != expected {
		t.Errorf("\"getMatchingStringValues()\" FAILED, expected -> %v, got -> %v", expected, result)
	} else {
		t.Logf("\"getMatchingStringValues()\" SUCCEDED, expected -> %v, got -> %v", expected, result)
	}

	expected = true
	subjectAttributes["name"] = "sam"
	result = matchesAnyRule(subjectAttributes, []Rule{notOneOfRule})

	if result != expected {
		t.Errorf("\"getMatchingStringValues()\" FAILED, expected -> %v, got -> %v", expected, result)
	} else {
		t.Logf("\"getMatchingStringValues()\" SUCCEDED, expected -> %v, got -> %v", expected, result)
	}
}

func Test_matchesAnyRule_OneOfOperatorWithNumber(t *testing.T) {
	oneOfRule := Rule{conditions: []Condition{{operator: "ONE_OF", value: []string{"14", "15.11", "15"}, attribute: "number"}}}
	notOneOfRule := Rule{conditions: []Condition{{operator: "NOT_ONE_OF", value: []string{"10"}, attribute: "number"}}}

	subjectAttributes := make(Dictionary)
	subjectAttributes["number"] = "14"

	expected := true
	result := matchesAnyRule(subjectAttributes, []Rule{oneOfRule})

	if result != expected {
		t.Errorf("\"getMatchingStringValues()\" FAILED, expected -> %v, got -> %v", expected, result)
	} else {
		t.Logf("\"getMatchingStringValues()\" SUCCEDED, expected -> %v, got -> %v", expected, result)
	}

	subjectAttributes["number"] = 15.11
	result = matchesAnyRule(subjectAttributes, []Rule{oneOfRule})

	if result != expected {
		t.Errorf("\"getMatchingStringValues()\" FAILED, expected -> %v, got -> %v", expected, result)
	} else {
		t.Logf("\"getMatchingStringValues()\" SUCCEDED, expected -> %v, got -> %v", expected, result)
	}

	subjectAttributes["number"] = 15
	result = matchesAnyRule(subjectAttributes, []Rule{oneOfRule})

	if result != expected {
		t.Errorf("\"getMatchingStringValues()\" FAILED, expected -> %v, got -> %v", expected, result)
	} else {
		t.Logf("\"getMatchingStringValues()\" SUCCEDED, expected -> %v, got -> %v", expected, result)
	}

	expected = false
	subjectAttributes["number"] = "10"
	result = matchesAnyRule(subjectAttributes, []Rule{oneOfRule})

	if result != expected {
		t.Errorf("\"getMatchingStringValues()\" FAILED, expected -> %v, got -> %v", expected, result)
	} else {
		t.Logf("\"getMatchingStringValues()\" SUCCEDED, expected -> %v, got -> %v", expected, result)
	}

	expected = false
	subjectAttributes["number"] = "10"
	result = matchesAnyRule(subjectAttributes, []Rule{notOneOfRule})

	if result != expected {
		t.Errorf("\"getMatchingStringValues()\" FAILED, expected -> %v, got -> %v", expected, result)
	} else {
		t.Logf("\"getMatchingStringValues()\" SUCCEDED, expected -> %v, got -> %v", expected, result)
	}

	expected = true
	subjectAttributes["number"] = 11
	result = matchesAnyRule(subjectAttributes, []Rule{notOneOfRule})

	if result != expected {
		t.Errorf("\"getMatchingStringValues()\" FAILED, expected -> %v, got -> %v", expected, result)
	} else {
		t.Logf("\"getMatchingStringValues()\" SUCCEDED, expected -> %v, got -> %v", expected, result)
	}
}

func Test_getMatchingStringValues_Success(t *testing.T) {
	expected := []string{"A"}
	result := getMatchingStringValues("A", []string{"A", "B", "C"})

	if !equal(expected, result) {
		t.Errorf("\"getMatchingStringValues()\" FAILED, expected -> %v, got -> %v", expected, result)
	} else {
		t.Logf("\"getMatchingStringValues()\" SUCCEDED, expected -> %v, got -> %v", expected, result)
	}
}

func Test_getMatchingStringValues_Fail(t *testing.T) {
	expected := []string{"B"}
	result := getMatchingStringValues("A", []string{"A", "B", "C"})

	if equal(expected, result) {
		t.Errorf("\"getMatchingStringValues()\" FAILED, expected -> %v, got -> %v", expected, result)
	} else {
		t.Logf("\"getMatchingStringValues()\" SUCCEDED, expected -> %v, got -> %v", expected, result)
	}
}

func Test_isOneOf_Success(t *testing.T) {
	expected := true
	result := isOneOf("A", []string{"A", "B", "C"})

	if expected != result {
		t.Errorf("\"isOneOf()\" FAILED, expected -> %v, got -> %v", expected, result)
	} else {
		t.Logf("\"isOneOf()\" SUCCEDED, expected -> %v, got -> %v", expected, result)
	}
}

func Test_isOneOf_Fail(t *testing.T) {
	expected := false
	result := isOneOf("D", []string{"A", "B", "C"})

	if expected != result {
		t.Errorf("\"isOneOf()\" FAILED, expected -> %v, got -> %v", expected, result)
	} else {
		t.Logf("\"isOneOf()\" SUCCEDED, expected -> %v, got -> %v", expected, result)
	}
}

func Test_isNotOneOf_Success(t *testing.T) {
	expected := true
	result := isNotOneOf("D", []string{"A", "B", "C"})

	if expected != result {
		t.Errorf("\"isNotOneOf()\" FAILED, expected -> %v, got -> %v", expected, result)
	} else {
		t.Logf("\"isNotOneOf()\" SUCCEDED, expected -> %v, got -> %v", expected, result)
	}
}

func Test_isNotOneOf_Fail(t *testing.T) {
	expected := false
	result := isNotOneOf("A", []string{"A", "B", "C"})

	if expected != result {
		t.Errorf("\"isNotOneOf()\" FAILED, expected -> %v, got -> %v", expected, result)
	} else {
		t.Logf("\"isNotOneOf()\" SUCCEDED, expected -> %v, got -> %v", expected, result)
	}
}

func Test_evaluateNumericCondition_Success(t *testing.T) {
	expected := false
	result := evaluateNumericCondition(40, Condition{operator: "LT", value: 30.0})

	if expected != result {
		t.Errorf("\"evaluateNumericCondition()\" FAILED, expected -> %v, got -> %v", expected, result)
	} else {
		t.Logf("\"evaluateNumericCondition()\" SUCCEDED, expected -> %v, got -> %v", expected, result)
	}
}

func Test_evaluateNumericCondition_Fail(t *testing.T) {
	expected := true
	result := evaluateNumericCondition(25, Condition{operator: "LT", value: 30.0})

	if expected != result {
		t.Errorf("\"evaluateNumericCondition()\" FAILED, expected -> %v, got -> %v", expected, result)
	} else {
		t.Logf("\"evaluateNumericCondition()\" SUCCEDED, expected -> %v, got -> %v", expected, result)
	}
}

func Test_evaluateNumericCondition_IncorrectOperator(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	evaluateNumericCondition(25, Condition{operator: "LTGT", value: 30.0})
}
