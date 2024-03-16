package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	BlenderPath    string         `yaml:"blender_path" env-required:"true"`
	BaseURL        string         `yaml:"base_url" env-required:"true"`
	SleepTime      time.Duration  `yaml:"sleep_time" env-default:"5s"`
	UpdateStatus   UpdateStatus   `yaml:"update_status"`
	ResponseStatus ResponseStatus `yaml:"response_status"`
}

type UpdateStatus struct {
	InProgress string `yaml:"in_progress" env-default:"in-progress"`
	Error      string `yaml:"error" env-default:"error"`
}

type ResponseStatus struct {
	Empty string `yaml:"empty" env-default:"empty"`
}

func MustLoad(configPath string) *Config {
	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("error opening config file: %s", err)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("error reading conig: %s", err)
	}

	return &cfg
}
