package eppoclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var greaterThanCondition = condition{Operator: "GT", Value: 10.0, Attribute: "age"}
var lessThanCondition = condition{Operator: "LT", Value: 100.0, Attribute: "age"}
var numericRule = rule{Conditions: []condition{greaterThanCondition, lessThanCondition}}

var greaterThanAppVersionCondition = condition{Operator: "GTE", Value: "1.0.0", Attribute: "appVersion"}
var lessThanAppVersionCondition = condition{Operator: "LT", Value: "2.2.0", Attribute: "appVersion"}
var semverRule = rule{Conditions: []condition{greaterThanAppVersionCondition, lessThanAppVersionCondition}}

var matchesEmailCondition = condition{Operator: "MATCHES", Value: ".*@email.com", Attribute: "email"}
var textRule = rule{AllocationKey: "allocation-key", Conditions: []condition{matchesEmailCondition}}
var ruleWithEmptyConditions = rule{Conditions: []condition{}}
var expectedNoMatchErrorMessage = "no matching rule"

func Test_findMatchingRule_withEmptyRules(t *testing.T) {
	subjectAttributes := make(dictionary)
	subjectAttributes["age"] = 20
	subjectAttributes["country"] = "US"

	_, err := findMatchingRule(subjectAttributes, []rule{})

	assert.EqualError(t, err, expectedNoMatchErrorMessage)
}

func Test_findMatchingRule_whenNoRulesMatch(t *testing.T) {
	subjectAttributes := make(dictionary)
	subjectAttributes["age"] = 99
	subjectAttributes["country"] = "US"
	subjectAttributes["email"] = "test@example.com"

	_, err := findMatchingRule(subjectAttributes, []rule{textRule})

	assert.EqualError(t, err, expectedNoMatchErrorMessage)
}

func Test_findMatchingRule_Success(t *testing.T) {
	subjectAttributes := make(dictionary)
	subjectAttributes["age"] = 99.0

	result, _ := findMatchingRule(subjectAttributes, []rule{numericRule})

	assert.Equal(t, numericRule, result)
}

func Test_findMatchingSemVerRule_Success(t *testing.T) {
	subjectAttributes := make(dictionary)
	subjectAttributes["age"] = 99.0
	subjectAttributes["appVersion"] = "1.1.0"

	result, _ := findMatchingRule(subjectAttributes, []rule{semverRule})

	assert.Equal(t, semverRule, result)
}

func Test_findMatchingRule_NoAttributeForCondition(t *testing.T) {
	subjectAttributes := make(dictionary)

	_, err := findMatchingRule(subjectAttributes, []rule{numericRule})

	assert.EqualError(t, err, expectedNoMatchErrorMessage)
}

func Test_findMatchingRule_NoConditionsForRule(t *testing.T) {
	subjectAttributes := make(dictionary)

	result, _ := findMatchingRule(subjectAttributes, []rule{ruleWithEmptyConditions})

	assert.Equal(t, ruleWithEmptyConditions, result)
}

func Test_findMatchingRule_NumericOperatorWithString(t *testing.T) {
	subjectAttributes := make(dictionary)
	subjectAttributes["age"] = "something"

	_, err := findMatchingRule(subjectAttributes, []rule{numericRule})

	assert.EqualError(t, err, expectedNoMatchErrorMessage)
}

func Test_findMatchingRule_NumericValueAndRegex(t *testing.T) {
	cdn := condition{Operator: "MATCHES", Value: "[0-9]+", Attribute: "age"}
	rl := rule{Conditions: []condition{cdn}}

	subjectAttributes := make(dictionary)
	subjectAttributes["age"] = 99

	result, _ := findMatchingRule(subjectAttributes, []rule{rl})

	assert.Equal(t, rl, result)
}

type MatchesAnyRuleTest []struct {
	a             dictionary
	b             []rule
	expectedRule  rule
	expectedError string
}

