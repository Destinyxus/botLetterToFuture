package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	TelegramToken string `env:"TELEGRAM_TOKEN" env-required:"TELEGRAM_TOKEN"`
	Postgres      Postgres
	EmailSender   EmailSender
	Errors        Errors
	Responses     Responses
}

type Postgres struct {
	Host     string `env:"DB_HOST" env-required:"DB_HOST"`
	Port     string `env:"DB_PORT" env-required:"DB_PORT"`
	User     string `env:"DB_USER" env-required:"DB_USER"`
	Password string `env:"DB_PASSWORD" env-required:"DB_PASSWORD"`
	DBName   string `env:"DB_NAME" env-required:"DB_NAME"`
}

type EmailSender struct {
	EmailToken  string `env:"EMAIL_TOKEN" env-required:"EMAIL_TOKEN"`
	ClientEmail string `env:"CLIENT_EMAIL" env-required:"CLIENT_EMAIL"`
	HostEmail   string `env:"HOST_EMAIL" env-required:"HOST_EMAIL"`
	SMTPAddress string `env:"SMTP" env-required:"SMTP"`
}

type Errors struct {
	SizeLetter           string `toml:"SizeLetter"`
	InvalidFormatMessage string `toml:"InvalidFormatMessage"`
	NotValidCommand      string `toml:"NotValidCommand"`
}

type Responses struct {
	AboutDescription string `toml:"AboutDescription"`
	Result           string `toml:"Result"`
	StopCommand      string `toml:"StopCommand"`
	SendLetter       string `toml:"SendLetter"`
}

func New(path string) (*Config, error) {
	cfg := new(Config)

	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		return nil, fmt.Errorf("reading config from %s: %w", path, err)
	}

	return cfg, nil
}
