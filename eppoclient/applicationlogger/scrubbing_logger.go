package applicationlogger

import (
	"regexp"
)

type ScrubbingLogger struct {
	innerLogger Logger
}

func NewScrubbingLogger(innerLogger Logger) *ScrubbingLogger {
	return &ScrubbingLogger{innerLogger: innerLogger}
}

func (s *ScrubbingLogger) scrub(args ...interface{}) []interface{} {
	scrubbedArgs := make([]interface{}, len(args))
	for i, arg := range args {
		strArg, ok := arg.(string)
		if ok {
			strArg = maskSensitiveInfo(strArg)
			scrubbedArgs[i] = strArg
		} else {
			scrubbedArgs[i] = arg
		}
	}
	return scrubbedArgs
}

// maskSensitiveInfo replaces sensitive information (like apiKey or sdkKey)
// in the error message with 'XXXXXX' to prevent exposure of these keys in
// logs or error messages.
func maskSensitiveInfo(errMsg string) string {
	// Scrub apiKey and sdkKey from error messages containing URLs
	// Matches any string that starts with apiKey or sdkKey followed by any characters until the next & or the end of the string
	re := regexp.MustCompile(`(apiKey|sdkKey)=[^&]*`)
	return re.ReplaceAllString(errMsg, "$1=XXXXXX")
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
