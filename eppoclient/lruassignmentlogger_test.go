package eppoclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_LruAssignmentLogger_cacheAssignment(t *testing.T) {
	innerLogger := new(mockLogger)
	innerLogger.On("LogAssignment", mock.Anything).Return()

	logger, err := NewLruAssignmentLogger(innerLogger, 1000)
	assert.NoError(t, err)

	event := AssignmentEvent{
		Experiment:        "testExperiment",
		FeatureFlag:       "testFeatureFlag",
		Allocation:        "testAllocation",
		Variation:         "123.45",
		Subject:           "testSubject",
		Timestamp:         "testTimestamp",
		SubjectAttributes: Attributes{"testKey": "testValue"},
	}

	logger.LogAssignment(event)
	logger.LogAssignment(event)

	innerLogger.AssertNumberOfCalls(t, "LogAssignment", 1)
}

func Test_LruAssignmentLogger_timestampAndAttributesAreNotImportant(t *testing.T) {
	innerLogger := new(mockLogger)
	innerLogger.On("LogAssignment", mock.Anything).Return()

	logger, err := NewLruAssignmentLogger(innerLogger, 1000)
	assert.NoError(t, err)

	logger.LogAssignment(AssignmentEvent{
		FeatureFlag:       "testFeatureFlag",
		Allocation:        "testAllocation",
		Variation:         "testVariation",
		Subject:           "testSubject",
		Experiment:        "testExperiment",
		Timestamp:         "t1",
		SubjectAttributes: Attributes{"testKey": "testValue1"},
	})
	logger.LogAssignment(AssignmentEvent{
		FeatureFlag:       "testFeatureFlag",
		Allocation:        "testAllocation",
		Variation:         "testVariation",
		Subject:           "testSubject",
		Experiment:        "testExperiment",
		Timestamp:         "t2",
		SubjectAttributes: Attributes{"testKey": "testValue2"},
	})

	innerLogger.AssertNumberOfCalls(t, "LogAssignment", 1)
}

func Test_LruAssignmentLogger_panicsAreNotCached(t *testing.T) {
	innerLogger := new(mockLogger)
	innerLogger.On("LogAssignment", mock.Anything).Panic("test panic")

	logger, err := NewLruAssignmentLogger(innerLogger, 1000)
	assert.NoError(t, err)

	event := AssignmentEvent{
		Experiment:        "testExperiment",
		FeatureFlag:       "testFeatureFlag",
		Allocation:        "testAllocation",
		Variation:         "123.45",
		Subject:           "testSubject",
		Timestamp:         "testTimestamp",
		SubjectAttributes: Attributes{"testKey": "testValue"},
	}

	assert.Panics(t, func() {
		logger.LogAssignment(event)
	})
	assert.Panics(t, func() {
		logger.LogAssignment(event)
	})

	innerLogger.AssertNumberOfCalls(t, "LogAssignment", 2)
}

func Test_LruAssignmentLogger_changeInAllocationCausesLogging(t *testing.T) {
	innerLogger := new(mockLogger)
	innerLogger.On("LogAssignment", mock.Anything).Return()

	logger, err := NewLruAssignmentLogger(innerLogger, 1000)
	assert.NoError(t, err)

	logger.LogAssignment(AssignmentEvent{
		Experiment:        "testExperiment",
		FeatureFlag:       "testFeatureFlag",
		Allocation:        "testAllocation1",
		Variation:         "variation",
		Subject:           "testSubject",
		Timestamp:         "testTimestamp",
		SubjectAttributes: Attributes{"testKey": "testValue"},
	})
	logger.LogAssignment(AssignmentEvent{
		Experiment:        "testExperiment",
		FeatureFlag:       "testFeatureFlag",
		Allocation:        "testAllocation2",
		Variation:         "variation",
		Subject:           "testSubject",
		Timestamp:         "testTimestamp",
		SubjectAttributes: Attributes{"testKey": "testValue"},
	})

	innerLogger.AssertNumberOfCalls(t, "LogAssignment", 2)
}

func Test_LruAssignmentLogger_changeInVariationCausesLogging(t *testing.T) {
	innerLogger := new(mockLogger)
	innerLogger.On("LogAssignment", mock.Anything).Return()

	logger, err := NewLruAssignmentLogger(innerLogger, 1000)
	assert.NoError(t, err)

	logger.LogAssignment(AssignmentEvent{
		Experiment:        "testExperiment",
		FeatureFlag:       "testFeatureFlag",
		Allocation:        "testAllocation",
		Variation:         "variation1",
		Subject:           "testSubject",
		Timestamp:         "testTimestamp",
		SubjectAttributes: Attributes{"testKey": "testValue"},
	})
	logger.LogAssignment(AssignmentEvent{
		Experiment:        "testExperiment",
		FeatureFlag:       "testFeatureFlag",
		Allocation:        "testAllocation",
		Variation:         "variation2",
		Subject:           "testSubject",
		Timestamp:         "testTimestamp",
		SubjectAttributes: Attributes{"testKey": "testValue"},
	})

	innerLogger.AssertNumberOfCalls(t, "LogAssignment", 2)
}

