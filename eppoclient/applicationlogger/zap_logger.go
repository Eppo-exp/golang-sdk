package applicationlogger

import (
	"go.uber.org/zap"
)

/**
 * The default logger for the SDK.
 */
type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger(logger *zap.Logger) *ZapLogger {
	return &ZapLogger{logger: logger}
}

func (z *ZapLogger) Debug(args ...interface{}) {
	z.logger.Sugar().Debug(args...)
}

func (z *ZapLogger) Info(args ...interface{}) {
	z.logger.Sugar().Info(args...)
}

func (z *ZapLogger) Warn(args ...interface{}) {
	z.logger.Sugar().Warn(args...)
}

func (z *ZapLogger) Error(args ...interface{}) {
	z.logger.Sugar().Error(args...)
}
