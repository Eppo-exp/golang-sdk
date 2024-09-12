package eppoclient

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var numericRule = rule{Conditions: []condition{
	{Operator: "GT", Value: 10.0, Attribute: "age"},
	{Operator: "LT", Value: 100.0, Attribute: "age"},
}}

var semverRule = rule{Conditions: []condition{
	{Operator: "GTE", Value: "1.2.0", Attribute: "appVersion"},
	{Operator: "LT", Value: "2.2.0", Attribute: "appVersion"},
}}

var textRule = rule{Conditions: []condition{
	{Operator: "MATCHES", Value: ".*@email.com", Attribute: "email"},
}}

var ruleWithEmptyConditions = rule{Conditions: []condition{}}

func init() {
	numericRule.precompute()
	semverRule.precompute()
	textRule.precompute()
}

func Test_TextRule_NoMatch(t *testing.T) {
	subjectAttributes := make(Attributes)
	subjectAttributes["age"] = 99
	subjectAttributes["country"] = "US"
	subjectAttributes["email"] = "test@example.com"

	assert.False(t, textRule.matches(subjectAttributes))
}

func Test_numericRule_Success(t *testing.T) {
	subjectAttributes := make(Attributes)
	subjectAttributes["age"] = 99.0

	assert.True(t, numericRule.matches(subjectAttributes))
}

func Test_numericRule_WithString(t *testing.T) {
	subjectAttributes := make(Attributes)
	subjectAttributes["age"] = "99.0"

	assert.True(t, numericRule.matches(subjectAttributes))
}

func Test_semverRule_Success(t *testing.T) {
	subjectAttributes := make(Attributes)
	subjectAttributes["age"] = 99.0
	subjectAttributes["appVersion"] = "1.15.0"

	assert.True(t, semverRule.matches(subjectAttributes))
}

func Test_numericRule_NoAttributeForcondition(t *testing.T) {
	subjectAttributes := make(Attributes)
	assert.False(t, numericRule.matches(subjectAttributes))
}

func Test_ruleWithEmptycondition_NoConditionsForRule(t *testing.T) {
	subjectAttributes := make(Attributes)

	assert.True(t, ruleWithEmptyConditions.matches(subjectAttributes))
}

func Test_numericRule_NumericOperatorWithString(t *testing.T) {
	subjectAttributes := make(Attributes)
	subjectAttributes["age"] = "something"

	assert.False(t, numericRule.matches(subjectAttributes))
}

func Test_regex_NumericValueAndRegex(t *testing.T) {
	rule := rule{Conditions: []condition{{Operator: "MATCHES", Value: "[0-9]+", Attribute: "age"}}}

	subjectAttributes := make(Attributes)
	subjectAttributes["age"] = 99

	result := rule.matches(subjectAttributes)

	assert.True(t, result)
}

type MatchesRuleTest []struct {
	attributes Attributes
	rule       rule
	expected   bool
}

