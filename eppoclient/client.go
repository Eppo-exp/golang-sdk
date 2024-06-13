package eppoclient

import (
	"fmt"
)

type Attributes map[string]interface{}

// Client for eppo.cloud. Instance of this struct will be created on calling InitClient.
// EppoClient will then immediately start polling experiments data from Eppo.
type EppoClient struct {
	configurationStore *configurationStore
	configRequestor    *configurationRequestor
	poller             *poller
	logger             IAssignmentLogger
}

func newEppoClient(configurationStore *configurationStore, configRequestor *configurationRequestor, poller *poller, assignmentLogger IAssignmentLogger) *EppoClient {
	return &EppoClient{
		configurationStore: configurationStore,
		configRequestor:    configRequestor,
		poller:             poller,
		logger:             assignmentLogger,
	}
}

func (ec *EppoClient) GetBoolAssignment(flagKey string, subjectKey string, subjectAttributes Attributes, defaultValue bool) (bool, error) {
	variation, err := ec.getAssignment(flagKey, subjectKey, subjectAttributes, booleanVariation)
	if err != nil || variation == nil {
		return defaultValue, err
	}
	result, ok := variation.(bool)
	if !ok {
		return defaultValue, fmt.Errorf("failed to cast %v to bool", variation)
	}
	return result, err
}

func (ec *EppoClient) GetNumericAssignment(flagKey string, subjectKey string, subjectAttributes Attributes, defaultValue float64) (float64, error) {
	variation, err := ec.getAssignment(flagKey, subjectKey, subjectAttributes, numericVariation)
	if err != nil || variation == nil {
		return defaultValue, err
	}
	result, ok := variation.(float64)
	if !ok {
		return defaultValue, fmt.Errorf("failed to cast %v to float64", variation)
	}
	return result, err
}

func (ec *EppoClient) GetIntegerAssignment(flagKey string, subjectKey string, subjectAttributes Attributes, defaultValue int64) (int64, error) {
	variation, err := ec.getAssignment(flagKey, subjectKey, subjectAttributes, integerVariation)
	if err != nil || variation == nil {
		return defaultValue, err
	}
	result, ok := variation.(int64)
	if !ok {
		return defaultValue, fmt.Errorf("failed to cast %v to int64", variation)
	}
	return result, err
}

func (ec *EppoClient) GetStringAssignment(flagKey string, subjectKey string, subjectAttributes Attributes, defaultValue string) (string, error) {
	variation, err := ec.getAssignment(flagKey, subjectKey, subjectAttributes, stringVariation)
	if err != nil || variation == nil {
		return defaultValue, err
	}
	result, ok := variation.(string)
	if !ok {
		return defaultValue, fmt.Errorf("failed to cast %v to string", variation)
	}
	return result, err
}

func (ec *EppoClient) GetJSONAssignment(flagKey string, subjectKey string, subjectAttributes Attributes, defaultValue interface{}) (interface{}, error) {
	variation, err := ec.getAssignment(flagKey, subjectKey, subjectAttributes, jsonVariation)
	if err != nil || variation == nil {
		return defaultValue, err
	}
	return variation, err
}

func (ec *EppoClient) getAssignment(flagKey string, subjectKey string, subjectAttributes Attributes, variationType variationType) (interface{}, error) {
	if subjectKey == "" {
		panic("no subject key provided")
	}

	if flagKey == "" {
		panic("no flag key provided")
	}

	if ec.configRequestor != nil && !ec.configRequestor.IsAuthorized() {
		panic("Unauthorized: please check your SDK key")
	}

	flag, err := ec.configurationStore.getFlagConfiguration(flagKey)
	if err != nil {
		return nil, err
	}

	err = flag.verifyType(variationType)
	if err != nil {
		return nil, err
	}

	assignmentValue, assignmentEvent, err := flag.eval(subjectKey, subjectAttributes)
	if err != nil {
		return nil, err
	}

	if assignmentEvent != nil {
		func() {
			// need to catch panics from Logger and continue
			defer func() {
				r := recover()
				if r != nil {
					fmt.Println("panic occurred:", r)
				}
			}()

			// Log assignment
			ec.logger.LogAssignment(*assignmentEvent)
		}()
	}

	return assignmentValue, nil
}
