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
	Auth struct {
		GoogleAuthClientID     string `env-required:"true" yaml:"google_auth_client_id" env:"GOOGLE_AUTH_CLIENT_ID"`
		GoogleAuthClientSecret string `env-required:"true" yaml:"google_auth_client_secret" env:"GOOGLE_AUTH_CLIENT_SECRET"`
		GoogleAuthRedirectURL  string `env-required:"true" yaml:"google_auth_redirect_url" env:"GOOGLE_AUTH_REDIRECT_URL"`
		JwtSecret              string `env-required:"true" yaml:"jwt-secret" env:"JWT_SECRET"`
		TokenExpiredIn         int    `env-required:"true" yaml:"token-expired-in" env:"TOKEN_EXPIRED_IN"`
		JwtRefreshSecret       string `env-required:"true" yaml:"jwt_refresh_secret" env:"JWT_REFRESH_SECRET"`
		JwtRefreshKeyExpire    int    `env-required:"true" yaml:"jwt_refresh_key_expire_hours_count" env:"JWT_REFRESH_KEY_EXPIRE_HOURS_COUNT"`
		TokenMaxAge            int    `env-required:"true" yaml:"token-maxage" env:"TOKEN_MAXAGE"`
		RefreshTokenMaxAge     int    `env-required:"true" yaml:"refresh-token-maxage" env:"REFRESH_TOKEN_MAXAGE"`
		FrontendOrigin         string `env-required:"true" yaml:"frontend_origin" env:"FRONTEND_ORIGIN"`
	} `yaml:"auth"`
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
