package config

import (
	"errors"
	"fmt"
	"mythic-plus-crawler/assets"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	BlizzardAPI struct {
		ClientID       string `yaml:"client-id" env:"CLIENT_ID"`
		ClientSecret   string `yaml:"client-secret" env:"CLIENT_SECRET"`
		RequestTimeout int    `yaml:"request-timeout" env:"REQUEST_TIMEOUT"`
	} `yaml:"blizzard-api" env-prefix:"BLIZZARD_"`

	Database struct {
		Host     string `yaml:"host" env:"HOST"`
		Port     uint16 `yaml:"port" env:"PORT"`
		User     string `yaml:"username" env:"USER"`
		Password string `yaml:"password" env:"PASSWORD"`
		Database string `yaml:"database" env:"DATABASE"`
	} `yaml:"database" env-prefix:"DB_"`
}

func LoadConfig() (*Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	configPath := filepath.Join(wd, "config.yml")

	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(configPath)
		if err != nil {
			return nil, fmt.Errorf("error while creating default config: %w", err)
		}

		_, err = file.WriteString(assets.DefaultConfig)

		if err != nil {
			return nil, fmt.Errorf("error while writing default conifg: %w", err)
		}
	}

	var cfg Config
	err = cleanenv.ReadConfig("config.yml", &cfg)
	if err != nil {
		return nil, fmt.Errorf("error while loading config: %w", err)
	}

	return &cfg, nil
}
