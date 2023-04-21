package telegram

import (
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Destinyxus/botLetterToFuture/internal/apiserver"
	"github.com/Destinyxus/botLetterToFuture/internal/encryptedLetter"
	"github.com/Destinyxus/botLetterToFuture/internal/model"

	"github.com/Destinyxus/botLetterToFuture/pkg"
	"github.com/Destinyxus/botLetterToFuture/pkg/config"
)

const (
	MaxMessageLimit = 4095
)

var (
	commands      = NewCommands()
	temporaryUser = model.TemporaryUser()
	configUser    = model.NewConfigUser()
	userToDelete  = model.NewDeleteModel()
)

func Init() {
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

	for update := range updates {
		if update.Message != nil && update.Message.Text != "" {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			switch update.Message.Command() {
			case "start":
				handleStart(bot, update.Message.Command(), update.Message.Chat.ID, cfg)
			case "help":
				handleHelp(bot, update.Message.Command(), update.Message.Chat.ID, cfg)
			case "goletter":
				handleGoLetter(bot, update.Message.Command(), update.Message.Chat.ID, cfg)
			case "stop":
				handleStop(bot, update.Message.Command(), update.Message.Chat.ID, update.Message.MessageID, cfg)
				continue
			default:
				if commands.Start == true {
					handleIsStart(bot, update.Message.Chat.ID, cfg)
				}
				if commands.Help == true {
					handleIsHelp(bot, update.Message.Chat.ID, cfg)
				}
				if commands.Goletter == true {
					handleGoLetterFinal(bot, update.Message.Command(), update.Message.Chat.ID, update.Message.MessageID, update.Message.Text, update.Message.Chat.UserName, cfg, server)
				}
			}
		} else if update.Message != nil {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "К сожалению, я принимаю только текстовый формат."+
				" Нажмите /help"))
		}
	}

}

func handleStart(bot *tgbotapi.BotAPI, updateCommand string, messageChatID int64, cfg *config.Config) {
	commands.CommandMode(updateCommand)
	msg := tgbotapi.NewMessage(messageChatID, cfg.Messages.Start)
	bot.Send(msg)
}

func handleHelp(bot *tgbotapi.BotAPI, updateCommand string, messageChatID int64, cfg *config.Config) {
	commands.CommandMode(updateCommand)
	msg := tgbotapi.NewMessage(messageChatID, cfg.Messages.HelpText)
	bot.Send(msg)
}

func handleGoLetter(bot *tgbotapi.BotAPI, updateCommand string, messageChatID int64, cfg *config.Config) {
	commands.CommandMode(updateCommand)
	msg := tgbotapi.NewMessage(messageChatID, cfg.Messages.Goletter)
	bot.Send(msg)
	temporaryUser.Letter = ""
	temporaryUser.Email = ""
	temporaryUser.Date = ""
}

func handleStop(bot *tgbotapi.BotAPI, updateCommand string, messageChatID int64, messageID int, cfg *config.Config) {
	commands.CommandMode(updateCommand)
	deleteMsg := tgbotapi.NewDeleteMessage(messageChatID, messageID)
	bot.Send(deleteMsg)
	msg := tgbotapi.NewMessage(messageChatID, "Сессия окончена, если вы желаете начать сначала, нажмите /start.")
	bot.Send(msg)
	pkg.UpdateStruct(temporaryUser)

	letterToDelete := tgbotapi.NewDeleteMessage(messageChatID, userToDelete.LetterId)
	dateToDelete := tgbotapi.NewDeleteMessage(messageChatID, userToDelete.DateId)
	emailToDelete := tgbotapi.NewDeleteMessage(messageChatID, userToDelete.EmailId)

	bot.Send(letterToDelete)
	bot.Send(dateToDelete)
	bot.Send(emailToDelete)

}

func handleIsStart(bot *tgbotapi.BotAPI, messageChatID int64, cfg *config.Config) {
	msg := tgbotapi.NewMessage(messageChatID, cfg.Messages.StartTrue)
	bot.Send(msg)
}

func handleIsHelp(bot *tgbotapi.BotAPI, messageChatID int64, cfg *config.Config) {
	msg := tgbotapi.NewMessage(messageChatID, cfg.Messages.HelpTrue)
	bot.Send(msg)
}

func handleGoLetterFinal(bot *tgbotapi.BotAPI, updateCommand string, messageChatID int64, messageID int, textMessage string, userName string, cfg *config.Config, server *apiserver.APIServer) {
	textMsg := textMessage
	if len(textMsg) > MaxMessageLimit {
		msg := tgbotapi.NewMessage(messageChatID, cfg.Messages.SizeLetter)
		bot.Send(msg)
		return
	}
	if temporaryUser.Letter == "" {
		userToDelete.LetterId = messageID
		temporaryUser.Letter = textMessage
		msg := tgbotapi.NewMessage(messageChatID, cfg.Messages.Email)
		bot.Send(msg)
	} else if temporaryUser.Email == "" {
		if pkg.ValidateEmail(textMessage) != false {
			userToDelete.EmailId = messageID
			temporaryUser.Email = textMessage
			msg := tgbotapi.NewMessage(messageChatID, cfg.Messages.Date)
			bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(messageChatID, cfg.Messages.InvalidEmail)
			bot.Send(msg)
			return
		}
	} else if temporaryUser.Date == "" {
		if pkg.DateValidation(textMessage) != false {
			userToDelete.DateId = messageID
			temporaryUser.Date = textMessage
		} else {
			msg := tgbotapi.NewMessage(messageChatID, cfg.Messages.InvalidDate)
			bot.Send(msg)
			return
		}
		msg := tgbotapi.NewMessage(messageChatID, cfg.Messages.Result)
		bot.Send(msg)
		enc := encryptedLetter.NewEncrypter()
		encrypt, err := enc.Encrypt(temporaryUser.Letter, cfg.HashKey)
		if err != nil {
			return
		}
		if err != nil {
			log.Fatal(err)
		}
		configUser.UserName = userName
		err = server.Store.CreateALetter(model.NewUser(temporaryUser.Email, temporaryUser.Date, encrypt), configUser)
		if err != nil {
			log.Fatal(err)
		}

		pkg.UpdateStruct(temporaryUser)
		commands.CommandMode("reset")

		go func() {
			time.Sleep(time.Second * 3)
			letterToDelete := tgbotapi.NewDeleteMessage(messageChatID, userToDelete.LetterId)
			dateToDelete := tgbotapi.NewDeleteMessage(messageChatID, userToDelete.DateId)
			emailToDelete := tgbotapi.NewDeleteMessage(messageChatID, userToDelete.EmailId)
			bot.Send(letterToDelete)
			bot.Send(dateToDelete)
			bot.Send(emailToDelete)
		}()
	}
}
