package eppoclient

const default_base_url = "https://eppo.cloud/api"

type Config struct {
	BaseUrl          string
	ApiKey           string
	AssignmentLogger AssignmentLogger
}

func (cfg *Config) validate() {
	if cfg.ApiKey == "" {
		panic("api key not set")
	}

	if cfg.BaseUrl == "" {
		cfg.BaseUrl = default_base_url
	}
}
