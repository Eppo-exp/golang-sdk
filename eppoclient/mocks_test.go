package eppoclient

import (
	"github.com/stretchr/testify/mock"
)

type mockLogger struct {
	mock.Mock
}

func (ml *mockLogger) LogAssignment(event AssignmentEvent) {
	ml.MethodCalled("LogAssignment", event)
}

func (ml *mockLogger) LogBanditAction(event BanditEvent) {
	ml.MethodCalled("LogBanditAction", event)
}

// `mockNonBanditLogger` is missing `LogBanditAction` and therefore
// does not implement `BanditActionLogger`.
type mockNonBanditLogger struct {
	mock.Mock
}

func (ml *mockNonBanditLogger) LogAssignment(event AssignmentEvent) {
	ml.MethodCalled("LogAssignment", event)
}
