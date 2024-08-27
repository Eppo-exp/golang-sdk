package eppoclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_LruBanditLogger_cacheBanditAction(t *testing.T) {
	innerLogger := new(mockLogger)
	innerLogger.On("LogAssignment", mock.Anything).Return()
	innerLogger.On("LogBanditAction", mock.Anything).Return()

	logger, err := NewLruBanditLogger(innerLogger, 1000)
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

	innerLogger.AssertNumberOfCalls(t, "LogBanditAction", 1)
}

func Test_LruBanditLogger_okIfInnerLoggerIsNotBandit(t *testing.T) {
	innerLogger := new(mockNonBanditLogger)
	innerLogger.On("LogAssignment", mock.Anything).Return()

	logger, err := NewLruBanditLogger(innerLogger, 1000)
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

	innerLogger.AssertNumberOfCalls(t, "LogBanditAction", 0)
}
