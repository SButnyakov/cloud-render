package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env           string `env:"ENV" env-default:"local"`
	HTTPServer    `yaml:"http_server"`
	DB            `yaml:"db"`
	Redis         `yaml:"redis"`
	Paths         `yaml:"paths"`
	External      `yaml:"external"`
	Payments      `yaml:"payments"`
	Subscriptions `yaml:"subscriptions"`
}

type HTTPServer struct {
	Host        string        `yaml:"host" env-default:"localhost"`
	Port        int           `yaml:"port" env-default:"8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"30s"`
}

type DB struct {
	Name           string `yaml:"name" env-required:"true"`
	User           string `yaml:"user" env-default:"user"`
	Password       string `yaml:"password" env-default:"password"`
	Host           string `yaml:"host" env-default:"postgres"`
	Port           string `yaml:"port" env-default:"5432"`
	MigrationsPath string `yaml:"migrations_path", env-required:"true"`
}

type Paths struct {
	SignUp    string `yaml:"signup" env-required:"false"`
	SignIn    string `yaml:"signin" env-required:"false"`
	Refresh   string `yaml:"refresh" env-required:"false"`
	Edit      string `yaml:"edit" env-required:"false"`
	Info      string `yaml:"info" env-required:"false"`
	User      string `yaml:"user" env-required:"false"`
	Subscribe string `yaml:"subscribe" env-required:"false"`
}

type Redis struct {
	Address           string `yaml:"address" env-default:"localhost:6379"`
	QueueName         string `yaml:"queue_name" env-default:"render-list"`
	PriorityQueueName string `yaml:"priority_queue_name" env-default:"render-list"`
	Password          string `yaml:"password" env-default:"password"`
}

type External struct {
	SSOUserInfo string `yanl:"sso_user_info" env-required:"false"`
}

type Payments struct {
	SubPremiumMonth string `yaml:"sub_premium_month" env-default:"sub-premium-month"`
}

type Subscriptions struct {
	Premium string `yaml:"premium" env-default:"premium"`
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
