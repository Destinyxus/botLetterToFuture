package botclient

import (
	"github.com/Destinyxus/botLetterToFuture/internal/bot_commander"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type botClient struct {
	botAPI *tgbotapi.BotAPI
}

func New(botAPI *tgbotapi.BotAPI) bot_commander.BotClient {
	botAPI.Debug = true
	return &botClient{botAPI: botAPI}
}

func (bc *botClient) GetUpdatesChan(email, letter string) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	bc.botAPI.GetUpdatesChan(u)

}
