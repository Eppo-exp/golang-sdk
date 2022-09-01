package eppoclient

import (
	"github.com/stretchr/testify/mock"
)

type mockLogger struct {
	mock.Mock
}

func (ml *mockLogger) LogAssignment(event AssignmentEvent) {
	ml.Called(event)
}

type mockConfigRequestor struct {
	mock.Mock
}

func (mcr *mockConfigRequestor) GetConfiguration(experimentKey string) (experimentConfiguration, error) {
	args := mcr.Called(experimentKey)

	return args.Get(0).(experimentConfiguration), args.Error(1)
}

func (mcr *mockConfigRequestor) FetchAndStoreConfigurations() {
}
