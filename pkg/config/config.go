package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	TelegramToken string `env:"TELEGRAM_TOKEN" env-required:"TELEGRAM_TOKEN"`
	SendGridKey   string `env:"SENDGRID_KEY" env-required:"SENDGRID_KEY"`
	StoreURL      string `env:"STORE_URL" env-required:"STORE_URL"`
	HashKey       string `env:"HASH_KEY" env-required:"HASH_KEY"`
	Errors        Errors
	Responses     Responses
}

type Errors struct {
	StartTrue    string `toml:"StartTrue"`
	HelpTrue     string `toml:"HelpTrue"`
	SizeLetter   string `toml:"SizeLetter"`
	InvalidEmail string `toml:"InvalidEmail"`
	InvalidDate  string `toml:"InvalidDate"`
}

type Responses struct {
	Start    string `toml:"Start"`
	HelpText string `toml:"HelpText"`
	Goletter string `toml:"Goletter"`
	Email    string `toml:"Email"`
	Date     string `toml:"Date"`
	Result   string `toml:"Result"`
}

func New(path string) (*Config, error) {
	cfg := new(Config)

	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		return nil, fmt.Errorf("reading config from %s: %w", path, err)
	}

	return cfg, nil
}
