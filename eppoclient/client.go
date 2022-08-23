package eppoclient

import (
	"encoding/json"
	"time"
)

type EppoClient struct {
	configRequestor IConfigRequestor
	poller          Poller
	logger          AssignmentLogger
}

type AssignmentEvent struct {
	Experiment        string
	Variation         string
	Subject           string
	Timestamp         string
	SubjectAttributes Dictionary
}

func NewEppoClient(configRequestor IConfigRequestor, assignmentLogger AssignmentLogger) *EppoClient {
	var ec = &EppoClient{}

	var poller = NewPoller(10, configRequestor.FetchAndStoreConfigurations)
	ec.poller = *poller
	ec.configRequestor = configRequestor
	ec.logger = assignmentLogger

	return ec
}

func (ec *EppoClient) GetAssignment(subjectKey string, experimentKey string, subjectAttributes Dictionary) (string, error) {
	if subjectKey == "" {
		panic("no subject key provided")
	}

	if experimentKey == "" {
		panic("no experiment key provided")
	}

	experimentConfig, err := ec.configRequestor.GetConfiguration(experimentKey)

	if err != nil {
		return "", err
	}

	if !experimentConfig.Enabled ||
		!subjectAttributesSatisfyRules(subjectAttributes, experimentConfig.Rules) ||
		!isInExperimentSample(subjectKey, experimentKey, experimentConfig) {
		return "", err
	}

	assignmentKey := "assignment-" + subjectKey + "-" + experimentKey
	shard := getShard(assignmentKey, int64(experimentConfig.SubjectShards))
	variations := experimentConfig.Variations
	var variationShard Variation

	for _, variation := range variations {
		if isShardInRange(int(shard), variation.ShardRange) {
			variationShard = variation
		}
	}

	assignedVariation := variationShard.Name

	assignmentEvent := &AssignmentEvent{
		Experiment:        experimentKey,
		Variation:         assignedVariation,
		Subject:           subjectKey,
		Timestamp:         time.Now().String(),
		SubjectAttributes: subjectAttributes,
	}

	aeJson, _ := json.Marshal(assignmentEvent)

	ec.logger.LogAssignment(string(aeJson))

	return assignedVariation, err
}

func subjectAttributesSatisfyRules(subjectAttributes Dictionary, rules []Rule) bool {
	if len(rules) == 0 {
		return true
	}

	return matchesAnyRule(subjectAttributes, rules)
}

func isInExperimentSample(subjectKey string, experimentKey string, experimentConfig ExperimentConfiguration) bool {
	shardKey := "exposure-" + subjectKey + "-" + experimentKey
	shard := getShard(shardKey, int64(experimentConfig.SubjectShards))

	return shard <= int64(experimentConfig.PercentExposure)*int64(experimentConfig.SubjectShards)
}
