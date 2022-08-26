package eppoclient

import (
	"github.com/stretchr/testify/mock"
)

type MockLogger struct {
	mock.Mock
}

func (ml *MockLogger) LogAssignment(event map[string]string) {
	ml.Called(event)
}

type MockConfigRequestor struct {
	mock.Mock
}

func (mcr *MockConfigRequestor) GetConfiguration(experimentKey string) (experimentConfiguration, error) {
	args := mcr.Called(experimentKey)

	return args.Get(0).(experimentConfiguration), args.Error(1)
}

func (mcr *MockConfigRequestor) FetchAndStoreConfigurations() {
}
