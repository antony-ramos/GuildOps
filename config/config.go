package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App     `yaml:"app"`
		Discord `yaml:"discord"`
		Log     `yaml:"logger"`
		Metrics `yaml:"metrics"`
		PG      `yaml:"postgres"`
	}

	// App -.
	App struct {
		Name    string `env:"APP_NAME"    env-required:"true" yaml:"name"`
		Version string `env:"APP_VERSION" env-required:"true" yaml:"version"`
		Env     string `env:"APP_ENV"     env-required:"true" yaml:"environment"`
	}

	// Discord -.
	Discord struct {
		Token          string `env:"DISCORD_TOKEN"    env-required:"true" yaml:"token"`
		GuildID        int    `env:"DISCORD_GUILD_ID" env-required:"true" yaml:"guild_id"`
		DeleteCommands bool   `env:"DISCORD_DELETE_COMMANDS" env-required:"true" yaml:"delete_commands"`
	}

	// Log -.
	Log struct {
		Level string `env:"LOG_LEVEL" env-required:"true" yaml:"level"`
	}

	// Metrics -.
	Metrics struct {
		Port string `env:"METRICS_PORT" env-required:"true" yaml:"port"`
	}

	// PG -.
	PG struct {
		PoolMax      int           `env:"PG_POOL_MAX"      env-required:"true" yaml:"pool_max"`
		URL          string        `env:"PG_URL"           env-required:"true" yaml:"url"`
		ConnAttempts int           `env:"PG_CONN_ATTEMPTS" env-required:"true" yaml:"conn_attempts"`
		ConnTimeOut  time.Duration `env:"PG_CONN_TIMEOUT"  env-required:"true" yaml:"conn_timeout"`
	}
)

// NewConfig returns app config.
func NewConfig(configPath string) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
