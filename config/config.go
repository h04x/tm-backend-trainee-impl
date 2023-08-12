package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"path"
)

type Config struct {
	App `yaml:"app"`
	Pg  `yaml:"postgres"`
}

type App struct {
	Listen string `env-required:"true" yaml:"listen" env:"APP_LISTEN"`
}

type Pg struct {
	Url string `env-required:"true" yaml:"URL" env:"PG_URL"`
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
