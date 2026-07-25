package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_DefaultsAndOverrides(t *testing.T) {
	// 1. Clear environment
	os.Unsetenv("PORT")
	os.Unsetenv("REQUEST_TIMEOUT")

	// 2. Load
	cfg, err := LoadConfig(".")
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// Since we might have a local .env file in the path, let's verify overrides or defaults
	// We'll set env variables directly and ensure they override
	os.Setenv("PORT", "9999")
	os.Setenv("REQUEST_TIMEOUT", "30s")
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("REQUEST_TIMEOUT")
	}()

	cfgOverride, err := LoadConfig(".")
	assert.NoError(t, err)
	assert.Equal(t, 9999, cfgOverride.Port)
	assert.Equal(t, 30*time.Second, cfgOverride.RequestTimeout)
}
