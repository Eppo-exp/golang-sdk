package eppoclient

import (
	"context"
	"fmt"
	"time"
)

type Attributes map[string]interface{}

// EppoClient Client for eppo.cloud. Instance of this struct will be created on calling InitClient.
// EppoClient will then immediately start polling experiments data from Eppo.
type EppoClient struct {
	configurationStore *configurationStore
	configRequestor    *configurationRequestor
	poller             *poller
	logger             IAssignmentLogger
	loggerContext      IAssignmentLoggerContext
	applicationLogger  ApplicationLogger
}

func newEppoClient(
	configurationStore *configurationStore,
	configRequestor *configurationRequestor,
	poller *poller,
	assignmentLogger IAssignmentLogger,
	assignmentLoggerContext IAssignmentLoggerContext,
	applicationLogger ApplicationLogger,
) *EppoClient {
	return &EppoClient{
		configurationStore: configurationStore,
		configRequestor:    configRequestor,
		poller:             poller,
		logger:             assignmentLogger,
		loggerContext:      assignmentLoggerContext,
		applicationLogger:  applicationLogger,
	}
}

// Returns a channel that gets closed after client has been
// *successfully* initialized.
//
// It is recommended to apply a timeout to initialization as otherwise
// it may hang up indefinitely.
//
//	select {
//	case <-client.Initialized():
//	case <-time.After(5 * time.Second):
//	}
func (ec *EppoClient) Initialized() <-chan struct{} {
	return ec.configurationStore.Initialized()
}

func (ec *EppoClient) GetBoolAssignment(
	flagKey, subjectKey string,
	subjectAttributes Attributes,
	defaultValue bool,
) (bool, error) {
	return ec.getBoolAssignment(context.Background(), flagKey, subjectKey, subjectAttributes, defaultValue)
}

func (ec *EppoClient) GetBoolAssignmentContext(
	ctx context.Context,
	flagKey, subjectKey string,
	subjectAttributes Attributes,
	defaultValue bool,
) (bool, error) {
	return ec.getBoolAssignment(ctx, flagKey, subjectKey, subjectAttributes, defaultValue)
}

func (ec *EppoClient) getBoolAssignment(
	ctx context.Context,
	flagKey, subjectKey string,
	subjectAttributes Attributes,
	defaultValue bool,
) (bool, error) {
	variation, err := ec.getAssignment(ctx, ec.configurationStore.getConfiguration(), flagKey, subjectKey, subjectAttributes, booleanVariation)
	if err != nil || variation == nil {
		return defaultValue, err
	}
	result, ok := variation.(bool)
	if !ok {
		ec.applicationLogger.Errorf("failed to cast %v to bool", variation)
		return defaultValue, fmt.Errorf("failed to cast %v to bool", variation)
	}
	return result, err
}

func (ec *EppoClient) GetNumericAssignment(
	flagKey, subjectKey string,
	subjectAttributes Attributes,
	defaultValue float64,
) (float64, error) {
	return ec.getNumericAssignment(context.Background(), flagKey, subjectKey, subjectAttributes, defaultValue)
}

func (ec *EppoClient) GetNumericAssignmentContext(
	ctx context.Context,
	flagKey, subjectKey string,
	subjectAttributes Attributes,
	defaultValue float64,
) (float64, error) {
	return ec.getNumericAssignment(ctx, flagKey, subjectKey, subjectAttributes, defaultValue)
}

func (ec *EppoClient) getNumericAssignment(
	ctx context.Context,
	flagKey, subjectKey string,
	subjectAttributes Attributes,
	defaultValue float64,
) (float64, error) {
	variation, err := ec.getAssignment(ctx, ec.configurationStore.getConfiguration(), flagKey, subjectKey, subjectAttributes, numericVariation)
	if err != nil || variation == nil {
		return defaultValue, err
	}
	result, ok := variation.(float64)
	if !ok {
		ec.applicationLogger.Errorf("failed to cast %v to float64", variation)
		return defaultValue, fmt.Errorf("failed to cast %v to float64", variation)
	}
	return result, err
}

func (ec *EppoClient) GetIntegerAssignment(
	flagKey, subjectKey string,
	subjectAttributes Attributes,
	defaultValue int64,
) (int64, error) {
	return ec.getIntegerAssignment(context.Background(), flagKey, subjectKey, subjectAttributes, defaultValue)
}

