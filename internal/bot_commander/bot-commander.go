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
	EmailSender EmailSender
	Repo        Repository
	userState   mapwmutex.MapWmutex[int64, bool]
	Logger      *logrus.Logger
	tg          *tgbotapi.BotAPI
	DateIndex   map[time.Time]struct{}
	cfg         config.Config
}

type EmailSender interface {
	SendEmail(email, letter string) error
}

type Repository interface {
	InsertLetter(letter, email string, date time.Time) error
	GetLetter(date time.Time) ([]Letter, error)
	GetActualDates() (map[time.Time]struct{}, error)
	DeprecateLetter(date time.Time) error
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

var DateIndexesError = errors.New("date indexes err")

func New(
	repo Repository,
	cfg config.Config,
	options ...Option,
) (*BotCommander, error) {
	b := &BotCommander{
		userState: *mapwmutex.NewMapWmutex[int64, bool](0),
		DateIndex: make(map[time.Time]struct{}),
		Repo:      repo,
		cfg:       cfg,
	}

	for _, o := range options {
		if err := o(b); err != nil {
			fmt.Println("")
		}
	}

	if err := b.DatesDump(); err != nil {
		return &BotCommander{}, DateIndexesError
	}

	return b, nil
}

func (b *BotCommander) Start(ctx context.Context, wg *sync.WaitGroup) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.tg.GetUpdatesChan(u)

	wg.Add(1)

	go func() {
		defer wg.Done()

		<-ctx.Done()
		b.tg.StopReceivingUpdates()
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()

		for update := range updates {
			if update.Message == nil {
				continue
			}

			if err := b.handleCommand(
				update.Message.From.ID,
				update.Message.Chat.ID,
				update.Message.MessageID,
				update.Message.Text); err != nil {

				b.Logger.Debug("error while handling command")
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
			b.Logger.Debugf("sending keyboard message: %v", err)
		}
	case "/open":
		b.userState.Store(userId, false)
		msg.ReplyMarkup = numericKeyboard

		if _, err := b.tg.Send(msg); err != nil {
			b.Logger.Debugf("sending keyboard message: %v", err)
		}
	case "/about me":
		b.userState.Store(userId, false)

		msg.Text = b.cfg.Responses.AboutDescription

		if _, err := b.tg.Send(msg); err != nil {
			b.Logger.Debugf("sending about me info: %v", err)
		}
	case "send the Letter":
		b.userState.Store(userId, true)

		msg.Text = b.cfg.Responses.SendLetter

		if _, err := b.tg.Send(msg); err != nil {
			b.Logger.Debugf("sending the offer to send message: %v", err)
		}
	case "/stop":
		b.userState.Store(userId, false)

		msg.Text = b.cfg.Responses.StopCommand

		if _, err := b.tg.Send(msg); err != nil {
			b.Logger.Debugf("sending the stop info: %v", err)
		}
	default:
		b.processMessage(userId, chatId, messageID, message, msg)
	}

	return nil
}

func (b *BotCommander) processMessage(userId int64, chatId int64, messageID int, message string, msg tgbotapi.MessageConfig) {
	if state := b.userState.Load(userId); state {
		letter, err := ValidateMessage(message)
		if err == nil {
			if err = b.Repo.InsertLetter(letter.Letter, letter.Email, letter.Date); err != nil {
				log.Fatal(err)
			}

			b.DateIndex[letter.Date] = struct{}{}

			b.userState.Store(userId, false)

			msg.Text = b.cfg.Responses.Result

			if _, err = b.tg.Send(msg); err != nil {
				b.Logger.Debugf("sending the success message: %v", err)
			}

			if _, err = b.tg.Send(tgbotapi.NewDeleteMessage(chatId, messageID)); err != nil {
				b.Logger.Debugf("deleting user's message: %v", err)
			}
		} else if errors.Is(err, ErrNotValidEmailOrDate) {
			msg.Text = b.cfg.Errors.InvalidFormatMessage

			if _, err = b.tg.Send(msg); err != nil {
				b.Logger.Debugf("sending the invalid message: %v", err)
			}
		}
	} else {
		msg.Text = b.cfg.Errors.NotValidCommand

		if _, err := b.tg.Send(msg); err != nil {
			b.Logger.Debugf("sending the not valid command message: %v", err)
		}
	}
}

func (b *BotCommander) CheckForActualDate() error {
	now := time.Now().Format(DateFormat)

	currentDate, err := time.Parse(DateFormat, now)
	if err != nil {
		return err
	}

	if _, actual := b.DateIndex[currentDate]; actual {
		letters, err := b.Repo.GetLetter(currentDate)
		if err != nil {
			return fmt.Errorf("getting the letter with date: %w", err)
		}

		for _, letter := range letters {
			if err = b.EmailSender.SendEmail(letter.Email, letter.Letter); err != nil {
				return fmt.Errorf("sending the letters to emails: %w", err)
			}
		}

		b.Logger.Infof("successful sending letters to emails: %d", len(letters))

		delete(b.DateIndex, currentDate)

		if err = b.Repo.DeprecateLetter(currentDate); err != nil {
			return fmt.Errorf("deprecating not actual letters with date: %w", err)
		}
	}

	return nil
}

func (b *BotCommander) DatesDump() error {
	dates, err := b.Repo.GetActualDates()
	if err != nil {
		return fmt.Errorf("getting actual dates dump: %w", err)
	}

	b.DateIndex = dates

	b.Logger.Infof("successfully dumped date indexes: %d", len(b.DateIndex))

	return nil
}
