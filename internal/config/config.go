package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel string `toml:"log_level" env:"LOG_LEVEL" default:"info"`
	Server   Server
	Auth     Auth
}

func Load(path string, cfg *Config) error {
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	return nil
}
