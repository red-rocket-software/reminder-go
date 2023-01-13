package config

import (
	"fmt"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTP struct {
		IP   string `env-required:"true" yaml:"ip" env:"APP_IP"`
		Port string `env-required:"true" yaml:"port" env:"APP_PORT"`
	}
	Postgres struct {
		Password string `env-default:"secret" env-required:"true" yaml:"password" env:"DB_PASSWORD"`
		Username string `env-default:"root" env-required:"true" yaml:"username" env:"DB_USERNAME"`
		Host     string `env-default:"localhost" env-required:"true" yaml:"host" env:"DB_HOST"`
		Port     string `env-default:"5432" env-required:"true" yaml:"port" env:"DB_PORT"`
		Database string `env-default:"test_remind" env-required:"true" yaml:"database" env:"DB_DATABASE"`
	} `yaml:"postgresql"`
}

func GetConfig() *Config {

	log.Print("config init")

	c := &Config{}

	if err := cleanenv.ReadConfig("config.yaml", c); err != nil {
		fmt.Println("error read config")
	}

	return c
}
