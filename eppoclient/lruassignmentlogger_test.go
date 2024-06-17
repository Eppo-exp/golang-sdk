package eppoclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_LruAssignmentLogger_cacheAssignment(t *testing.T) {
	innerLogger := new(mockLogger)
	innerLogger.On("LogAssignment", mock.Anything).Return()

	logger := NewLruAssignmentLogger(innerLogger, 1000)

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

	logger := NewLruAssignmentLogger(innerLogger, 1000)

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

	logger := NewLruAssignmentLogger(innerLogger, 1000)

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

	logger := NewLruAssignmentLogger(innerLogger, 1000)

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

	logger := NewLruAssignmentLogger(innerLogger, 1000)

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

	logger := NewLruAssignmentLogger(innerLogger, 1000)

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

	logger := NewLruAssignmentLogger(innerLogger, 1000)

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
