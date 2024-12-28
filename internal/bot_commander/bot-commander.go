package bot_commander

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Destinyxus/botLetterToFuture/internal/config"
	"github.com/Destinyxus/botLetterToFuture/internal/mapwmutex"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	ErrDumpingDatesIndexes = errors.New("error dumping dates indexes")

	infoAboutDescription string
	infoResult           string
	infoStopCommand      string
	infoSendLetter       string

	errInvalidFormatMessage error
	errNotValidCommand      error

	numericKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("send the Letter"),
			tgbotapi.NewKeyboardButton("/about me"),
		),
	)
)

func New(ctx context.Context, storage Storage, cfg config.BotResponses, options ...Option) (*BotCommander, error) {
	b := &BotCommander{
		userState: *mapwmutex.NewMapWmutex[int64, bool](0),
		dateIndex: make(map[time.Time]struct{}),
		storage:   storage,
	}

	for _, opt := range options {
		opt(b)
	}

	if err := b.DatesDump(ctx); err != nil {
		return &BotCommander{}, fmt.Errorf("init botcommander error: %w", ErrDumpingDatesIndexes)
	}

	infoAboutDescription = cfg.AboutDescription
	infoResult = cfg.Result
	infoStopCommand = cfg.StopCommand
	infoSendLetter = cfg.SendLetter

	errInvalidFormatMessage = errors.New(cfg.InvalidFormatMessage)
	errNotValidCommand = errors.New(cfg.NotValidCommand)

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

			ctx := context.Background()

			if err := b.handleCommand(
				ctx,
				update.Message.From.ID,
				update.Message.Chat.ID,
				update.Message.MessageID,
				update.Message.Text); err != nil {

				b.logger.Debug("error while handling command")
			}

			continue
		}
	}()

	return nil
}

func (b *BotCommander) handleCommand(ctx context.Context, userId, chatId int64, messageID int, message string) error {
	msg := tgbotapi.NewMessage(chatId, message)

	switch message {
	case "/start":
		b.userState.Store(userId, false)
		msg.ReplyMarkup = numericKeyboard

		if _, err := b.tg.Send(msg); err != nil {
			b.logger.Debugf("sending keyboard message: %v", err)
		}
	case "/open":
		b.userState.Store(userId, false)
		msg.ReplyMarkup = numericKeyboard

		if _, err := b.tg.Send(msg); err != nil {
			b.logger.Debugf("sending keyboard message: %v", err)
		}
	case "/about me":
		b.userState.Store(userId, false)

		msg.Text = infoAboutDescription

		if _, err := b.tg.Send(msg); err != nil {
			b.logger.Debugf("sending about me info: %v", err)
		}
	case "send the Letter":
		b.userState.Store(userId, true)

		msg.Text = infoSendLetter

		if _, err := b.tg.Send(msg); err != nil {
			b.logger.Debugf("sending the offer to send message: %v", err)
		}
	case "/stop":
		b.userState.Store(userId, false)

		msg.Text = infoStopCommand

		if _, err := b.tg.Send(msg); err != nil {
			b.logger.Debugf("sending the stop info: %v", err)
		}
	default:
		b.processMessage(ctx, userId, chatId, messageID, message, msg)
	}

	return nil
}

func (b *BotCommander) processMessage(ctx context.Context, userId int64, chatId int64, messageID int, message string, msg tgbotapi.MessageConfig) {
	if state := b.userState.Load(userId); state {
		letter, err := ValidateMessage(message)
		if err == nil {
			if err = b.storage.InsertLetter(ctx, letter.Letter, letter.Email, letter.Date); err != nil {
				log.Fatal(err)
			}

			b.dateIndex[letter.Date] = struct{}{}

			b.userState.Store(userId, false)

			msg.Text = infoResult

			if _, err = b.tg.Send(msg); err != nil {
				b.logger.Debugf("sending the success message: %v", err)
			}

			if _, err = b.tg.Send(tgbotapi.NewDeleteMessage(chatId, messageID)); err != nil {
				b.logger.Debugf("deleting user's message: %v", err)
			}
		} else if errors.Is(err, errInvalidEmailOrDate) {
			msg.Text = errInvalidFormatMessage.Error()

			if _, err = b.tg.Send(msg); err != nil {
				b.logger.Debugf("sending the invalid message: %v", err)
			}
		}
	} else {
		msg.Text = errNotValidCommand.Error()

		if _, err := b.tg.Send(msg); err != nil {
			b.logger.Debugf("sending the not valid command message: %v", err)
		}
	}
}

func (b *BotCommander) CheckForActualDate(ctx context.Context) error {
	currentDate, err := time.Parse(dateFormat, time.Now().Format(dateFormat))
	if err != nil {
		return fmt.Errorf("error parsing time: %w", err)
	}

	if _, actual := b.dateIndex[currentDate]; actual {
		letters, err := b.storage.GetLetter(ctx, currentDate)
		if err != nil {
			return fmt.Errorf("error getting letter: %w", err)
		}

		for _, letter := range letters {
			if err = b.emailClient.SendEmail(letter.Email, letter.Letter); err != nil {
				return fmt.Errorf("error sending letters to emails: %w", err)
			}
		}

		b.logger.Infof("successful sending letters to emails: %d", len(letters))

		delete(b.dateIndex, currentDate)

		if err = b.storage.DeprecateLetter(ctx, currentDate); err != nil {
			return fmt.Errorf("error deprecating letter: %w", err)
		}
	}

	return nil
}

func (b *BotCommander) DatesDump(ctx context.Context) error {
	dates, err := b.storage.GetActualDates(ctx)
	if err != nil {
		return fmt.Errorf("getting actual dates dump: %w", err)
	}

	b.dateIndex = dates

	b.logger.Infof("successfully dumped date indexes: %d", len(b.dateIndex))

	return nil
}
