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

	cfg.validate()

	assert.Equal(t, 10*time.Second, cfg.PollerInterval)
}