func (ec *EppoClient) GetIntegerAssignmentContext(
	ctx context.Context,
	flagKey, subjectKey string,
	subjectAttributes Attributes,
	defaultValue int64,
) (int64, error) {
	return ec.getIntegerAssignment(context.Background(), flagKey, subjectKey, subjectAttributes, defaultValue)
}

func (ec *EppoClient) getIntegerAssignment(
	ctx context.Context,
	flagKey, subjectKey string,
	subjectAttributes Attributes,
	defaultValue int64,
) (int64, error) {
	variation, err := ec.getAssignment(ctx, ec.configurationStore.getConfiguration(), flagKey, subjectKey, subjectAttributes, integerVariation)
	if err != nil || variation == nil {
		return defaultValue, err
	}
	result, ok := variation.(int64)
	if !ok {
		ec.applicationLogger.Errorf("failed to cast %v to int64", variation)
		return defaultValue, fmt.Errorf("failed to cast %v to int64", variation)
	}
	return result, err
}

func (ec *EppoClient) GetStringAssignment(
	flagKey, subjectKey string,
	subjectAttributes Attributes,
	defaultValue string,
) (string, error) {
	return ec.getStringAssignment(context.Background(), flagKey, subjectKey, subjectAttributes, defaultValue)
}

func (ec *EppoClient) GetStringAssignmentContext(
	ctx context.Context,
	flagKey, subjectKey string,
	subjectAttributes Attributes,
	defaultValue string,
) (string, error) {
	return ec.getStringAssignment(ctx, flagKey, subjectKey, subjectAttributes, defaultValue)
}

func (ec *EppoClient) getStringAssignment(
	ctx context.Context,
	flagKey, subjectKey string,
	subjectAttributes Attributes,
	defaultValue string,
) (string, error) {
	variation, err := ec.getAssignment(ctx, ec.configurationStore.getConfiguration(), flagKey, subjectKey, subjectAttributes, stringVariation)
	if err != nil || variation == nil {
		return defaultValue, err
	}
	result, ok := variation.(string)
	if !ok {
		ec.applicationLogger.Errorf("failed to cast %v to string", variation)
		return defaultValue, fmt.Errorf("failed to cast %v to string", variation)
	}
	return result, err
}

func (ec *EppoClient) GetJSONAssignment(
	flagKey, subjectKey string,
	subjectAttributes Attributes,
	defaultValue any,
) (any, error) {
	return ec.getJSONAssignment(context.Background(), flagKey, subjectKey, subjectAttributes, defaultValue)
}

func (ec *EppoClient) GetJSONAssignmentContext(
	ctx context.Context,
	flagKey, subjectKey string,
	subjectAttributes Attributes,
	defaultValue any,
) (any, error) {
	return ec.getJSONAssignment(ctx, flagKey, subjectKey, subjectAttributes, defaultValue)
}

func (ec *EppoClient) getJSONAssignment(
	ctx context.Context,
	flagKey, subjectKey string,
	subjectAttributes Attributes,
	defaultValue any,
) (any, error) {
	variation, err := ec.getAssignment(ctx, ec.configurationStore.getConfiguration(), flagKey, subjectKey, subjectAttributes, jsonVariation)
	if err != nil || variation == nil {
		return defaultValue, err
	}
	result, ok := variation.(jsonVariationValue)
	if !ok {
		ec.applicationLogger.Errorf("failed to cast %v to json. This should never happen. Please report bug to Eppo", variation)
		return defaultValue, fmt.Errorf("failed to cast %v to json. This should never happen. Please report bug to Eppo", variation)
	}
	return result.Parsed, err
}

func (ec *EppoClient) GetJSONBytesAssignment(
	flagKey, subjectKey string,
	subjectAttributes Attributes,
	defaultValue []byte,
) ([]byte, error) {
	return ec.getJSONBytesAssignment(context.Background(), flagKey, subjectKey, subjectAttributes, defaultValue)
}

func (ec *EppoClient) GetJSONBytesAssignmentContext(
	ctx context.Context,
	flagKey, subjectKey string,
	subjectAttributes Attributes,
	defaultValue []byte,
) ([]byte, error) {
	return ec.getJSONBytesAssignment(ctx, flagKey, subjectKey, subjectAttributes, defaultValue)
}

