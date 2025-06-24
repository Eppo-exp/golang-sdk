package eppoclient

import (
	"context"

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

type mockLoggerContext struct {
	mock.Mock
}

func (ml *mockLoggerContext) LogAssignment(ctx context.Context, event AssignmentEvent) {
	ml.MethodCalled("LogAssignment", ctx, event)
}

func (ml *mockLoggerContext) LogBanditAction(ctx context.Context, event BanditEvent) {
	ml.MethodCalled("LogBanditAction", ctx, event)
}

// `mockNonBanditLogger` is missing `LogBanditAction` and therefore
// does not implement `BanditActionLogger`.
type mockNonBanditLogger struct {
	mock.Mock
}

func (ml *mockNonBanditLogger) LogAssignment(event AssignmentEvent) {
	ml.MethodCalled("LogAssignment", event)
}
