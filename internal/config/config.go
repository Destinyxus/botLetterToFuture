package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

// env-required:"SENDGRID_KEY"
// env-required:"TELEGRAM_TOKEN
// env-required:"SENDGRID_ADDRESS
type Config struct {
	TelegramToken   string `env:"TELEGRAM_TOKEN"`
	SendGridKey     string `env:"SENDGRID_KEY"`
	SendGridAddress string `env:"SENDGRID_ADDRESS"`
	LetterName      string `toml:"LetterName"`
	Errors          Errors
	Responses       Responses
}

type Errors struct {
	SizeLetter           string `toml:"SizeLetter"`
	InvalidFormatMessage string `toml:"InvalidFormatMessage"`
}

type Responses struct {
	AboutDescription string `toml:"AboutDescription"`
	Result           string `toml:"Result"`
}

func New(path string) (*Config, error) {
	cfg := new(Config)

	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		return nil, fmt.Errorf("reading config from %s: %w", path, err)
	}

	return cfg, nil
}
