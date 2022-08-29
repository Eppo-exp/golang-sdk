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
	configRequestor IConfigRequestor
	poller          Poller
	logger          IAssignmentLogger
}

type AssignmentEvent struct {
	Experiment        string
	Variation         string
	Subject           string
	Timestamp         string
	SubjectAttributes Dictionary
}

func NewEppoClient(configRequestor IConfigRequestor, assignmentLogger IAssignmentLogger) *EppoClient {
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

	override := getSubjectVariationOverride(experimentConfig, subjectKey)

	if override != "" {
		return override, nil
	}

	if !experimentConfig.Enabled ||
		!subjectAttributesSatisfyRules(subjectAttributes, experimentConfig.Rules) ||
		!isInExperimentSample(subjectKey, experimentKey, experimentConfig) {
		return "", errors.New("not in sample")
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

	func() {
		// need to catch panics from Logger and continue
		defer func() {
			r := recover()
			if r != nil {
				fmt.Println("panic occurred:", r)
			}
		}()

		event := map[string]string{}
		json.Unmarshal([]byte(aeJson), &event)

		ec.logger.LogAssignment(event)
	}()

	return assignedVariation, nil
}

func getSubjectVariationOverride(experimentConfig ExperimentConfiguration, subject string) string {
	hash := md5.Sum([]byte(subject))
	hashOutput := hex.EncodeToString(hash[:])

	if val, ok := experimentConfig.Overrides[hashOutput]; ok {
		return val.(string)
	}

	return ""
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

	return float64(shard) <= float64(experimentConfig.PercentExposure)*float64(experimentConfig.SubjectShards)
}
