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

func maskSensitiveInfo(errMsg string) string {
	re := regexp.MustCompile(`(apiKey|sdkKey)=[^&]*`)
	return re.ReplaceAllString(errMsg, "$1=XXXXXX")
}

func (s *ScrubbingLogger) Debug(args ...interface{}) {
	s.innerLogger.Debug(s.scrub(args...)...)
}

func (s *ScrubbingLogger) Info(args ...interface{}) {
	s.innerLogger.Info(s.scrub(args...)...)
}

func (s *ScrubbingLogger) Warn(args ...interface{}) {
	s.innerLogger.Warn(s.scrub(args...)...)
}

func (s *ScrubbingLogger) Error(args ...interface{}) {
	s.innerLogger.Error(s.scrub(args...)...)
}
