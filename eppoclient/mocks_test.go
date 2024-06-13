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

func (mcr *mockConfigRequestor) FetchAndStoreConfigurations() {
}
