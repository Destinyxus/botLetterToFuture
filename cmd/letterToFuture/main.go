package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"

	"LetterToFuture/internal/apiserver"
	"LetterToFuture/internal/model"
)

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

	var goletter bool = false

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if update.Message.Command() == "start" {
				goletter = false
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет, "+update.Message.From.UserName+"! Это бот \"Письма в будущее\". "+
					"Чтобы ознакомиться поподробнее с тем, что я умею, вызовете комманду /help.")
				bot.Send(msg)

			} else if update.Message.Command() == "help" {
				goletter = false
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Я отправляю письма в будущее. Для того чтобы начать вызовете команду /goletter")
				bot.Send(msg)
			} else if update.Message.Command() == "goletter" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Отправьте текст своего письма:")
				bot.Send(msg)
				goletter = true
				continue
			}
			if goletter == true {
				textMsg := update.Message.Text

				server.Store.CreateALetter(model.NewModel(textMsg, "2222-11-11", ""))
			}

		}
	}
}
