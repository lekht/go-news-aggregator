package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App    `yaml:"app"`
		Server `yaml:"server"`
		PG     `yaml:"postgres"`
		RSS
	}

	App struct {
		Name    string `env-required:"true" yaml:"name" env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	Server struct {
		Port string `env-required:"true" yaml:"port" env:"SERVER_PORT" envDefault:":80"`
	}

	PG struct {
		User     string `env-required:"true" yaml:"user" env:"PG_USER"`
		Password string `env-required:"true" yaml:"password" env:"PG_PASSWORD"`
		Host     string `env-required:"true" yaml:"host" env:"PG_HOST"`
		Port     string `env-required:"true" yaml:"port" env:"PG_PORT"`
		DB       string `env-required:"true" yaml:"db" env:"PG_DB"`
	}

	RSS struct {
		URLs   []string `json:"rss"`
		Period int64    `json:"request_period"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	b, err := os.ReadFile("./config/config.json")
	if err != nil {
		return nil, fmt.Errorf("config - json - error: %w", err)
	}
	var config RSS
	err = json.Unmarshal(b, &config)
	if err != nil {
		return nil, fmt.Errorf("config -json unmarshal - error: %w", err)
	}

	cfg.RSS = config

	return cfg, nil
}
