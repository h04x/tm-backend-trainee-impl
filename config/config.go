package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"path"
)

type Config struct {
	App `yaml:"app"`
	Pg
}

type App struct {
	Port string `env-required:"true" yaml:"port" env:"APP_PORT"`
}

type Pg struct {
	Url string `env-required:"true" env:"POSTGRES_URL"`
}

func NewConfig(configPath string) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(path.Join("./", configPath), cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	err = cleanenv.UpdateEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("error updating env: %w", err)
	}

	return cfg, nil
}
