package bot_commander

import (
	"github.com/Destinyxus/botLetterToFuture/internal/emailSender"
	"github.com/Destinyxus/botLetterToFuture/internal/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Option func(bot *BotCommander)

func WithLogger(l logger.Logger) Option {
	return func(bot *BotCommander) {
		bot.Logger = l
	}
}

func WithTgAPI(api *tgbotapi.BotAPI) Option {
	return func(bot *BotCommander) {
		api.Debug = true
		bot.tg = api
	}
}

func WithEmailSender(ec *emailSender.EmailClient) Option {
	return func(bot *BotCommander) {
		bot.EmailSender = ec
	}
}
