package eppoclient

import (
	"time"

	"github.com/Eppo-exp/golang-sdk/v4/eppoclient/applicationlogger"
	"go.uber.org/zap"
)

const default_base_url = "https://fscdn.eppo.cloud/api"

type Config struct {
	BaseUrl           string
	SdkKey            string
	AssignmentLogger  IAssignmentLogger
	PollerInterval    time.Duration
	ApplicationLogger applicationlogger.Logger
}

func (cfg *Config) validate() {
	if cfg.SdkKey == "" {
		panic("sdk key not set")
	}

	if cfg.BaseUrl == "" {
		cfg.BaseUrl = default_base_url
	}

	if cfg.PollerInterval <= 0 {
		cfg.PollerInterval = 10 * time.Second
	}

	if cfg.ApplicationLogger == nil {
		defaultLogger, _ := zap.NewProduction()
		cfg.ApplicationLogger = applicationlogger.NewZapLogger(defaultLogger)
	}
}
