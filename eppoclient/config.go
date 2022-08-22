package eppoclient

type Config struct {
	baseUrl          string `default:"https://eppo.cloud/api"`
	apiKey           string
	assignmentLogger AssignmentLogger
}

func (cfg *Config) validate() {
	if cfg.apiKey == "" {
		panic("api key not set")
	}
}
