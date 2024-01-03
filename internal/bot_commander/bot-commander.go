package bot_commander

import (
	"context"
	"errors"
	"fmt"
	"github.com/Destinyxus/botLetterToFuture/internal/config"
	"github.com/Destinyxus/botLetterToFuture/internal/mapwmutex"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"log"
	"sync"
	"time"
)

type BotCommander struct {
	emailSender EmailSender
	repo        Repository
	userState   mapwmutex.MapWmutex[int64, bool]
	logger      *logrus.Logger
	tg          *tgbotapi.BotAPI
	dateIndex   map[time.Time]struct{}
	cfg         config.Config
}

type EmailSender interface {
	SendEmail(email, letter string) error
}

type Repository interface {
	InsertLetter(letter, email string, date time.Time) error
	GetLetter(date time.Time) ([]Letter, error)
}

type Letter struct {
	Id       int       `db:"id"`
	Letter   string    `db:"letter"`
	Email    string    `db:"email"`
	Date     time.Time `db:"date"`
	IsActual bool      `db:"isactual"`
}

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("send the Letter"),
		tgbotapi.NewKeyboardButton("/about me"),
	),
)

func New(
	repo Repository,
	cfg config.Config,
	options ...Option) *BotCommander {
	b := &BotCommander{
		userState: *mapwmutex.NewMapWmutex[int64, bool](0),
		dateIndex: make(map[time.Time]struct{}),
		repo:      repo,
		cfg:       cfg,
	}

	for _, o := range options {
		if err := o(b); err != nil {
			fmt.Println("")
		}
	}

	return b
}

func (b *BotCommander) Start(ctx context.Context, wg *sync.WaitGroup) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.tg.GetUpdatesChan(u)

	wg.Add(1)

	go func() {
		defer wg.Done()

		for update := range updates {
			select {
			case <-ctx.Done():
				b.logger.Info("finishing app")

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
				msg.Text = "sorry, some unexpected error, i am going to sleep"

				if _, err := b.tg.Send(msg); err != nil {
					log.Panic(err)
				}

				return
			default:

			}

			if update.Message == nil {
				continue
			}

			if err := b.handleCommand(update.Message.From.ID, update.Message.Chat.ID, update.Message.MessageID, update.Message.Text); err != nil {
				log.Fatal(err)
			}

			continue
		}
	}()

	return nil
}

func (b *BotCommander) handleCommand(userId, chatId int64, messageID int, message string) error {
	msg := tgbotapi.NewMessage(chatId, message)

	switch message {
	case "/start":
		b.userState.Store(userId, false)
		msg.ReplyMarkup = numericKeyboard

		if _, err := b.tg.Send(msg); err != nil {
			b.logger.Debugf("sending keyboard message: %v", err)
		}

		break
	case "/open":
		b.userState.Store(userId, false)
		msg.ReplyMarkup = numericKeyboard

		if _, err := b.tg.Send(msg); err != nil {
			b.logger.Debugf("sending keyboard message: %v", err)
		}

		break
	case "/about me":
		b.userState.Store(userId, false)

		msg.Text = b.cfg.Responses.AboutDescription

		if _, err := b.tg.Send(msg); err != nil {
			b.logger.Debugf("sending about me info: %v", err)
		}

		break
	case "send the Letter":
		b.userState.Store(userId, true)

		msg.Text = b.cfg.Responses.SendLetter

		if _, err := b.tg.Send(msg); err != nil {
			b.logger.Debugf("sending the offer to send message: %v", err)
		}

		break
	case "/stop":
		b.userState.Store(userId, false)

		msg.Text = b.cfg.Responses.StopCommand

		if _, err := b.tg.Send(msg); err != nil {
			b.logger.Debugf("sending the stop info: %v", err)
		}

		break
	default:
		if state := b.userState.Load(userId); state {
			letter, err := ValidateMessage(message)
			if err == nil {
				if err = b.repo.InsertLetter(letter.Letter, letter.Email, letter.Date); err != nil {
					log.Fatal(err)
				}

				b.dateIndex[letter.Date] = struct{}{}

				b.userState.Store(userId, false)

				msg.Text = b.cfg.Responses.Result

				if _, err = b.tg.Send(msg); err != nil {
					b.logger.Debugf("sending the success message: %v", err)
				}

				if _, err = b.tg.Send(tgbotapi.NewDeleteMessage(chatId, messageID)); err != nil {
					b.logger.Debugf("deleting user's message: %v", err)
				}

				break
			} else if errors.Is(err, ErrNotValidEmailOrDate) {
				msg.Text = b.cfg.Errors.InvalidFormatMessage

				if _, err = b.tg.Send(msg); err != nil {
					b.logger.Debugf("sending the invalid message: %v", err)
				}

				break
			}
		} else {
			msg.Text = b.cfg.Errors.NotValidCommand

			if _, err := b.tg.Send(msg); err != nil {
				b.logger.Debugf("sending the not valid command message: %v", err)
			}

			break
		}
	}

	return nil
}

func (b *BotCommander) CheckForActualDate() error {
	now := time.Now().Format(dateFormat)

	currentDate, err := time.Parse(dateFormat, now)
	if err != nil {
		return err
	}

	if _, actual := b.dateIndex[currentDate]; actual {
		letters, err := b.repo.GetLetter(currentDate)
		if err != nil {
			return err
		}

		for _, letter := range letters {
			if err = b.emailSender.SendEmail(letter.Email, letter.Letter); err != nil {
				return err
			}
		}

		delete(b.dateIndex, currentDate)

		return nil
	}

	t := time.Date(
		0001,
		1,
		1,
		00,
		00,
		00,
		00,
		time.UTC)

	letters, err := b.repo.GetLetter(t)
	if err != nil {
		return err
	}

	for _, letter := range letters {
		if err = b.emailSender.SendEmail(letter.Email, letter.Letter); err != nil {
			return err
		}
	}

	return nil
}
