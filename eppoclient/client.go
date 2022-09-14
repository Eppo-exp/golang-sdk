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

	experimentConfig, err := ec.configRequestor.GetConfiguration(flagKey)
	if err != nil {
		return "", err
	}

	override := getSubjectVariationOverride(experimentConfig, subjectKey)

	if override != "" {
		return override, nil
	}

	if !experimentConfig.Enabled ||
		!subjectAttributesSatisfyRules(subjectAttributes, experimentConfig.Rules) ||
		!isInExperimentSample(subjectKey, flagKey, experimentConfig) {
		return "", errors.New("not in sample")
	}

	assignmentKey := "assignment-" + subjectKey + "-" + flagKey
	shard := getShard(assignmentKey, int64(experimentConfig.SubjectShards))
	variations := experimentConfig.Variations
	var variationShard Variation

	for _, variation := range variations {
		if isShardInRange(int(shard), variation.ShardRange) {
			variationShard = variation
		}
	}

	assignedVariation := variationShard.Name

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

func subjectAttributesSatisfyRules(subjectAttributes dictionary, rules []rule) bool {
	if len(rules) == 0 {
		return true
	}

	return matchesAnyRule(subjectAttributes, rules)
}

func isInExperimentSample(subjectKey string, flagKey string, experimentConfig experimentConfiguration) bool {
	shardKey := "exposure-" + subjectKey + "-" + flagKey
	shard := getShard(shardKey, int64(experimentConfig.SubjectShards))

	return float64(shard) <= float64(experimentConfig.PercentExposure)*float64(experimentConfig.SubjectShards)
}
