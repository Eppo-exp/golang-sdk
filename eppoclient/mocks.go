package eppoclient

import (
	"github.com/stretchr/testify/mock"
)

type MockLogger struct {
	mock.Mock
}

func (ml *MockLogger) LogAssignment(event string) {
	ml.Called(event)
}

type MockConfigRequestor struct {
	mock.Mock
}

func (mcr *MockConfigRequestor) GetConfiguration(experimentKey string) (ExperimentConfiguration, error) {
	args := mcr.Called(experimentKey)

	return args.Get(0).(ExperimentConfiguration), args.Error(1)
}

func (mcr *MockConfigRequestor) FetchAndStoreConfigurations() {
}
