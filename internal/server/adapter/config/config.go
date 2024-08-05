package config

import "github.com/caarlos0/env/v11"

type LogLevel string

const (
	LogLevelDebug = LogLevel("DEBUG")
	LogLevelInfo  = LogLevel("INFO")
)

type Config struct {
	LogLevel        LogLevel `env:"LOG_LEVEL" envDefault:"DEBUG"`
	TokenExpiration int      `env:"TOKEN_EXPIRATION" envDefault:"24"`
}

func New() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
