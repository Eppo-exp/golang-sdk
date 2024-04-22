package eppoclient

import "errors"

const default_base_url = "https://fscdn.eppo.cloud/api"

type Config struct {
	BaseUrl          string
	ApiKey           string
	AssignmentLogger IAssignmentLogger
}

func (cfg *Config) validate() error {
	if cfg.ApiKey == "" {
		return errors.New("api key not set")
	}

	if cfg.BaseUrl == "" {
		cfg.BaseUrl = default_base_url
	}

	return nil
}
