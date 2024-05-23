package eppoclient

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
)

// Client for eppo.cloud. Instance of this struct will be created on calling InitClient.
// EppoClient will then immediately start polling experiments data from Eppo.
type EppoClient struct {
	configRequestor iConfigRequestor
	poller          poller
	logger          IAssignmentLogger
}

func newEppoClient(configRequestor iConfigRequestor, assignmentLogger IAssignmentLogger) *EppoClient {
	var ec = &EppoClient{}

	var poller = newPoller(10, configRequestor.FetchAndStoreConfigurations)
	ec.poller = *poller
	ec.configRequestor = configRequestor
	ec.logger = assignmentLogger

	return ec
}

// GetAssignment is maintained for backwards capability. It will return a string value for the assignment.
func (ec *EppoClient) GetAssignment(subjectKey string, flagKey string, subjectAttributes dictionary) (string, error) {
	return ec.GetStringAssignment(subjectKey, flagKey, subjectAttributes)
}

func (ec *EppoClient) GetBoolAssignment(subjectKey string, flagKey string, subjectAttributes dictionary) (bool, error) {
	variation, err := ec.getAssignment(subjectKey, flagKey, subjectAttributes, BoolType)
	return variation.BoolValue, err
}

func (ec *EppoClient) GetNumericAssignment(subjectKey string, flagKey string, subjectAttributes dictionary) (float64, error) {
	variation, err := ec.getAssignment(subjectKey, flagKey, subjectAttributes, NumericType)
	return variation.NumericValue, err
}

func (ec *EppoClient) GetStringAssignment(subjectKey string, flagKey string, subjectAttributes dictionary) (string, error) {
	variation, err := ec.getAssignment(subjectKey, flagKey, subjectAttributes, StringType)
	return variation.StringValue, err
}

func (ec *EppoClient) GetJSONStringAssignment(subjectKey string, flagKey string, subjectAttributes dictionary) (string, error) {
	variation, err := ec.getAssignment(subjectKey, flagKey, subjectAttributes, StringType)
	return variation.StringValue, err
}

func (ec *EppoClient) getAssignment(subjectKey string, flagKey string, subjectAttributes dictionary, valueType ValueType) (Value, error) {
	if subjectKey == "" {
		panic("no subject key provided")
	}

	if flagKey == "" {
		panic("no flag key provided")
	}

	config, err := ec.configRequestor.GetConfiguration(flagKey)
	if err != nil {
		return Null(), err
	}

	override := getSubjectVariationOverride(config, subjectKey, valueType)
	if override != Null() {
		return override, nil
	}

	// Check if disabled
	if !config.Enabled {
		return Null(), errors.New("the experiment or flag is not enabled")
	}

	// Find matching rule
	rule, err := findMatchingRule(subjectAttributes, config.Rules)
	if err != nil {
		return Null(), err
	}

	// Check if in sample population
	allocation := config.Allocations[rule.AllocationKey]
	if !isInExperimentSample(subjectKey, flagKey, config.SubjectShards, allocation.PercentExposure) {
		return Null(), errors.New("subject not part of the sample population")
	}

	// Get assigned variation
	assignmentKey := "assignment-" + subjectKey + "-" + flagKey
	shard := getShard(assignmentKey, config.SubjectShards)
	variations := allocation.Variations
	var variationShard Variation

	for _, variation := range variations {
		if isShardInRange(shard, variation.ShardRange) {
			variationShard = variation
		}
	}

	assignedVariation := variationShard.Value

	func() {
		// need to catch panics from Logger and continue
		defer func() {
			r := recover()
			if r != nil {
				fmt.Println("panic occurred:", r)
			}
		}()

		// Log assignment
		assignmentEvent := AssignmentEvent{
			Experiment:        flagKey + "-" + rule.AllocationKey,
			FeatureFlag:       flagKey,
			Allocation:        rule.AllocationKey,
			Variation:         assignedVariation,
			Subject:           subjectKey,
			Timestamp:         TimeNow(),
			SubjectAttributes: subjectAttributes,
		}
		ec.logger.LogAssignment(assignmentEvent)
	}()

	return assignedVariation, nil
}

func getSubjectVariationOverride(experimentConfig experimentConfiguration, subject string, valueType ValueType) Value {
	hash := md5.Sum([]byte(subject))
	hashOutput := hex.EncodeToString(hash[:])

	if val, ok := experimentConfig.Overrides[hashOutput]; ok {
		return val
	}

	return Null()
}

func isInExperimentSample(subjectKey string, flagKey string, subjectShards int64, percentExposure float32) bool {
	shardKey := "exposure-" + subjectKey + "-" + flagKey
	shard := getShard(shardKey, subjectShards)

	return float64(shard) <= float64(percentExposure)*float64(subjectShards)
}
