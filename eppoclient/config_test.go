package eppoclient

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_config_defaultPollerInterval(t *testing.T) {
	cfg := Config{
		SdkKey: "blah",
	}

	err := cfg.validate()
	assert.NoError(t, err)
	assert.Equal(t, 10*time.Second, cfg.PollerInterval)
}
