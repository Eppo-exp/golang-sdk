package eppoclient

import (
	"fmt"
	"regexp"

	"go.uber.org/zap"
)

// sensitiveInfoRe matches apiKey or sdkKey query parameters in URLs.
// Compiled once at package init for performance.
var sensitiveInfoRe = regexp.MustCompile(`(apiKey|sdkKey)=[^&\s"]*`)

type ApplicationLogger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
}

// ZapLogger The default logger for the Eppo SDK
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

func (z *ZapLogger) Infof(template string, args ...interface{}) {
	z.logger.Sugar().Infof(template, args...)
}

func (z *ZapLogger) Warn(args ...interface{}) {
	z.logger.Sugar().Warn(args...)
}

func (z *ZapLogger) Warnf(template string, args ...interface{}) {
	z.logger.Sugar().Warnf(template, args...)
}

func (z *ZapLogger) Error(args ...interface{}) {
	z.logger.Sugar().Error(args...)
}

func (z *ZapLogger) Errorf(template string, args ...interface{}) {
	z.logger.Sugar().Errorf(template, args...)
}

// ScrubbingLogger is an ApplicationLogger that scrubs sensitive information from logs
type ScrubbingLogger struct {
	innerLogger ApplicationLogger
}

func NewScrubbingLogger(innerLogger ApplicationLogger) *ScrubbingLogger {
	return &ScrubbingLogger{innerLogger: innerLogger}
}

func (s *ScrubbingLogger) scrub(args ...interface{}) []interface{} {
	scrubbedArgs := make([]interface{}, len(args))
	for i, arg := range args {
		switch v := arg.(type) {
		case string:
			scrubbedArgs[i] = maskSensitiveInfo(v)
		case error:
			scrubbedArgs[i] = fmt.Errorf("%s", maskSensitiveInfo(v.Error()))
		default:
			scrubbedArgs[i] = arg
		}
	}
	return scrubbedArgs
}

// maskSensitiveInfo replaces sensitive information (like apiKey or sdkKey)
// in the error message with 'XXXXXX' to prevent exposure of these keys in
// logs or error messages.
func maskSensitiveInfo(errMsg string) string {
	return sensitiveInfoRe.ReplaceAllString(errMsg, "$1=XXXXXX")
}

func (s *ScrubbingLogger) Debug(args ...interface{}) {
	s.innerLogger.Debug(s.scrub(args...)...)
}

func (s *ScrubbingLogger) Info(args ...interface{}) {
	s.innerLogger.Info(s.scrub(args...)...)
}

func (s *ScrubbingLogger) Infof(template string, args ...interface{}) {
	s.innerLogger.Infof(maskSensitiveInfo(template), s.scrub(args...)...)
}

func (s *ScrubbingLogger) Warn(args ...interface{}) {
	s.innerLogger.Warn(s.scrub(args...)...)
}

func (s *ScrubbingLogger) Warnf(template string, args ...interface{}) {
	s.innerLogger.Warnf(maskSensitiveInfo(template), s.scrub(args...)...)
}

func (s *ScrubbingLogger) Error(args ...interface{}) {
	s.innerLogger.Error(s.scrub(args...)...)
}

func (s *ScrubbingLogger) Errorf(template string, args ...interface{}) {
	s.innerLogger.Errorf(maskSensitiveInfo(template), s.scrub(args...)...)
}
