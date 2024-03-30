package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		PG      `yaml:"postgres"`
		Redis   `yaml:"redis"`
		Server  `yaml:"server"`
		General `yaml:"general"`
	}
	PG struct {
		URL          string `yaml:"url" env:"PG_URL" env-required:"true"`
		ConnPoolSize int    `yaml:"maxConnPoolSize" env:"PG_MAX_POOL_SIZE"`
	}
	Redis struct {
		URL    string `yaml:"url" env:"REDIS_URL" env-required:"true"`
		PidTTL int64  `yaml:"pidTTL"`
	}
	Server struct {
		Port    string `yaml:"port" env:"HTTP_SERVER_PORT"`
		LogPath string `yaml:"logPath"`
	}
	General struct {
		DefaultShellPath string `yaml:"defaultShellPath"`
	}
)

func New(filePath string) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(filePath, cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	//	err = cleanenv.UpdateEnv(cfg)
	//	if err != nil {
	//		return nil, fmt.Errorf("error updating env: %w", err)
	//	}

	return cfg, nil
}
