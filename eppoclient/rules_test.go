package eppoclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var greaterThanCondition = condition{operator: "GT", value: 10.0, attribute: "age"}
var lessThanCondition = condition{operator: "LT", value: 100.0, attribute: "age"}
var numericRule = rule{conditions: []condition{greaterThanCondition, lessThanCondition}}

var matchesEmailCondition = condition{operator: "MATCHES", value: ".*@email.com", attribute: "email"}
var textRule = rule{conditions: []condition{matchesEmailCondition}}
var ruleWithEmptyConditions = rule{conditions: []condition{}}

func Test_matchesAnyRule_withEmptyRules(t *testing.T) {
	expected := false

	subjectAttributes := make(Dictionary)
	subjectAttributes["age"] = 20
	subjectAttributes["country"] = "US"

	result := matchesAnyRule(subjectAttributes, []rule{})

	assert.Equal(t, expected, result)
}

func Test_matchesAnyRule_whenNoRulesMatch(t *testing.T) {
	expected := false

	subjectAttributes := make(Dictionary)
	subjectAttributes["age"] = 99
	subjectAttributes["country"] = "US"
	subjectAttributes["email"] = "test@example.com"

	result := matchesAnyRule(subjectAttributes, []rule{textRule})

	assert.Equal(t, expected, result)
}

func Test_matchesAnyRule_Success(t *testing.T) {
	expected := true

	subjectAttributes := make(Dictionary)
	subjectAttributes["age"] = 99.0

	result := matchesAnyRule(subjectAttributes, []rule{numericRule})

	assert.Equal(t, expected, result)
}

func Test_matchesAnyRule_NoAttributeForCondition(t *testing.T) {
	expected := false

	subjectAttributes := make(Dictionary)

	result := matchesAnyRule(subjectAttributes, []rule{numericRule})

	assert.Equal(t, expected, result)
}

func Test_matchesAnyRule_NoConditionsForRule(t *testing.T) {
	expected := true

	subjectAttributes := make(Dictionary)

	result := matchesAnyRule(subjectAttributes, []rule{ruleWithEmptyConditions})

	assert.Equal(t, expected, result)
}

func Test_matchesAnyRule_NumericOperatorWithString(t *testing.T) {
	expected := false

	subjectAttributes := make(Dictionary)
	subjectAttributes["age"] = "something"

	result := matchesAnyRule(subjectAttributes, []rule{numericRule})

	assert.Equal(t, expected, result)
}

func Test_matchesAnyRule_NumericValueAndRegex(t *testing.T) {
	expected := true

	cdn := condition{operator: "MATCHES", value: "[0-9]+", attribute: "age"}
	rl := rule{conditions: []condition{cdn}}

	subjectAttributes := make(Dictionary)
	subjectAttributes["age"] = 99

	result := matchesAnyRule(subjectAttributes, []rule{rl})

	assert.Equal(t, expected, result)
}

type MatchesAnyRuleTest []struct {
	a    Dictionary
	b    []rule
	want bool
}

func Test_matchesAnyRule_oneOfOperatorWithBoolean(t *testing.T) {
	oneOfRule := rule{conditions: []condition{{operator: "ONE_OF", value: []string{"true"}, attribute: "enabled"}}}
	notOneOfRule := rule{conditions: []condition{{operator: "NOT_ONE_OF", value: []string{"True"}, attribute: "enabled"}}}

	subjectAttributesEnabled := make(Dictionary)
	subjectAttributesEnabled["enabled"] = "true"

	subjectAttributesDisabled := make(Dictionary)
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
	oneOfRule := rule{conditions: []condition{{operator: "ONE_OF", value: []string{"1Ab", "Ron"}, attribute: "name"}}}
	subjectAttributes0 := make(Dictionary)
	subjectAttributes0["name"] = "ron"

	subjectAttributes1 := make(Dictionary)
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
	notOneOfRule := rule{conditions: []condition{{operator: "NOT_ONE_OF", value: []string{"bbB", "1.1.ab"}, attribute: "name"}}}
	subjectAttributes0 := make(Dictionary)
	subjectAttributes0["name"] = "BBB"

	subjectAttributes1 := make(Dictionary)
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
	oneOfRule := rule{conditions: []condition{{operator: "ONE_OF", value: []string{"john", "ron"}, attribute: "name"}}}
	notOneOfRule := rule{conditions: []condition{{operator: "NOT_ONE_OF", value: []string{"ron"}, attribute: "name"}}}

	subjectAttributesJohn := make(Dictionary)
	subjectAttributesJohn["name"] = "john"

	subjectAttributesRon := make(Dictionary)
	subjectAttributesRon["name"] = "ron"

	subjectAttributesSam := make(Dictionary)
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
	oneOfRule := rule{conditions: []condition{{operator: "ONE_OF", value: []string{"14", "15.11", "15"}, attribute: "number"}}}
	notOneOfRule := rule{conditions: []condition{{operator: "NOT_ONE_OF", value: []string{"10"}, attribute: "number"}}}

	subjectAttributes0 := make(Dictionary)
	subjectAttributes0["number"] = "14"

	subjectAttributes1 := make(Dictionary)
	subjectAttributes1["number"] = 15.11

	subjectAttributes2 := make(Dictionary)
	subjectAttributes2["number"] = 15

	subjectAttributes3 := make(Dictionary)
	subjectAttributes3["number"] = "10"

	subjectAttributes4 := make(Dictionary)
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
	result := evaluateNumericCondition(40, condition{operator: "LT", value: 30.0})

	assert.Equal(t, expected, result)
}

func Test_evaluateNumericCondition_Fail(t *testing.T) {
	expected := true
	result := evaluateNumericCondition(25, condition{operator: "LT", value: 30.0})

	assert.Equal(t, expected, result)
}

func Test_evaluateNumericCondition_IncorrectOperator(t *testing.T) {
	assert.Panics(t, func() { evaluateNumericCondition(25, condition{operator: "LTGT", value: 30.0}) })
}
