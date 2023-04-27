package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config describes all app configuration
type Config struct {
	HTTP struct {
		IP           string `env-required:"true" yaml:"ip" env:"APP_IP"`
		ReminderPort string `env-required:"true" yaml:"reminder_port" env:"REMINDER_PORT"`
		AuthPort     string `env-required:"true" yaml:"auth_port" env:"AUTH_PORT"`
	} `yaml:"http"`
	Postgres struct {
		Password string `env-default:"secret" env-required:"true" yaml:"password" env:"DB_PASSWORD"`
		Username string `env-default:"root" env-required:"true" yaml:"username" env:"DB_USERNAME"`
		Host     string `env-default:"localhost" env-required:"true" yaml:"host" env:"DB_HOST"`
		Port     string `env-default:"5432" env-required:"true" yaml:"port" env:"DB_PORT"`
		Database string `env-default:"test_remind" env-required:"true" yaml:"database" env:"DB_DATABASE"`
	} `yaml:"postgresql"`
	Email struct {
		EmailSenderName     string `env-required:"true" yaml:"email_sender_name" env:"EMAIL_SENDER_NAME"`
		EmailSenderAddress  string `env-required:"true" yaml:"email_sender_address" env:"EMAIL_SENDER_ADDRESS"`
		EmailSenderPassword string `env-required:"true" yaml:"email_sender_password" env:"EMAIL_SENDER_PASSWORD"`
	} `yaml:"email"`
}

func GetConfig() *Config {

	log.Print("config init")

	c := &Config{}

	if err := cleanenv.ReadConfig("config.yaml", c); err != nil {
		log.Fatalf("error read config: %v", err)
	}

	return c
}
