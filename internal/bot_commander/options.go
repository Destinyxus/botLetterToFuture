package bot_commander

import (
	"log"

	"github.com/Destinyxus/botLetterToFuture/internal/emailSender"
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

func WithTgAPI(token string) Option {
	return func(bot *BotCommander) error {
		b, err := tgbotapi.NewBotAPI(token)
		if err != nil {
			log.Panic(err)
		}

		b.Debug = true

		bot.tg = b

		return nil
	}
}

func WithEmailSender(token string, mailName string, address string) Option {
	return func(bot *BotCommander) error {
		em := emailSender.NewEmail(token, mailName, address)

		bot.emailSender = em

		return nil
	}
}