func Test_LruAssignmentLogger_allocationOscillationLogsAll(t *testing.T) {
	innerLogger := new(mockLogger)
	innerLogger.On("LogAssignment", mock.Anything).Return()

	logger, err := NewLruAssignmentLogger(innerLogger, 1000)
	assert.NoError(t, err)

	logger.LogAssignment(AssignmentEvent{
		Experiment:        "testExperiment",
		FeatureFlag:       "testFeatureFlag",
		Allocation:        "testAllocation1",
		Variation:         "variation",
		Subject:           "testSubject",
		Timestamp:         "t1",
		SubjectAttributes: Attributes{"testKey": "testValue"},
	})
	logger.LogAssignment(AssignmentEvent{
		Experiment:        "testExperiment",
		FeatureFlag:       "testFeatureFlag",
		Allocation:        "testAllocation2",
		Variation:         "variation",
		Subject:           "testSubject",
		Timestamp:         "t2",
		SubjectAttributes: Attributes{"testKey": "testValue"},
	})
	logger.LogAssignment(AssignmentEvent{
		Experiment:        "testExperiment",
		FeatureFlag:       "testFeatureFlag",
		Allocation:        "testAllocation1",
		Variation:         "variation",
		Subject:           "testSubject",
		Timestamp:         "t3",
		SubjectAttributes: Attributes{"testKey": "testValue"},
	})
	logger.LogAssignment(AssignmentEvent{
		Experiment:        "testExperiment",
		FeatureFlag:       "testFeatureFlag",
		Allocation:        "testAllocation2",
		Variation:         "variation",
		Subject:           "testSubject",
		Timestamp:         "t4",
		SubjectAttributes: Attributes{"testKey": "testValue"},
	})

	innerLogger.AssertNumberOfCalls(t, "LogAssignment", 4)
}

func Test_LruAssignmentLogger_variationOscillationLogsAll(t *testing.T) {
	innerLogger := new(mockLogger)
	innerLogger.On("LogAssignment", mock.Anything).Return()

	logger, err := NewLruAssignmentLogger(innerLogger, 1000)
	assert.NoError(t, err)

	logger.LogAssignment(AssignmentEvent{
		Experiment:        "testExperiment",
		FeatureFlag:       "testFeatureFlag",
		Allocation:        "testAllocation",
		Variation:         "variation1",
		Subject:           "testSubject",
		Timestamp:         "t1",
		SubjectAttributes: Attributes{"testKey": "testValue"},
	})
	logger.LogAssignment(AssignmentEvent{
		Experiment:        "testExperiment",
		FeatureFlag:       "testFeatureFlag",
		Allocation:        "testAllocation",
		Variation:         "variation2",
		Subject:           "testSubject",
		Timestamp:         "t2",
		SubjectAttributes: Attributes{"testKey": "testValue"},
	})
	logger.LogAssignment(AssignmentEvent{
		Experiment:        "testExperiment",
		FeatureFlag:       "testFeatureFlag",
		Allocation:        "testAllocation",
		Variation:         "variation1",
		Subject:           "testSubject",
		Timestamp:         "t3",
		SubjectAttributes: Attributes{"testKey": "testValue"},
	})
	logger.LogAssignment(AssignmentEvent{
		Experiment:        "testExperiment",
		FeatureFlag:       "testFeatureFlag",
		Allocation:        "testAllocation",
		Variation:         "variation2",
		Subject:           "testSubject",
		Timestamp:         "t4",
		SubjectAttributes: Attributes{"testKey": "testValue"},
	})

	innerLogger.AssertNumberOfCalls(t, "LogAssignment", 4)
}

func Test_LruAssignmentLogger_proxyLogBanditAction(t *testing.T) {
	innerLogger := new(mockLogger)
	innerLogger.On("LogAssignment", mock.Anything).Return()
	innerLogger.On("LogBanditAction", mock.Anything).Return()

	logger, err := NewLruAssignmentLogger(innerLogger, 1000)
	assert.NoError(t, err)

	event := BanditEvent{
		FlagKey:                      "flag",
		BanditKey:                    "bandit",
		Subject:                      "subject",
		Action:                       "action",
		ActionProbability:            0.1,
		OptimalityGap:                0.1,
		ModelVersion:                 "model-version",
		Timestamp:                    "timestamp",
		SubjectNumericAttributes:     map[string]float64{},
		SubjectCategoricalAttributes: map[string]string{},
		ActionNumericAttributes:      map[string]float64{},
		ActionCategoricalAttributes:  map[string]string{},
		MetaData:                     map[string]string{},
	}

	banditLogger := logger.(BanditActionLogger)

	banditLogger.LogBanditAction(event)
	banditLogger.LogBanditAction(event)

	innerLogger.AssertNumberOfCalls(t, "LogBanditAction", 2)
}
