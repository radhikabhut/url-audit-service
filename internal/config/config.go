package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Port                int           `mapstructure:"PORT"`
	RequestTimeout      time.Duration `mapstructure:"REQUEST_TIMEOUT"`
	CacheTTL            time.Duration `mapstructure:"CACHE_TTL"`
	MaxConcurrentAudits int           `mapstructure:"MAX_CONCURRENT_AUDITS"`
	RateLimit           float64       `mapstructure:"RATE_LIMIT"`
	RateLimitWindow     time.Duration `mapstructure:"RATE_LIMIT_WINDOW"`

	// Server, Logging & Defaults
	Host      string `mapstructure:"HOST"`
	LogLevel  string `mapstructure:"LOG_LEVEL"`
	LogFormat string `mapstructure:"LOG_FORMAT"`
}

func LoadConfig(path string) (*Config, error) {
	v := viper.New()

	v.SetConfigType("env")
	v.SetConfigName(".env")
	v.AddConfigPath(path)
	v.AddConfigPath(".")

	v.AutomaticEnv()

	// Default values
	v.SetDefault("PORT", 8080)
	v.SetDefault("REQUEST_TIMEOUT", 10*time.Second)
	v.SetDefault("CACHE_TTL", 5*time.Minute)
	v.SetDefault("MAX_CONCURRENT_AUDITS", 100)
	v.SetDefault("RATE_LIMIT", 10.0)
	v.SetDefault("RATE_LIMIT_WINDOW", 1*time.Second)
	v.SetDefault("HOST", "0.0.0.0")
	v.SetDefault("LOG_LEVEL", "info")
	v.SetDefault("LOG_FORMAT", "json")

	_ = v.ReadInConfig()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
