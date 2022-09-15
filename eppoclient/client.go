package eppoclient

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Client for eppo.cloud. Instance of this struct will be created on calling InitClient.
// EppoClient will then immediately start polling experiments data from eppo.cloud.
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

func (ec *EppoClient) GetAssignment(subjectKey string, flagKey string, subjectAttributes dictionary) (string, error) {
	if subjectKey == "" {
		panic("no subject key provided")
	}

	if flagKey == "" {
		panic("no flag key provided")
	}

	config, err := ec.configRequestor.GetConfiguration(flagKey)
	if err != nil {
		return "", err
	}

	override := getSubjectVariationOverride(config, subjectKey)

	if override != "" {
		return override, nil
	}

	// Check if disabled
	if !config.Enabled {
		return "", errors.New("the experiment or flag is not enabled")
	}

	// Find matching rule
	rule, err := findMatchingRule(subjectAttributes, config.Rules)
	if err != nil {
		return "", err
	}

	// Check if in sample population
	allocation := config.Allocations[rule.AllocationKey]
	if !isInExperimentSample(subjectKey, flagKey, config.SubjectShards, allocation.PercentExposure) {
		return "", errors.New("subject not part of the sample population")
	}

	// Get assigned variation
	assignmentKey := "assignment-" + subjectKey + "-" + flagKey
	shard := getShard(assignmentKey, int64(config.SubjectShards))
	variations := allocation.Variations
	var variationShard Variation

	for _, variation := range variations {
		if isShardInRange(int(shard), variation.ShardRange) {
			variationShard = variation
		}
	}

	assignedVariation := variationShard.Value.(string)

	assignmentEvent := AssignmentEvent{
		Experiment:        flagKey,
		Variation:         assignedVariation,
		Subject:           subjectKey,
		Timestamp:         time.Now().String(),
		SubjectAttributes: subjectAttributes,
	}

	_, jsonErr := json.Marshal(assignmentEvent)

	if jsonErr != nil {
		panic("incorrect json")
	}

	func() {
		// need to catch panics from Logger and continue
		defer func() {
			r := recover()
			if r != nil {
				fmt.Println("panic occurred:", r)
			}
		}()

		ec.logger.LogAssignment(assignmentEvent)
	}()

	return assignedVariation, nil
}

func getSubjectVariationOverride(experimentConfig experimentConfiguration, subject string) string {
	hash := md5.Sum([]byte(subject))
	hashOutput := hex.EncodeToString(hash[:])

	if val, ok := experimentConfig.Overrides[hashOutput]; ok {
		return val.(string)
	}

	return ""
}

func isInExperimentSample(subjectKey string, flagKey string, subjectShards int, percentExposure float32) bool {
	shardKey := "exposure-" + subjectKey + "-" + flagKey
	shard := getShard(shardKey, int64(subjectShards))

	return float64(shard) <= float64(percentExposure)*float64(subjectShards)
}
