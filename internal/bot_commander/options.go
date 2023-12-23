package bot_commander

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type Option func(bot *BotCommander) error

func WithLogger() Option {
	return func(bot *BotCommander) error {
		l := logrus.New()
		l.SetFormatter(
			&logrus.TextFormatter{DisableColors: false},
		)
		l.SetLevel(logrus.InfoLevel)

		bot.logger = l

		return nil
	}
}

func WithTgAPI() Option {
	return func(bot *BotCommander) error {
		b, err := tgbotapi.NewBotAPI("5982458978:AAHPIJjXWs4-Nu3JTnmlxtjQ8Yya90kZnNk")
		if err != nil {
			log.Panic(err)
		}

		b.Debug = true

		bot.tg = b

		return nil
	}
}