func Test_findMatchingRule_oneOfOperatorWithBoolean(t *testing.T) {
	oneOfRule := rule{Conditions: []condition{{Operator: "ONE_OF", Value: []string{"true"}, Attribute: "enabled"}}}
	notOneOfRule := rule{Conditions: []condition{{Operator: "NOT_ONE_OF", Value: []string{"True"}, Attribute: "enabled"}}}

	subjectAttributesEnabled := make(dictionary)
	subjectAttributesEnabled["enabled"] = "true"

	subjectAttributesDisabled := make(dictionary)
	subjectAttributesDisabled["enabled"] = "false"

	var tests = MatchesAnyRuleTest{
		{subjectAttributesEnabled, []rule{oneOfRule}, oneOfRule, ""},
		{subjectAttributesDisabled, []rule{oneOfRule}, rule{}, expectedNoMatchErrorMessage},
		{subjectAttributesEnabled, []rule{notOneOfRule}, rule{}, expectedNoMatchErrorMessage},
		{subjectAttributesDisabled, []rule{notOneOfRule}, notOneOfRule, ""},
	}

	for _, tt := range tests {
		result, err := findMatchingRule(tt.a, tt.b)

		assert.Equal(t, tt.expectedRule, result)
		if tt.expectedError != "" {
			assert.EqualError(t, err, tt.expectedError)
		}
	}
}

func Test_findMatchingRule_OneOfOperatorCaseInsensitive(t *testing.T) {
	oneOfRule := rule{Conditions: []condition{{Operator: "ONE_OF", Value: []string{"1Ab", "Ron"}, Attribute: "name"}}}
	subjectAttributes0 := make(dictionary)
	subjectAttributes0["name"] = "ron"

	subjectAttributes1 := make(dictionary)
	subjectAttributes1["name"] = "1AB"

	var tests = MatchesAnyRuleTest{
		{subjectAttributes0, []rule{oneOfRule}, oneOfRule, ""},
		{subjectAttributes1, []rule{oneOfRule}, oneOfRule, ""},
	}

	for _, tt := range tests {
		result, err := findMatchingRule(tt.a, tt.b)

		assert.Equal(t, tt.expectedRule, result)
		if tt.expectedError != "" {
			assert.EqualError(t, err, tt.expectedError)
		}
	}
}

func Test_findMatchingRule_NotOneOfOperatorCaseInsensitive(t *testing.T) {
	notOneOfRule := rule{Conditions: []condition{{Operator: "NOT_ONE_OF", Value: []string{"bbB", "1.1.ab"}, Attribute: "name"}}}
	subjectAttributes0 := make(dictionary)
	subjectAttributes0["name"] = "BBB"

	subjectAttributes1 := make(dictionary)
	subjectAttributes1["name"] = "1.1.AB"

	var tests = MatchesAnyRuleTest{
		{subjectAttributes0, []rule{notOneOfRule}, rule{}, expectedNoMatchErrorMessage},
		{subjectAttributes1, []rule{notOneOfRule}, rule{}, expectedNoMatchErrorMessage},
	}

	for _, tt := range tests {
		result, err := findMatchingRule(tt.a, tt.b)

		assert.Equal(t, tt.expectedRule, result)
		assert.EqualError(t, err, tt.expectedError)
	}
}

func Test_findMatchingRule_OneOfOperatorWithString(t *testing.T) {
	oneOfRule := rule{Conditions: []condition{{Operator: "ONE_OF", Value: []string{"john", "ron"}, Attribute: "name"}}}
	notOneOfRule := rule{Conditions: []condition{{Operator: "NOT_ONE_OF", Value: []string{"ron"}, Attribute: "name"}}}

	subjectAttributesJohn := make(dictionary)
	subjectAttributesJohn["name"] = "john"

	subjectAttributesRon := make(dictionary)
	subjectAttributesRon["name"] = "ron"

	subjectAttributesSam := make(dictionary)
	subjectAttributesSam["name"] = "sam"

	var tests = MatchesAnyRuleTest{
		{subjectAttributesJohn, []rule{oneOfRule}, oneOfRule, ""},
		{subjectAttributesRon, []rule{oneOfRule}, oneOfRule, ""},
		{subjectAttributesSam, []rule{oneOfRule}, rule{}, expectedNoMatchErrorMessage},
		{subjectAttributesRon, []rule{notOneOfRule}, rule{}, expectedNoMatchErrorMessage},
		{subjectAttributesSam, []rule{notOneOfRule}, notOneOfRule, ""},
	}

	for _, tt := range tests {
		result, err := findMatchingRule(tt.a, tt.b)

		assert.Equal(t, tt.expectedRule, result)
		if tt.expectedError != "" {
			assert.EqualError(t, err, tt.expectedError)
		}
	}
}

