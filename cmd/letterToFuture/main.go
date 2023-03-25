package main

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"

	"LetterToFuture/internal/apiserver"
	"LetterToFuture/internal/model"
	"LetterToFuture/internal/telegram"
)

const MAX_MESSAGE_LIMIT = 4095

func main() {
	server := apiserver.NewAPIServer()

	if err := server.Start(); err != nil {
		server.Logger.Fatal(err)
	}
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	token := os.Getenv("TELEGRAM_TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	commands := telegram.NewCommands()

	model2 := model.TemporaryModel()

	for update := range updates {
		if update.Message.Text != "" { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			switch update.Message.Command() {
			case "start":
				commands.CommandMode(update.Message.Command())
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет, "+update.Message.From.UserName+"! Это бот \"Письма в будущее\". "+
					"Чтобы ознакомиться поподробнее с тем, что я умею, вызовете комманду /help.")
				bot.Send(msg)
			case "help":
				commands.CommandMode(update.Message.Command())
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Я отправляю письма в будущее. Для того чтобы начать вызовете команду /goletter")
				bot.Send(msg)
			case "goletter":
				commands.CommandMode(update.Message.Command())
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Отправьте текст своего письма:")
				bot.Send(msg)
			default:
				if commands.Start == true {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Нажмите /help")
					bot.Send(msg)
				}
				if commands.Help == true {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Нажмите /goletter")
					bot.Send(msg)
				}
				if commands.Goletter == true {
					textMsg := update.Message.Text
					if len(textMsg) > MAX_MESSAGE_LIMIT {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Извини, я не могу сохранить такое большое письмо. Сократите его и попробуйте еще раз!")
						bot.Send(msg)
						continue
					} else {
						if model2.Letter == "" {
							model2.Letter = update.Message.Text
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Теперь введите вашу почту:")

							bot.Send(msg)
						} else if model2.Email == "" {
							model2.Email = update.Message.Text
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Теперь введите дату:")
							bot.Send(msg)
						} else if model2.Date == "" {
							model2.Date = update.Message.Text
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Спасибо!")
							bot.Send(msg)
							fmt.Println(model2)
						}
						//server.Store.CreateALetter(model.NewModel("dsfafdasfasdfasd", "2222-11-11", textMsg))
					}

				}
			}
		} else {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "К сожалению, я принимаю только текстовый формат."+
				" Нажмите /help"))
		}
	}

}
