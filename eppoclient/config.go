package eppoclient

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

const defaultBaseUrl = "https://fscdn.eppo.cloud/api"
const defaultPollerInterval = 10 * time.Second

type Config struct {
	BaseUrl           string
	SdkKey            string
	AssignmentLogger  IAssignmentLogger
	PollerInterval    time.Duration
	ApplicationLogger ApplicationLogger
	HttpClient        *http.Client
}

func (cfg *Config) validate() error {
	if cfg.SdkKey == "" {
		return fmt.Errorf("SDK key not set")
	}

	if cfg.BaseUrl == "" {
		cfg.BaseUrl = defaultBaseUrl
	}

	if cfg.PollerInterval <= 0 {
		cfg.PollerInterval = defaultPollerInterval
	}

	if cfg.ApplicationLogger == nil {
		defaultLogger, err := zap.NewProduction(zap.IncreaseLevel(zap.WarnLevel))
		if err != nil {
			return fmt.Errorf("failed to create default logger: %v", err)
		}
		cfg.ApplicationLogger = NewZapLogger(defaultLogger)
	}

	return nil
}