func (tests MatchesRuleTest) run(t *testing.T) {
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			result := tt.rule.matches(tt.attributes)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_ruleMatches_oneOfOperatorWithBoolean(t *testing.T) {
	oneOfRule := rule{Conditions: []condition{{Operator: "ONE_OF", Value: []string{"true"}, Attribute: "enabled"}}}
	notOneOfRule := rule{Conditions: []condition{{Operator: "NOT_ONE_OF", Value: []string{"true"}, Attribute: "enabled"}}}

	subjectAttributesEnabled := make(Attributes)
	subjectAttributesEnabled["enabled"] = "true"

	subjectAttributesDisabled := make(Attributes)
	subjectAttributesDisabled["enabled"] = "false"

	var tests = MatchesRuleTest{
		{subjectAttributesEnabled, oneOfRule, true},
		{subjectAttributesDisabled, oneOfRule, false},
		{subjectAttributesEnabled, notOneOfRule, false},
		{subjectAttributesDisabled, notOneOfRule, true},
	}
	tests.run(t)
}

func Test_ruleMatches_OneOfOperatorCaseSensitive(t *testing.T) {
	oneOfRule := rule{Conditions: []condition{{Operator: "ONE_OF", Value: []string{"1Ab", "Ron"}, Attribute: "name"}}}

	subjectAttributes0 := make(Attributes)
	subjectAttributes0["name"] = "ron"

	subjectAttributes1 := make(Attributes)
	subjectAttributes1["name"] = "1AB"

	MatchesRuleTest{
		{subjectAttributes0, oneOfRule, false},
		{subjectAttributes1, oneOfRule, false},
	}.run(t)
}

func Test_ruleMatches_NotOneOfOperatorCaseSensitive(t *testing.T) {
	notOneOfRule := rule{Conditions: []condition{{Operator: "NOT_ONE_OF", Value: []string{"bbB", "1.1.ab"}, Attribute: "name"}}}
	subjectAttributes0 := make(Attributes)
	subjectAttributes0["name"] = "BBB"

	subjectAttributes1 := make(Attributes)
	subjectAttributes1["name"] = "1.1.AB"

	MatchesRuleTest{
		{subjectAttributes0, notOneOfRule, true},
		{subjectAttributes1, notOneOfRule, true},
	}.run(t)
}

func Test_ruleMatches_OneOfOperatorWithString(t *testing.T) {
	oneOfRule := rule{Conditions: []condition{{Operator: "ONE_OF", Value: []string{"john", "ron"}, Attribute: "name"}}}
	notOneOfRule := rule{Conditions: []condition{{Operator: "NOT_ONE_OF", Value: []string{"ron"}, Attribute: "name"}}}

	subjectAttributesJohn := make(Attributes)
	subjectAttributesJohn["name"] = "john"

	subjectAttributesRon := make(Attributes)
	subjectAttributesRon["name"] = "ron"

	subjectAttributesSam := make(Attributes)
	subjectAttributesSam["name"] = "sam"

	MatchesRuleTest{
		{subjectAttributesJohn, oneOfRule, true},
		{subjectAttributesRon, oneOfRule, true},
		{subjectAttributesSam, oneOfRule, false},
		{subjectAttributesRon, notOneOfRule, false},
		{subjectAttributesSam, notOneOfRule, true},
	}.run(t)
}

func Test_matchesRule_OneOfOperatorWithNumber(t *testing.T) {
	oneOfRule := rule{Conditions: []condition{{Operator: "ONE_OF", Value: []string{"14", "15.11", "15"}, Attribute: "number"}}}
	notOneOfRule := rule{Conditions: []condition{{Operator: "NOT_ONE_OF", Value: []string{"10"}, Attribute: "number"}}}

	subjectAttributes0 := make(Attributes)
	subjectAttributes0["number"] = "14"

	subjectAttributes1 := make(Attributes)
	subjectAttributes1["number"] = 15.11

	subjectAttributes2 := make(Attributes)
	subjectAttributes2["number"] = 15

	subjectAttributes3 := make(Attributes)
	subjectAttributes3["number"] = "10"

	subjectAttributes4 := make(Attributes)
	subjectAttributes4["number"] = 11

	MatchesRuleTest{
		{subjectAttributes0, oneOfRule, true},
		{subjectAttributes1, oneOfRule, true},
		{subjectAttributes2, oneOfRule, true},
		{subjectAttributes3, oneOfRule, false},
		{subjectAttributes3, notOneOfRule, false},
		{subjectAttributes4, notOneOfRule, true},
	}.run(t)
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
	result := !isOneOf("D", []string{"A", "B", "C"})

	assert.Equal(t, expected, result)
}

func Test_isNotOneOf_Fail(t *testing.T) {
	expected := false
	result := !isOneOf("A", []string{"A", "B", "C"})

	assert.Equal(t, expected, result)
}

func Test_evaluateNumericcondition_Fail(t *testing.T) {
	expected := false
	result, err := evaluateNumericCondition(40, 30.0, condition{Operator: "LT", Value: 30.0})
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func Test_evaluateNumericcondition_Success(t *testing.T) {
	expected := true
	result, err := evaluateNumericCondition(25, 30.0, condition{Operator: "LT", Value: 30.0})
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func Test_evaluateNumericcondition_IncorrectOperator(t *testing.T) {
	result, err := evaluateNumericCondition(25, 30.0, condition{Operator: "LTGT", Value: 30.0})
	assert.Error(t, err)
	assert.False(t, result)
}

func Test_isNull_missingAttribute(t *testing.T) {
	result := condition{Operator: "IS_NULL", Attribute: "name", Value: true}.matches(
		Attributes{})
	assert.True(t, result)
}
func Test_isNotNull_missingAttribute(t *testing.T) {
	result := condition{Operator: "IS_NULL", Attribute: "name", Value: false}.matches(
		Attributes{})
	assert.False(t, result)
}
func Test_isNull_nilAttribute(t *testing.T) {
	result := condition{Operator: "IS_NULL", Attribute: "name", Value: true}.matches(
		Attributes{
			"name": nil,
		})
	assert.True(t, result)
}
func Test_isNotNull_nilAttribute(t *testing.T) {
	result := condition{Operator: "IS_NULL", Attribute: "name", Value: false}.matches(
		Attributes{
			"name": nil,
		})
	assert.False(t, result)
}
func Test_isNull_attributePresent(t *testing.T) {
	result := condition{Operator: "IS_NULL", Attribute: "name", Value: true}.matches(
		Attributes{
			"name": "Alex",
		})
	assert.False(t, result)
}
func Test_isNotNull_attributePresent(t *testing.T) {
	result := condition{Operator: "IS_NULL", Attribute: "name", Value: false}.matches(
		Attributes{
			"name": "Alex",
		})
	assert.True(t, result)
}

func Test_handles_all_numeric_types(t *testing.T) {
	condition := condition{Operator: "GT", Attribute: "powerLevel", Value: "9000"}
	condition.precompute()

	// Floats
	assert.True(t, condition.matches(Attributes{"powerLevel": 9001.0}))
	assert.False(t, condition.matches(Attributes{"powerLevel": 9000.0}))
	assert.True(t, condition.matches(Attributes{"powerLevel": float64(9001)}))
	assert.False(t, condition.matches(Attributes{"powerLevel": float64(-9001.0)}))
	assert.True(t, condition.matches(Attributes{"powerLevel": float32(9001)}))
	assert.False(t, condition.matches(Attributes{"powerLevel": float32(8999)}))
	// Signed Integers
	assert.True(t, condition.matches(Attributes{"powerLevel": 9001}))
	assert.False(t, condition.matches(Attributes{"powerLevel": 9000}))
	assert.False(t, condition.matches(Attributes{"powerLevel": int8(1)}))
	assert.True(t, condition.matches(Attributes{"powerLevel": int16(9001)}))
	assert.False(t, condition.matches(Attributes{"powerLevel": int16(-9002)}))
	assert.True(t, condition.matches(Attributes{"powerLevel": int32(10000)}))
	assert.False(t, condition.matches(Attributes{"powerLevel": int32(0)}))
	assert.True(t, condition.matches(Attributes{"powerLevel": int64(9001)}))
	assert.False(t, condition.matches(Attributes{"powerLevel": int64(8999)}))
	// Unsigned Integers
	assert.False(t, condition.matches(Attributes{"powerLevel": uint8(1)}))
	assert.True(t, condition.matches(Attributes{"powerLevel": uint16(9001)}))
	assert.False(t, condition.matches(Attributes{"powerLevel": uint16(8999)}))
	assert.True(t, condition.matches(Attributes{"powerLevel": uint32(10000)}))
	assert.False(t, condition.matches(Attributes{"powerLevel": uint32(0)}))
	assert.True(t, condition.matches(Attributes{"powerLevel": uint64(9001)}))
	assert.False(t, condition.matches(Attributes{"powerLevel": uint64(8999)}))
	// Strings
	assert.True(t, condition.matches(Attributes{"powerLevel": "9001"}))
	assert.True(t, condition.matches(Attributes{"powerLevel": "9000.1"}))
	assert.False(t, condition.matches(Attributes{"powerLevel": "9000"}))
	assert.False(t, condition.matches(Attributes{"powerLevel": ".2"}))
}

func Test_invalid_numeric_types(t *testing.T) {
	condition := condition{Operator: "GT", Attribute: "powerLevel", Value: "9000"}
	condition.precompute()

	assert.False(t, condition.matches(Attributes{"powerLevel": "empty"}))
	assert.False(t, condition.matches(Attributes{"powerLevel": ""}))
	assert.False(t, condition.matches(Attributes{"powerLevel": false}))
	assert.False(t, condition.matches(Attributes{"powerLevel": true}))
}
