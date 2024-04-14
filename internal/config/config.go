package config

import (
	"github.com/kelseyhightower/envconfig"
	"log/slog"
	"sync"
)

type Config struct {
	Env      string `envconfig:"ENV"`
	LogLevel string `envconfig:"LOG_LEVEL"`
	MongoURI string `envconfig:"MONGODB_URI"`
	HTTPAddr string `envconfig:"HTTP_ADDR"`
	ApiKey   string `envconfig:"API_KEY"`
}

var (
	config Config
	once   sync.Once
)

// Get reads config from environment
func Get(logger *slog.Logger) *Config {
	once.Do(func() {
		err := envconfig.Process("", &config)
		if err != nil {
			logger.Error("failed to read config", "error", err)
		}

		logger.Info("config initialized successfully", "config", config)
	})
	return &config
}

func (c *Config) GetSlogLevel() slog.Level {
	switch c.LogLevel {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}
