package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `env:"ENV" env-default:"local"`
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Host        string        `env:"host" env-default:"localhost"`
	Port        int           `yaml:"port" env-default:"8080"`
	Timeout     time.Duration `env:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `env:"idle_timeout" env-default:"30s"`
}

func MustLoad(configPath string) *Config {
	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("error opening config file: %s", err)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("error reading envs: %s", err)
	}

	return &cfg
}
