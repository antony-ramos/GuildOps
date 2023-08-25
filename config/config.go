package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App     `yaml:"app"`
		Discord `yaml:"discord"`
		Log     `yaml:"logger"`
		Backend `yaml:"backend"`
	}

	// App -.
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	// Discord -.
	Discord struct {
		Token   string `env-required:"true" yaml:"token" env:"DISCORD_TOKEN"`
		GuildID string `env-required:"true" yaml:"guild_id" env:"DISCORD_GUILD_ID"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}

	// Backend -.
	Backend struct {
		URL string `env-required:"true"   yaml:"url"              env:"COVEN_BACKEND_URL"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
