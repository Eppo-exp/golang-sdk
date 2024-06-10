package eppoclient

import "time"

const default_base_url = "https://fscdn.eppo.cloud/api"

type Config struct {
	BaseUrl          string
	SdkKey           string
	AssignmentLogger IAssignmentLogger
	PollerInterval   time.Duration
}

func (cfg *Config) validate() {
	if cfg.SdkKey == "" {
		panic("api key not set")
	}

	if cfg.BaseUrl == "" {
		cfg.BaseUrl = default_base_url
	}

	if cfg.PollerInterval <= 0 {
		cfg.PollerInterval = 10 * time.Second
	}
}
