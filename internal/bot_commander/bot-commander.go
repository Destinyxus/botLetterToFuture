package bot_commander

import (
	"context"
	"errors"
	"fmt"
	"github.com/Destinyxus/botLetterToFuture/internal/mapwmutex"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"log"
	"sync"
)

type BotCommander struct {
	logger      *logrus.Logger
	tg          *tgbotapi.BotAPI
	emailSender EmailSender
	repo        Repository
	userInfo    mapwmutex.MapWmutex[int64, []Letter]
	userState   mapwmutex.MapWmutex[int64, bool]
	dateIndex   mapwmutex.MapWmutex[string, []int64]
}

type EmailSender interface {
	SendEmail(email, letter string) error
}

type Repository interface {
	InsertLetter(letter, email, date string) error
}

type Letter struct {
	message  string
	Email    string
	date     string
	isActual bool
}

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("send the letter"),
		tgbotapi.NewKeyboardButton("/about me"),
	),
)

func New(
	repo Repository,
	options ...Option) *BotCommander {
	b := &BotCommander{
		userState: *mapwmutex.NewMapWmutex[int64, bool](0),
		dateIndex: *mapwmutex.NewMapWmutex[string, []int64](0),
		repo:      repo,
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

			if err := b.handleCommand(update.Message.From.ID, update.Message.Chat.ID, update.Message.Text); err != nil {
				log.Fatal(err)
			}

			continue
		}
	}()

	return nil
}

func (b *BotCommander) handleCommand(userId, chatId int64, message string) error {
	msg := tgbotapi.NewMessage(chatId, message)

	switch message {
	case "/start":
		b.userState.Store(userId, false)
		msg.ReplyMarkup = numericKeyboard

		if _, err := b.tg.Send(msg); err != nil {
			log.Panic(err)
		}

		break
	case "/open":
		b.userState.Store(userId, false)
		msg.ReplyMarkup = numericKeyboard

		if _, err := b.tg.Send(msg); err != nil {
			log.Panic(err)
		}

		break
	case "/about me":
		b.userState.Store(userId, false)

		msg.Text = "Я отправляю письма в будущее. Как это работает? Ты присылаешь мне письмо, которое хочешь получить на какую-то конкретную дату в будущем. Затем, тебе будет дана возможность указать свою почту, на которую ты желаешь получить то самое послание. Имей в виду, что после процесса отправки своего письма, я его удалю из чата, для того, чтобы сохранить интригу и дать твоему мозгу возможность забыть о нем! Ты оставляешь сто рублей в зимней куртке и благополучно забываешь о них, а через сезон надеваешь ее и ого, ничего себе - ты нащупываешь те самые рубли и радуешься! Здесь принцип схожий) P.S Дата должна быть в формате год-месяц-день и не может выходить за рамки определенного мной периода - настоящее время - 2024-03-30 Для того чтобы начать вызовете команду /goletter. Кроме того, если по какой-то причине на одном из этапов отправки письма вы передумали отправлять его, нажмите /stop."

		if _, err := b.tg.Send(msg); err != nil {
			log.Panic(err)
		}

		break
	case "send the letter":
		b.userState.Store(userId, true)

		msg.Text = "send me the email, date and letter"

		if _, err := b.tg.Send(msg); err != nil {
			log.Panic(err)
		}

		break
	case "/stop":
		b.userState.Store(userId, false)

		msg.Text = "Stopping the letter creation.If you want to start again just /open the menu"

		if _, err := b.tg.Send(msg); err != nil {
			log.Panic(err)
		}

		break
	default:
		if state := b.userState.Load(userId); state {
			letter, err := ValidateMessage(message)
			if err == nil {
				if err = b.repo.InsertLetter(letter.message, letter.Email, letter.date); err != nil {
					log.Fatal(err)
				}

				if _, err = b.tg.Send(msg); err != nil {
					log.Panic(err)
				}

				b.userState.Store(userId, false)

				break
			} else if errors.Is(err, ErrNotValidEmailOrDate) {
				msg.Text = "invalid email or date"

				if _, err = b.tg.Send(msg); err != nil {
					log.Panic(err)
				}

				break
			}
		} else {
			msg.Text = "i dont know this command"

			if _, err := b.tg.Send(msg); err != nil {
				log.Panic(err)
			}

			break
		}
	}

	return nil
}

//func (b *BotCommander) CheckForActualDate() error {
//	currentTime := time.Now()
//
//	formattedDate := currentTime.Format("2006-01-02")
//
//	return nil
//}
