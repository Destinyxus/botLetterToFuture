package main

import (
	"fmt"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"

	"LetterToFuture/internal/apiserver"
	"LetterToFuture/internal/model"
	"LetterToFuture/internal/telegram"
	"LetterToFuture/pkg"
)

const MAX_MESSAGE_LIMIT = 4095

var (
	helpText = "Я отправляю письма в будущее. Как это работает? " +
		"Ты присылаешь мне письмо, которое хочешь получить на какую-то конкретную дату в будущем. Затем, тебе будет дана " +
		"возможность указать свою почту, на которую ты желаешь получить то самое послание. " +
		"Имей в виду, что после процесса отправки своего письма, я его удалю из чата, для того, чтобы сохранить интригу и дать твоему мозгу возможность забыть " +
		"о нем! Ты оставляешь сто рублей в зимней куртке и благополучно забываешь о них, а через сезон надеваешь ее и \"ого, ничего себе\" - ты нащупываешь те самые рубли и радуешься! " +
		"Здесь принцип схожий)" +
		" Для того чтобы начать вызовете команду /goletter"
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

	commands := telegram.NewCommands()

	model2 := model.TemporaryModel()

	modelDelete := model.NewDeleteModel()

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
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, helpText)
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
							fmt.Println(model2.Letter)
							modelDelete.MessageId = update.Message.MessageID
							model2.Letter = update.Message.Text
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Теперь введите вашу почту:")
							bot.Send(msg)
						} else if model2.Email == "" {
							if pkg.ValidateEmail(update.Message.Text) != false {
								model2.Email = update.Message.Text
							} else {
								msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Проверьте правильность введенной почты!")
								bot.Send(msg)
								continue
							}
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Отлично, теперь введите дату:")
							bot.Send(msg)
						} else if model2.Date == "" {
							if pkg.DateValidation(update.Message.Text) != false {
								model2.Date = update.Message.Text
							} else {
								msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Проверьте правильность введенной даты! Дата должна быть в формате"+
									" yyyy-mm-dd. Кроме того, проверьте еще раз точность своей даты, я могу сохранять ваши письма только до 2025-03-28")
								bot.Send(msg)
								continue
							}
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Спасибо! Ваше письмо будет удалено из чата через пару минут! Если вы хотите отправить еще одно"+
								" письмо в будущее, нажмите /goletter. Удачи!")
							bot.Send(msg)
							server.Store.CreateALetter(model.NewModel(model2.Email, model2.Date, model2.Letter))
							fmt.Println(model2)
							pkg.UpdateStruct(model2)

							go func() {
								time.Sleep(time.Minute * 5)
								letterToDelete := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, modelDelete.MessageId)
								bot.Send(letterToDelete)
							}()

						}
					}

				}
			}
		} else {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "К сожалению, я принимаю только текстовый формат."+
				" Нажмите /help"))
		}
	}

}
