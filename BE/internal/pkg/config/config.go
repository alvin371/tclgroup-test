package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config holds all application configuration.
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Logger   LoggerConfig
}

// AppConfig holds application-level settings.
type AppConfig struct {
	Env     string `env:"APP_ENV"     env-default:"development"`
	Address string `env:"APP_ADDRESS" env-default:":8080"`
}

// DatabaseConfig holds database connection settings.
type DatabaseConfig struct {
	DSN      string `env:"DATABASE_DSN"       env-required:"true"`
	MaxConns int32  `env:"DATABASE_MAX_CONNS" env-default:"10"`
	MinConns int32  `env:"DATABASE_MIN_CONNS" env-default:"2"`
}

// LoggerConfig holds logging settings.
type LoggerConfig struct {
	Level string `env:"LOG_LEVEL" env-default:"info"`
	Env   string `env:"LOG_ENV"   env-default:"production"`
}

// MustLoad loads config from environment variables and panics on error.
func MustLoad() *Config {
	cfg := &Config{}
	_ = cleanenv.ReadConfig(".env", cfg)
	if err := cleanenv.ReadEnv(cfg); err != nil {
		panic(fmt.Sprintf("config: %v", err))
	}
	return cfg
}

// Load loads config from environment variables and returns an error on failure.
func Load() (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("read env: %w", err)
	}
	return cfg, nil
}
