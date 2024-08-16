package eppoclient

import (
	"math"
	"sort"
	"strconv"
)

type ContextAttributes struct {
	Numeric     map[string]float64
	Categorical map[string]string
}

// Tries to map generic attributes to ContextAttributes depending on attribute types.
// - Integer and float types are mapped to numeric attributes.
// - Strings and bools are mapped to categorical attributes.
// - Rest of types are silently dropped.
func InferContextAttributes(attrs map[string]interface{}) ContextAttributes {
	result := ContextAttributes{
		Numeric:     map[string]float64{},
		Categorical: map[string]string{},
	}
	for key, value := range attrs {
		switch value := value.(type) {
		case int:
			result.Numeric[key] = float64(value)
		case int8:
			result.Numeric[key] = float64(value)
		case int16:
			result.Numeric[key] = float64(value)
		case int32:
			result.Numeric[key] = float64(value)
		case int64:
			result.Numeric[key] = float64(value)
		case uint:
			result.Numeric[key] = float64(value)
		case uint8:
			result.Numeric[key] = float64(value)
		case uint16:
			result.Numeric[key] = float64(value)
		case uint32:
			result.Numeric[key] = float64(value)
		case uint64:
			result.Numeric[key] = float64(value)
		case float32:
			result.Numeric[key] = float64(value)
		case float64:
			result.Numeric[key] = value
		case string:
			result.Categorical[key] = value
		case bool:
			result.Categorical[key] = strconv.FormatBool(value)
		}
	}
	return result
}

func (self ContextAttributes) toGenericAttributes() Attributes {
	result := make(Attributes)
	for key, value := range self.Numeric {
		result[key] = value
	}
	for key, value := range self.Categorical {
		result[key] = value
	}
	return result
}

type banditEvaluationContext struct {
	flagKey           string
	subjectKey        string
	subjectAttributes ContextAttributes
	actions           map[string]ContextAttributes
}

type banditEvaluationDetails struct {
	flagKey           string
	subjectKey        string
	subjectAttributes ContextAttributes
	actionKey         string
	actionAttributes  ContextAttributes
	actionScore       float64
	actionWeight      float64
	gamma             float64
	optimalityGap     float64
}

type action struct {
	key        string
	attributes ContextAttributes
}

func (model *banditModelData) evaluate(ctx banditEvaluationContext) banditEvaluationDetails {
	// There's currently no way to change totalShards in bandit evaluation.
	var totalShards int64 = 10_000

	nActions := len(ctx.actions)

	scores := make(map[string]float64, nActions)
	for actionKey, actionAttributes := range ctx.actions {
		scores[actionKey] = model.scoreAction(ctx.subjectAttributes, action{key: actionKey, attributes: actionAttributes})
	}

	bestAction, bestScore := "", math.Inf(-1)
	for actionKey, score := range scores {
		if score > bestScore || (score == bestScore && actionKey < bestAction) {
			bestAction, bestScore = actionKey, score
		}
	}

	weights := make(map[string]float64, nActions)
	{
		for actionKey, score := range scores {
			if actionKey == bestAction {
				// best action is assigned the remainder weight
				continue
			}

			// adjust probability floor for number of actions to control the sum
			minProbability := model.ActionProbabilityFloor / float64(nActions)
			weights[actionKey] = math.Max(minProbability, 1.0/(float64(nActions)+model.Gamma*(bestScore-score)))
		}

		remainderWeight := 1.0
		for _, weight := range weights {
			remainderWeight -= weight
		}
		weights[bestAction] = math.Max(0.0, remainderWeight)
	}

	// Pseudo-random deterministic shuffle of actions.
	shuffledActions := make([]string, 0, nActions)
	{
		for actionKey := range ctx.actions {
			shuffledActions = append(shuffledActions, actionKey)
		}

		shards := make(map[string]int64, nActions)
		for actionKey := range ctx.actions {
			shards[actionKey] = getShard(ctx.flagKey+"-"+ctx.subjectKey+"-"+actionKey, totalShards)
		}

		// Sort actions by their shard value. Use action key
		// as tie breaker.
		sort.Slice(shuffledActions, func(i, j int) bool {
			a1 := shuffledActions[i]
			a2 := shuffledActions[j]
			v1 := shards[a1]
			v2 := shards[a2]
			if v1 < v2 {
				return true
			} else if v1 > v2 {
				return false
			} else {
				// tie-breaking
				return a1 < a2
			}
		})
	}

	shardValue := float64(getShard(ctx.flagKey+"-"+ctx.subjectKey, totalShards)) / float64(totalShards)

	cumulativeWeight := 0.0
	var selectedAction string
	for _, selectedAction = range shuffledActions {
		cumulativeWeight += weights[selectedAction]
		if cumulativeWeight > shardValue {
			break
		}
	}

	optimalityGap := bestScore - scores[selectedAction]

	return banditEvaluationDetails{
		flagKey:           ctx.flagKey,
		subjectKey:        ctx.subjectKey,
		subjectAttributes: ctx.subjectAttributes,
		actionKey:         selectedAction,
		actionAttributes:  ctx.actions[selectedAction],
		actionScore:       scores[selectedAction],
		actionWeight:      weights[selectedAction],
		gamma:             model.Gamma,
		optimalityGap:     optimalityGap,
	}
}

func (model *banditModelData) scoreAction(subjectAttributes ContextAttributes, action action) float64 {
	coefficients, hasCoefficients := model.Coefficients[action.key]
	if !hasCoefficients {
		return model.DefaultActionScore
	}

	score := coefficients.Intercept
	score += scoreNumericAttributes(coefficients.ActionNumericCoefficients, action.attributes.Numeric)
	score += scoreCategoricalAttributes(coefficients.ActionCategoricalCoefficients, action.attributes.Categorical)
	score += scoreNumericAttributes(coefficients.SubjectNumericCoefficients, subjectAttributes.Numeric)
	score += scoreCategoricalAttributes(coefficients.SubjectCategoricalCoefficients, subjectAttributes.Categorical)
	return score
}

func scoreNumericAttributes(coefficients []banditNumericAttributeCoefficient, attributes map[string]float64) float64 {
	score := 0.0
	for _, coefficient := range coefficients {
		attribute, hasAttribute := attributes[coefficient.AttributeKey]
		if hasAttribute {
			score += coefficient.Coefficient * attribute
		} else {
			score += coefficient.MissingValueCoefficient
		}
	}
	return score
}

func scoreCategoricalAttributes(coefficients []banditCategoricalAttributeCoefficient, attributes map[string]string) float64 {
	score := 0.0
	for _, coefficient := range coefficients {
		attribute, hasAttribute := attributes[coefficient.AttributeKey]
		if hasAttribute {
			valueCoefficient, hasValueCoefficient := coefficient.ValueCoefficients[attribute]
			if hasValueCoefficient {
				score += valueCoefficient
			} else {
				score += coefficient.MissingValueCoefficient
			}
		} else {
			score += coefficient.MissingValueCoefficient
		}
	}
	return score
}
