package eppoclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var greaterThanCondition = condition{Operator: "GT", Value: 10.0, Attribute: "age"}
var lessThanCondition = condition{Operator: "LT", Value: 100.0, Attribute: "age"}
var numericRule = rule{Conditions: []condition{greaterThanCondition, lessThanCondition}}

var matchesEmailCondition = condition{Operator: "MATCHES", Value: ".*@email.com", Attribute: "email"}
var textRule = rule{Conditions: []condition{matchesEmailCondition}}
var ruleWithEmptyConditions = rule{Conditions: []condition{}}

func Test_matchesAnyRule_withEmptyRules(t *testing.T) {
	expected := false

	subjectAttributes := make(dictionary)
	subjectAttributes["age"] = 20
	subjectAttributes["country"] = "US"

	result := matchesAnyRule(subjectAttributes, []rule{})

	assert.Equal(t, expected, result)
}

func Test_matchesAnyRule_whenNoRulesMatch(t *testing.T) {
	expected := false

	subjectAttributes := make(dictionary)
	subjectAttributes["age"] = 99
	subjectAttributes["country"] = "US"
	subjectAttributes["email"] = "test@example.com"

	result := matchesAnyRule(subjectAttributes, []rule{textRule})

	assert.Equal(t, expected, result)
}

func Test_matchesAnyRule_Success(t *testing.T) {
	expected := true

	subjectAttributes := make(dictionary)
	subjectAttributes["age"] = 99.0

	result := matchesAnyRule(subjectAttributes, []rule{numericRule})

	assert.Equal(t, expected, result)
}

func Test_matchesAnyRule_NoAttributeForCondition(t *testing.T) {
	expected := false

	subjectAttributes := make(dictionary)

	result := matchesAnyRule(subjectAttributes, []rule{numericRule})

	assert.Equal(t, expected, result)
}

func Test_matchesAnyRule_NoConditionsForRule(t *testing.T) {
	expected := true

	subjectAttributes := make(dictionary)

	result := matchesAnyRule(subjectAttributes, []rule{ruleWithEmptyConditions})

	assert.Equal(t, expected, result)
}

func Test_matchesAnyRule_NumericOperatorWithString(t *testing.T) {
	expected := false

	subjectAttributes := make(dictionary)
	subjectAttributes["age"] = "something"

	result := matchesAnyRule(subjectAttributes, []rule{numericRule})

	assert.Equal(t, expected, result)
}

func Test_matchesAnyRule_NumericValueAndRegex(t *testing.T) {
	expected := true

	cdn := condition{Operator: "MATCHES", Value: "[0-9]+", Attribute: "age"}
	rl := rule{Conditions: []condition{cdn}}

	subjectAttributes := make(dictionary)
	subjectAttributes["age"] = 99

	result := matchesAnyRule(subjectAttributes, []rule{rl})

	assert.Equal(t, expected, result)
}

type MatchesAnyRuleTest []struct {
	a    dictionary
	b    []rule
	want bool
}

func Test_matchesAnyRule_oneOfOperatorWithBoolean(t *testing.T) {
	oneOfRule := rule{Conditions: []condition{{Operator: "ONE_OF", Value: []string{"true"}, Attribute: "enabled"}}}
	notOneOfRule := rule{Conditions: []condition{{Operator: "NOT_ONE_OF", Value: []string{"True"}, Attribute: "enabled"}}}

	subjectAttributesEnabled := make(dictionary)
	subjectAttributesEnabled["enabled"] = "true"

	subjectAttributesDisabled := make(dictionary)
	subjectAttributesDisabled["enabled"] = "false"

	var tests = MatchesAnyRuleTest{
		{subjectAttributesEnabled, []rule{oneOfRule}, true},
		{subjectAttributesDisabled, []rule{oneOfRule}, false},
		{subjectAttributesEnabled, []rule{notOneOfRule}, false},
		{subjectAttributesDisabled, []rule{notOneOfRule}, true},
	}

	for _, tt := range tests {
		result := matchesAnyRule(tt.a, tt.b)

		assert.Equal(t, tt.want, result)
	}
}

func Test_matchesAnyRule_OneOfOperatorCaseInsensitive(t *testing.T) {
	oneOfRule := rule{Conditions: []condition{{Operator: "ONE_OF", Value: []string{"1Ab", "Ron"}, Attribute: "name"}}}
	subjectAttributes0 := make(dictionary)
	subjectAttributes0["name"] = "ron"

	subjectAttributes1 := make(dictionary)
	subjectAttributes1["name"] = "1AB"

	var tests = MatchesAnyRuleTest{
		{subjectAttributes0, []rule{oneOfRule}, true},
		{subjectAttributes1, []rule{oneOfRule}, true},
	}

	for _, tt := range tests {
		result := matchesAnyRule(tt.a, tt.b)

		assert.Equal(t, tt.want, result)
	}
}

func Test_matchesAnyRule_NotOneOfOperatorCaseInsensitive(t *testing.T) {
	notOneOfRule := rule{Conditions: []condition{{Operator: "NOT_ONE_OF", Value: []string{"bbB", "1.1.ab"}, Attribute: "name"}}}
	subjectAttributes0 := make(dictionary)
	subjectAttributes0["name"] = "BBB"

	subjectAttributes1 := make(dictionary)
	subjectAttributes1["name"] = "1.1.AB"

	var tests = MatchesAnyRuleTest{
		{subjectAttributes0, []rule{notOneOfRule}, false},
		{subjectAttributes1, []rule{notOneOfRule}, false},
	}

	for _, tt := range tests {
		result := matchesAnyRule(tt.a, tt.b)

		assert.Equal(t, tt.want, result)
	}
}

func Test_matchesAnyRule_OneOfOperatorWithString(t *testing.T) {
	oneOfRule := rule{Conditions: []condition{{Operator: "ONE_OF", Value: []string{"john", "ron"}, Attribute: "name"}}}
	notOneOfRule := rule{Conditions: []condition{{Operator: "NOT_ONE_OF", Value: []string{"ron"}, Attribute: "name"}}}

	subjectAttributesJohn := make(dictionary)
	subjectAttributesJohn["name"] = "john"

	subjectAttributesRon := make(dictionary)
	subjectAttributesRon["name"] = "ron"

	subjectAttributesSam := make(dictionary)
	subjectAttributesSam["name"] = "sam"

	var tests = MatchesAnyRuleTest{
		{subjectAttributesJohn, []rule{oneOfRule}, true},
		{subjectAttributesRon, []rule{oneOfRule}, true},
		{subjectAttributesSam, []rule{oneOfRule}, false},
		{subjectAttributesRon, []rule{notOneOfRule}, false},
		{subjectAttributesSam, []rule{notOneOfRule}, true},
	}

	for _, tt := range tests {
		result := matchesAnyRule(tt.a, tt.b)

		assert.Equal(t, tt.want, result)
	}
}

func Test_matchesAnyRule_OneOfOperatorWithNumber(t *testing.T) {
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
		{subjectAttributes0, []rule{oneOfRule}, true},
		{subjectAttributes1, []rule{oneOfRule}, true},
		{subjectAttributes2, []rule{oneOfRule}, true},
		{subjectAttributes3, []rule{oneOfRule}, false},
		{subjectAttributes3, []rule{notOneOfRule}, false},
		{subjectAttributes4, []rule{notOneOfRule}, true},
	}

	for _, tt := range tests {
		result := matchesAnyRule(tt.a, tt.b)

		assert.Equal(t, tt.want, result)
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
