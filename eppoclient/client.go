package eppoclient

import (
	"fmt"
	"time"

	"github.com/Eppo-exp/golang-sdk/v5/eppoclient/applicationlogger"
)

type Attributes map[string]interface{}

// Client for eppo.cloud. Instance of this struct will be created on calling InitClient.
// EppoClient will then immediately start polling experiments data from Eppo.
type EppoClient struct {
	configurationStore *configurationStore
	configRequestor    *configurationRequestor
	poller             *poller
	logger             IAssignmentLogger
	applicationLogger  applicationlogger.Logger
}

func newEppoClient(
	configurationStore *configurationStore,
	configRequestor *configurationRequestor,
	poller *poller,
	assignmentLogger IAssignmentLogger,
	applicationLogger applicationlogger.Logger,
) *EppoClient {
	return &EppoClient{
		configurationStore: configurationStore,
		configRequestor:    configRequestor,
		poller:             poller,
		logger:             assignmentLogger,
		applicationLogger:  applicationLogger,
	}
}

func (ec *EppoClient) GetBoolAssignment(flagKey string, subjectKey string, subjectAttributes Attributes, defaultValue bool) (bool, error) {
	variation, err := ec.getAssignment(ec.configurationStore.getConfiguration(), flagKey, subjectKey, subjectAttributes, booleanVariation)
	if err != nil || variation == nil {
		return defaultValue, err
	}
	result, ok := variation.(bool)
	if !ok {
		ec.applicationLogger.Error("failed to cast %v to bool", variation)
		return defaultValue, fmt.Errorf("failed to cast %v to bool", variation)
	}
	return result, err
}

func (ec *EppoClient) GetNumericAssignment(flagKey string, subjectKey string, subjectAttributes Attributes, defaultValue float64) (float64, error) {
	variation, err := ec.getAssignment(ec.configurationStore.getConfiguration(), flagKey, subjectKey, subjectAttributes, numericVariation)
	if err != nil || variation == nil {
		return defaultValue, err
	}
	result, ok := variation.(float64)
	if !ok {
		ec.applicationLogger.Error("failed to cast %v to float64", variation)
		return defaultValue, fmt.Errorf("failed to cast %v to float64", variation)
	}
	return result, err
}

func (ec *EppoClient) GetIntegerAssignment(flagKey string, subjectKey string, subjectAttributes Attributes, defaultValue int64) (int64, error) {
	variation, err := ec.getAssignment(ec.configurationStore.getConfiguration(), flagKey, subjectKey, subjectAttributes, integerVariation)
	if err != nil || variation == nil {
		return defaultValue, err
	}
	result, ok := variation.(int64)
	if !ok {
		ec.applicationLogger.Error("failed to cast %v to int64", variation)
		return defaultValue, fmt.Errorf("failed to cast %v to int64", variation)
	}
	return result, err
}

func (ec *EppoClient) GetStringAssignment(flagKey string, subjectKey string, subjectAttributes Attributes, defaultValue string) (string, error) {
	variation, err := ec.getAssignment(ec.configurationStore.getConfiguration(), flagKey, subjectKey, subjectAttributes, stringVariation)
	if err != nil || variation == nil {
		return defaultValue, err
	}
	result, ok := variation.(string)
	if !ok {
		ec.applicationLogger.Error("failed to cast %v to string", variation)
		return defaultValue, fmt.Errorf("failed to cast %v to string", variation)
	}
	return result, err
}

func (ec *EppoClient) GetJSONAssignment(flagKey string, subjectKey string, subjectAttributes Attributes, defaultValue interface{}) (interface{}, error) {
	variation, err := ec.getAssignment(ec.configurationStore.getConfiguration(), flagKey, subjectKey, subjectAttributes, jsonVariation)
	if err != nil || variation == nil {
		return defaultValue, err
	}
	return variation, err
}

type BanditResult struct {
	Variation string
	Action    *string
}

