package eppoclient

type Config struct {
	BaseUrl          string `default:"https://eppo.cloud/api"`
	ApiKey           string
	AssignmentLogger IAssignmentLogger
}

func (cfg *Config) validate() {
	if cfg.ApiKey == "" {
		panic("api key not set")
	}
}
