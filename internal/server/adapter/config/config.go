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
	DatabaseDSN     string   `env:"DATABASE_DSN" envDefault:"host=localhost port=5432 user=myuser password=123 dbname=gophkeeper sslmode=disable"`
}

func New() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