func (ec *EppoClient) GetBanditAction(flagKey string, subjectKey string, subjectAttributes ContextAttributes, actions map[string]ContextAttributes, defaultVariation string) BanditResult {
	config := ec.configurationStore.getConfiguration()

	isBanditFlag := config.isBanditFlag(flagKey)

	if isBanditFlag && len(actions) == 0 {
		// No actions passed for a flag known to have an
		// active bandit, so we just return the default values
		// so that we don't log a variation or bandit
		// assignment.
		return BanditResult{
			Variation: defaultVariation,
			Action:    nil,
		}
	}

	// ignoring the error here as we can always proceed with default variation
	assignmentValue, _ := ec.getAssignment(config, flagKey, subjectKey, subjectAttributes.toGenericAttributes(), stringVariation)
	variation, ok := assignmentValue.(string)
	if !ok {
		variation = defaultVariation
	}

	banditVariation, ok := config.getBanditVariant(flagKey, variation)
	if !ok {
		return BanditResult{
			Variation: variation,
			Action:    nil,
		}
	}

	bandit, err := config.getBanditConfiguration(banditVariation.Key)
	if err != nil {
		// no bandit configuration
		return BanditResult{
			Variation: variation,
			Action:    nil,
		}
	}

	evaluation := bandit.ModelData.evaluate(banditEvaluationContext{
		flagKey:           flagKey,
		subjectKey:        subjectKey,
		subjectAttributes: subjectAttributes,
		actions:           actions,
	})

	if logger, ok := ec.logger.(BanditActionLogger); ok {
		event := BanditEvent{
			FlagKey:                      flagKey,
			BanditKey:                    bandit.BanditKey,
			Subject:                      subjectKey,
			Action:                       evaluation.actionKey,
			ActionProbability:            evaluation.actionWeight,
			OptimalityGap:                evaluation.optimalityGap,
			ModelVersion:                 bandit.ModelVersion,
			Timestamp:                    time.Now().UTC().Format(time.RFC3339),
			SubjectNumericAttributes:     evaluation.subjectAttributes.Numeric,
			SubjectCategoricalAttributes: evaluation.subjectAttributes.Categorical,
			ActionNumericAttributes:      evaluation.actionAttributes.Numeric,
			ActionCategoricalAttributes:  evaluation.actionAttributes.Categorical,
			MetaData: map[string]string{
				"sdkLanguage": "go",
				"sdkVersion":  __version__,
			},
		}

		func() {
			// need to catch panics from Logger and continue
			defer func() {
				r := recover()
				if r != nil {
					fmt.Println("panic occurred:", r)
				}
			}()

			logger.LogBanditAction(event)
		}()
	}

	return BanditResult{
		Variation: variation,
		Action:    &evaluation.actionKey,
	}
}

func (ec *EppoClient) getAssignment(config configuration, flagKey string, subjectKey string, subjectAttributes Attributes, variationType variationType) (interface{}, error) {
	if subjectKey == "" {
		return nil, fmt.Errorf("no subject key provided")
	}

	if flagKey == "" {
		return nil, fmt.Errorf("no flag key provided")
	}

	flag, err := config.getFlagConfiguration(flagKey)
	if err != nil {
		ec.applicationLogger.Info("failed to get flag configuration: %v", err)
		return nil, err
	}

	err = flag.verifyType(variationType)
	if err != nil {
		ec.applicationLogger.Info("failed to verify flag type: %v", err)
		return nil, err
	}

	assignmentValue, assignmentEvent, err := flag.eval(subjectKey, subjectAttributes, ec.applicationLogger)
	if err != nil {
		ec.applicationLogger.Info("failed to evaluate flag: %v", err)
		return nil, err
	}

	if assignmentEvent != nil {
		func() {
			// need to catch panics from Logger and continue
			defer func() {
				r := recover()
				if r != nil {
					ec.applicationLogger.Error("panic occurred: %v", r)
				}
			}()

			// Log assignment
			ec.logger.LogAssignment(*assignmentEvent)
		}()
	}

	return assignmentValue, nil
}
