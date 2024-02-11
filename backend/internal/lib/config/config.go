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
	DB         `yaml:"db"`
	Paths      `yaml:"paths"`
}

type HTTPServer struct {
	Host        string        `env:"host" env-default:"localhost"`
	Port        int           `yaml:"port" env-default:"8080"`
	Timeout     time.Duration `env:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `env:"idle_timeout" env-default:"30s"`
}

type DB struct {
	Name     string `yaml:"name" env-required:"true"`
	User     string `yaml:"user" env-default:"user"`
	Password string `yaml:"password" env-default:"password"`
	Host     string `yaml:"host" env-default:"postgres"`
	Port     string `yaml:"port" env-default:"5432"`
}

type Paths struct {
	SignUp  string `yaml:"signup" env-required:"false"`
	SignIn  string `yaml:"signin" env-required:"false"`
	Refresh string `yaml:"refresh" env-required:"false"`
	Edit    string `yaml:"edit" env-required:"false"`
	Info    string `yaml:"info" env-required:"false"`
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
