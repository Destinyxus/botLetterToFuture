package main

import (
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Destinyxus/botLetterToFuture/internal/apiserver"
	"github.com/Destinyxus/botLetterToFuture/internal/encryptedLetter"
	"github.com/Destinyxus/botLetterToFuture/internal/model"
	"github.com/Destinyxus/botLetterToFuture/internal/telegram"
	"github.com/Destinyxus/botLetterToFuture/pkg"
	"github.com/Destinyxus/botLetterToFuture/pkg/config"
)

const MAX_MESSAGE_LIMIT = 4095

func main() {

	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	server := apiserver.NewAPIServer(cfg)

	if err := server.Start(); err != nil {
		server.Logger.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
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

	modelDelete := model.NewDeleteModel()

	for update := range updates {
		if update.Message != nil && update.Message.Text != "" {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			switch update.Message.Command() {
			case "start":
				commands.CommandMode(update.Message.Command())
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, cfg.Messages.Start)
				bot.Send(msg)
			case "help":
				commands.CommandMode(update.Message.Command())
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, cfg.Messages.HelpText)
				bot.Send(msg)
			case "goletter":
				commands.CommandMode(update.Message.Command())
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, cfg.Messages.Goletter)
				bot.Send(msg)
				model2.Letter = ""
				model2.Email = ""
				model2.Date = ""
			case "stop":
				commands.CommandMode(update.Message.Command())
				deleteMsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
				bot.Send(deleteMsg)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Сессия окончена, если вы желаете начать сначала, нажмите /start.")
				bot.Send(msg)
				pkg.UpdateStruct(model2)
				letterToDelete := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, modelDelete.LetterId)
				dateToDelete := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, modelDelete.DateId)
				emailToDelete := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, modelDelete.EmailId)
				bot.Send(letterToDelete)
				bot.Send(dateToDelete)
				bot.Send(emailToDelete)
				continue
			default:
				if commands.Start == true {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, cfg.Messages.StartTrue)
					bot.Send(msg)
				}
				if commands.Help == true {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, cfg.Messages.HelpTrue)
					bot.Send(msg)
				}
				if commands.Goletter == true {
					textMsg := update.Message.Text
					if len(textMsg) > MAX_MESSAGE_LIMIT {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, cfg.Messages.SizeLetter)
						bot.Send(msg)
						continue
					} else {
						if model2.Letter == "" {
							modelDelete.LetterId = update.Message.MessageID
							model2.Letter = update.Message.Text
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, cfg.Messages.Email)
							bot.Send(msg)
						} else if model2.Email == "" {
							if pkg.ValidateEmail(update.Message.Text) != false {
								modelDelete.EmailId = update.Message.MessageID
								model2.Email = update.Message.Text
								msg := tgbotapi.NewMessage(update.Message.Chat.ID, cfg.Messages.Date)
								bot.Send(msg)
							} else {
								msg := tgbotapi.NewMessage(update.Message.Chat.ID, cfg.Messages.InvalidEmail)
								bot.Send(msg)
								continue
							}
						} else if model2.Date == "" {
							if pkg.DateValidation(update.Message.Text) != false {
								modelDelete.DateId = update.Message.MessageID
								model2.Date = update.Message.Text
							} else {
								msg := tgbotapi.NewMessage(update.Message.Chat.ID, cfg.Messages.InvalidDate)
								bot.Send(msg)
								continue
							}
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, cfg.Messages.Result)
							bot.Send(msg)
							enc := encryptedLetter.NewEncrypter()
							encrypt, err := enc.Encrypt(model2.Letter, cfg.HashKey)
							if err != nil {
								return
							}
							if err != nil {
								log.Fatal(err)
							}

							server.Store.CreateALetter(model.NewModel(model2.Email, model2.Date, encrypt))

							pkg.UpdateStruct(model2)
							commands.CommandMode("reset")

							go func() {
								time.Sleep(time.Minute * 5)
								letterToDelete := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, modelDelete.LetterId)
								dateToDelete := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, modelDelete.DateId)
								emailToDelete := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, modelDelete.EmailId)
								bot.Send(letterToDelete)
								bot.Send(dateToDelete)
								bot.Send(emailToDelete)
							}()

						}
					}

				}
			}
		} else if update.Message != nil {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "К сожалению, я принимаю только текстовый формат."+
				" Нажмите /help"))
		}
	}

}
