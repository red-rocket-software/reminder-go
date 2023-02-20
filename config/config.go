package config

import (
	"fmt"
	"log"
	"time"

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
	Auth struct {
		GoogleAuthClientId     string        `env-required:"true" yaml:"google_auth_client_id" env:"GOOGLE_AUTH_CLIENT_ID"`
		GoogleAuthClientSecret string        `env-required:"true" yaml:"google_auth_client_secret" env:"GOOGLE_AUTH_CLIENT_SECRET"`
		GoogleAuthRedirectUrl  string        `env-required:"true" yaml:"google_auth_redirect_url" env:"GOOGLE_AUTH_REDIRECT_URL"`
		JwtSecret              string        `env-required:"true" yaml:"jwt-secret" env:"JWT_SECRET"`
		TokenExpiredIn         time.Duration `env-required:"true" yaml:"token-expired-in" env:"TOKEN_EXPIRED_IN"`
		TokenMaxAge            int           `env-required:"true" yaml:"token-maxage" env:"TOKEN_MAXAGE"`
		FrontendOrigin         string        `env-required:"true" yaml:"frontend_origin" env:"FRONTEND_ORIGIN"`


		GithubAuthClientID     string `env-required:"true" yaml:"github_auth_client_id" env:"GITHUB_AUTH_CLIENT_ID"`
		GithubAuthClientSecret string `env-required:"true" yaml:"github_auth_client_secret" env:"GITHUB_AUTH_CLIENT_SECRET"`
		GithubAuthRedirectURL  string `env-required:"true" yaml:"github_auth_redirect_url" env:"GITHUB_AUTH_REDIRECT_URL"`
		LinkedinAuthClientID     string `env-require: "true" yaml: "linkedin_auth_client_id" env:"LINKEDIN_AUTH_CLIENT_ID"`
		LinkedinAuthClientSecret string `env-require: "true" yaml: "linkedin_auth_client_secret" env:"LINKEDIN_AUTH_CLIENT_SECRET"`
		LinkedinAuthRedirectURL  string `env-required:"true" yaml:"linkedin_auth_redirect_url" env:"LINKEDIN_AUTH_REDIRECT_URL"`

	} `yaml:"auth"`
}

func GetConfig() *Config {

	log.Print("config init")

	c := &Config{}

	if err := cleanenv.ReadConfig("config.yaml", c); err != nil {
		fmt.Println("error read config")
	}

	return c
}