func (ec *EppoClient) getJSONBytesAssignment(
	ctx context.Context,
	flagKey, subjectKey string,
	subjectAttributes Attributes,
	defaultValue []byte,
) ([]byte, error) {
	variation, err := ec.getAssignment(ctx, ec.configurationStore.getConfiguration(), flagKey, subjectKey, subjectAttributes, jsonVariation)
	if err != nil || variation == nil {
		return defaultValue, err
	}
	result, ok := variation.(jsonVariationValue)
	if !ok {
		ec.applicationLogger.Errorf("failed to cast %v to json. This should never happen. Please report bug to Eppo", variation)
		return defaultValue, fmt.Errorf("failed to cast %v to json. This should never happen. Please report bug to Eppo", variation)
	}
	return result.Raw, err
}

type BanditResult struct {
	Variation string
	Action    *string
}

func (ec *EppoClient) GetBanditAction(
	flagKey, subjectKey string,
	subjectAttributes ContextAttributes,
	actions map[string]ContextAttributes,
	defaultVariation string,
) BanditResult {
	return ec.getBanditAction(context.Background(), flagKey, subjectKey, subjectAttributes, actions, defaultVariation)
}

func (ec *EppoClient) GetBanditActionContext(
	ctx context.Context,
	flagKey, subjectKey string,
	subjectAttributes ContextAttributes,
	actions map[string]ContextAttributes,
	defaultVariation string,
) BanditResult {
	return ec.getBanditAction(ctx, flagKey, subjectKey, subjectAttributes, actions, defaultVariation)
}

func (ec *EppoClient) getBanditAction(
	ctx context.Context,
	flagKey, subjectKey string,
	subjectAttributes ContextAttributes,
	actions map[string]ContextAttributes,
	defaultVariation string,
) BanditResult {
	config := ec.configurationStore.getConfiguration()

	// ignoring the error here as we can always proceed with default variation
	assignmentValue, _ := ec.getAssignment(ctx, config, flagKey, subjectKey, subjectAttributes.toGenericAttributes(), stringVariation)
	variation, ok := assignmentValue.(string)
	if !ok {
		variation = defaultVariation
	}

	// If no actions have been passed, we will return the variation, even if it is a bandit key
	if len(actions) == 0 {
		return BanditResult{
			Variation: variation,
			Action:    nil,
		}
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

	ec.logBanditAction(ctx, BanditEvent{
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
	})

	return BanditResult{
		Variation: variation,
		Action:    &evaluation.actionKey,
	}
}

func (ec *EppoClient) getAssignment(ctx context.Context, config configuration, flagKey string, subjectKey string, subjectAttributes Attributes, variationType variationType) (interface{}, error) {
	if subjectKey == "" {
		return nil, fmt.Errorf("no subject key provided")
	}

	if flagKey == "" {
		return nil, fmt.Errorf("no flag key provided")
	}

	flag, err := config.getFlagConfiguration(flagKey)
	if err != nil {
		ec.applicationLogger.Infof("failed to get flag configuration: %v", err)
		return nil, err
	}

	err = flag.verifyType(variationType)
	if err != nil {
		ec.applicationLogger.Warnf("failed to verify flag type: %v", err)
		return nil, err
	}

	assignmentValue, assignmentEvent, err := flag.eval(subjectKey, subjectAttributes, ec.applicationLogger)
	if err != nil {
		ec.applicationLogger.Errorf("failed to evaluate flag: %v", err)
		return nil, err
	}

	ec.logAssignment(ctx, assignmentEvent)
	return assignmentValue, nil
}

func (ec *EppoClient) logAssignment(ctx context.Context, event *AssignmentEvent) {
	if event == nil {
		return
	}

	// need to catch panics from Logger and continue
	defer func() {
		r := recover()
		if r != nil {
			ec.applicationLogger.Errorf("panic occurred: %v", r)
		}
	}()

	switch {
	case ec.loggerContext != nil:
		// prioritise logger context if it's defined in the config
		ec.loggerContext.LogAssignment(ctx, *event)
	case ec.logger != nil:
		ec.logger.LogAssignment(*event)
	default:
		// no need to do anything if both logger are nil
	}
}

func (ec *EppoClient) logBanditAction(ctx context.Context, event BanditEvent) {
	// need to catch panics from Logger and continue
	defer func() {
		r := recover()
		if r != nil {
			fmt.Println("panic occurred:", r)
		}
	}()

	switch {
	case ec.loggerContext != nil:
		if l, ok := ec.loggerContext.(interface {
			LogBanditAction(context.Context, BanditEvent)
		}); ok {
			l.LogBanditAction(ctx, event)
		}
	case ec.logger != nil:
		if l, ok := ec.logger.(BanditActionLogger); ok {
			l.LogBanditAction(event)
		}
	}
}