func Test_findMatchingRule_OneOfOperatorWithNumber(t *testing.T) {
	oneOfRule := rule{Conditions: []condition{{Operator: "ONE_OF", Value: []string{"14", "15.11", "15"}, Attribute: "number"}}}
	notOneOfRule := rule{Conditions: []condition{{Operator: "NOT_ONE_OF", Value: []string{"10"}, Attribute: "number"}}}

	subjectAttributes0 := make(dictionary)
	subjectAttributes0["number"] = "14"

	subjectAttributes1 := make(dictionary)
	subjectAttributes1["number"] = 15.11

	subjectAttributes2 := make(dictionary)
	subjectAttributes2["number"] = 15

	subjectAttributes3 := make(dictionary)
	subjectAttributes3["number"] = "10"

	subjectAttributes4 := make(dictionary)
	subjectAttributes4["number"] = 11

	var tests = MatchesAnyRuleTest{
		{subjectAttributes0, []rule{oneOfRule}, oneOfRule, ""},
		{subjectAttributes1, []rule{oneOfRule}, oneOfRule, ""},
		{subjectAttributes2, []rule{oneOfRule}, oneOfRule, ""},
		{subjectAttributes3, []rule{oneOfRule}, rule{}, expectedNoMatchErrorMessage},
		{subjectAttributes3, []rule{notOneOfRule}, rule{}, expectedNoMatchErrorMessage},
		{subjectAttributes4, []rule{notOneOfRule}, notOneOfRule, ""},
	}

	for _, tt := range tests {
		result, err := findMatchingRule(tt.a, tt.b)

		assert.Equal(t, tt.expectedRule, result)
		if tt.expectedError != "" {
			assert.EqualError(t, err, tt.expectedError)
		}
	}
}

func Test_getMatchingStringValues_Success(t *testing.T) {
	expected := []string{"A"}
	result := getMatchingStringValues("A", []string{"A", "B", "C"})

	assert.Equal(t, expected, result)
}

func Test_getMatchingStringValues_Fail(t *testing.T) {
	expected := []string{"B"}
	result := getMatchingStringValues("A", []string{"A", "B", "C"})

	assert.NotEqual(t, expected, result)
}

func Test_isOneOf_Success(t *testing.T) {
	expected := true
	result := isOneOf("A", []string{"A", "B", "C"})

	assert.Equal(t, expected, result)
}

func Test_isOneOf_Fail(t *testing.T) {
	expected := false
	result := isOneOf("D", []string{"A", "B", "C"})

	assert.Equal(t, expected, result)
}

func Test_isNotOneOf_Success(t *testing.T) {
	expected := true
	result := isNotOneOf("D", []string{"A", "B", "C"})

	assert.Equal(t, expected, result)
}

func Test_isNotOneOf_Fail(t *testing.T) {
	expected := false
	result := isNotOneOf("A", []string{"A", "B", "C"})

	assert.Equal(t, expected, result)
}

func Test_evaluateNumericCondition_Success(t *testing.T) {
	expected := false
	result := evaluateNumericCondition(40, condition{Operator: "LT", Value: 30.0})

	assert.Equal(t, expected, result)
}

func Test_evaluateNumericCondition_Fail(t *testing.T) {
	expected := true
	result := evaluateNumericCondition(25, condition{Operator: "LT", Value: 30.0})

	assert.Equal(t, expected, result)
}

func Test_evaluateNumericCondition_IncorrectOperator(t *testing.T) {
	assert.Panics(t, func() { evaluateNumericCondition(25, condition{Operator: "LTGT", Value: 30.0}) })
}
