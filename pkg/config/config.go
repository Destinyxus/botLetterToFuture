package config

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/viper"
)

type Config struct {
	TelegramToken string
	SendGridKey   string
	StoreURL      string
	HashKey       string
	Messages      Messages
}

type Messages struct {
	Responses
	Errors
}

type Errors struct {
	StartTrue    string `mapstructure:"start_true"`
	HelpTrue     string `mapstructure:"help_true"`
	SizeLetter   string `mapstructure:"size_letter"`
	InvalidEmail string `mapstructure:"invalid_email"`
	InvalidDate  string `mapstructure:"invalid_date"`
}

type Responses struct {
	Start    string `mapstructure:"start"`
	HelpText string `mapstructure:"help_text"`
	Goletter string `mapstructure:"goletter"`
	Email    string `mapstructure:"email"`
	Date     string `mapstructure:"date"`
	Result   string `mapstructure:"result"`
}

func Init() (*Config, error) {

	viper.AddConfigPath("configs")
	viper.SetConfigName("main")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("messages.responses", &cfg.Messages.Responses); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("messages.errors", &cfg.Messages.Errors); err != nil {
		return nil, err
	}

	if err := parseEnv(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func parseEnv(cfg *Config) error {

	if err := viper.BindEnv("telegram_token"); err != nil {
		return err
	}
	cfg.TelegramToken = viper.GetString("telegram_token")

	if err := viper.BindEnv("sendgrid_api_key"); err != nil {
		return err
	}
	cfg.SendGridKey = viper.GetString("sendgrid_api_key")

	if err := viper.BindEnv("store_url"); err != nil {
		return err
	}
	cfg.StoreURL = viper.GetString("store_url")

	if err := viper.BindEnv("hash_key"); err != nil {
		return err
	}
	cfg.HashKey = viper.GetString("hash_key")

	return nil
}
