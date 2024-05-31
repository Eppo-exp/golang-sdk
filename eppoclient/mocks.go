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

func (mcr *mockConfigRequestor) GetConfiguration(flagKey string) (flagConfiguration, error) {
	args := mcr.Called(flagKey)

	return args.Get(0).(flagConfiguration), args.Error(1)
}

func (mcr *mockConfigRequestor) FetchAndStoreConfigurations() {
}
