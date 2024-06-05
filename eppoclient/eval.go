package eppoclient

import (
	"errors"
	"fmt"
	"time"
)

func (flag flagConfiguration) verifyType(ty variationType) error {
	if flag.VariationType == ty {
		return nil
	} else {
		return fmt.Errorf("unexpected variation type (expected: %v, actual: %v)", ty, flag.VariationType)
	}
}

func (flag flagConfiguration) eval(subjectKey string, subjectAttributes SubjectAttributes) (interface{}, *AssignmentEvent, error) {
	if !flag.Enabled {
		return nil, nil, errors.New("the experiment or flag is not enabled")
	}

	now := time.Now()
	augmentedSubjectAttributes := augmentWithSubjectKey(subjectAttributes, subjectKey)

	// Go doesn't have a good/quick way to set defaults during
	// parsing, so defaulting is handled here
	totalShards := flag.TotalShards
	if totalShards == 0 {
		totalShards = 10000
	}

	var allocation *allocation
	var split *split
	for _, a := range flag.Allocations {
		s := a.findMatchingSplit(subjectKey, augmentedSubjectAttributes, totalShards, now)
		if s != nil {
			allocation, split = &a, s
			break
		}
	}
	if allocation == nil || split == nil {
		return nil, nil, errors.New("subject is not part of any allocation")
	}

	variation, ok := flag.Variations[split.VariationKey]
	if !ok {
		return nil, nil, fmt.Errorf("cannot find variation: %v", split.VariationKey)
	}

	assignmentValue, err := flag.VariationType.valueToAssignmentValue(variation.Value)
	if err != nil {
		return nil, nil, err
	}

	var assignmentEvent *AssignmentEvent
	if allocation.DoLog == nil || *allocation.DoLog {
		assignmentEvent = &AssignmentEvent{
			FeatureFlag:       flag.Key,
			Allocation:        allocation.Key,
			Experiment:        flag.Key + "-" + allocation.Key,
			Variation:         variation.Key,
			Subject:           subjectKey,
			SubjectAttributes: subjectAttributes,
			Timestamp:         now.UTC().Format(time.RFC3339),
			MetaData: map[string]string{
				"sdkLanguage": "go",
				"sdkVersion":  __version__,
			},
			ExtraLogging: split.ExtraLogging,
		}
	}

	return assignmentValue, assignmentEvent, nil
}

// Augment `subjectAttributes` by setting "id" attribute to
// `subjectKey` if "id" is not already present.
//
// This is used so that rules can reference subject key in coditions.
func augmentWithSubjectKey(subjectAttributes SubjectAttributes, subjectKey string) SubjectAttributes {
	_, hasId := subjectAttributes["id"]
	if hasId {
		return subjectAttributes
	}

	augmentedSubjectAttributes := make(map[string]interface{}, len(subjectAttributes))
	for k, v := range subjectAttributes {
		augmentedSubjectAttributes[k] = v
	}
	augmentedSubjectAttributes["id"] = subjectKey

	return augmentedSubjectAttributes
}

func (allocation allocation) findMatchingSplit(subjectKey string, augmentedSubjectAttributes SubjectAttributes, totalShards int64, now time.Time) *split {
	if !allocation.StartAt.IsZero() && now.Before(allocation.StartAt) {
		return nil
	}
	if !allocation.EndAt.IsZero() && now.After(allocation.EndAt) {
		return nil
	}

	matchesRule := false
	for _, rule := range allocation.Rules {
		if rule.matches(augmentedSubjectAttributes) {
			matchesRule = true
			break
		}
	}

	if len(allocation.Rules) > 0 && !matchesRule {
		// Forbidden by rules
		return nil
	}

	for _, split := range allocation.Splits {
		if split.matches(subjectKey, totalShards) {
			return &split
		}
	}

	return nil
}

func (split split) matches(subjectKey string, totalShards int64) bool {
	for _, shard := range split.Shards {
		if !shard.matches(subjectKey, totalShards) {
			return false
		}
	}
	return true
}

func (shard shard) matches(subjectKey string, totalShards int64) bool {
	s := getShard(shard.Salt+"-"+subjectKey, totalShards)
	for _, r := range shard.Ranges {
		if isShardInRange(s, r) {
			return true
		}
	}
	return false
}
