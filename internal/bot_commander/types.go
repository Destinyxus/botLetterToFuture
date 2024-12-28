package bot_commander

import (
	"context"
	"time"

	"github.com/Destinyxus/botLetterToFuture/internal/mapwmutex"
)

type EmailClient interface {
	SendEmail(email, letter string) error
}

type Storage interface {
	InsertLetter(ctx context.Context, letter, email string, date time.Time) error
	GetLetter(ctx context.Context, date time.Time) ([]Letter, error)
	GetActualDates(ctx context.Context) (map[time.Time]struct{}, error)
	DeprecateLetter(ctx context.Context, date time.Time) error
}

type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
}

type BotClient interface {
}

type BotCommander struct {
	emailClient EmailClient
	storage     Storage
	userState   mapwmutex.MapWmutex[int64, bool]
	logger      Logger
	tg          BotClient
	dateIndex   map[time.Time]struct{}
}

type UpdatesChannel <-chan Update

type Update struct {
	Message *Message
}

type Message struct {
	MessageID int
	From      *User
	Chat      *Chat
	Text      string
}

type User struct {
	ID int64
}

type Chat struct {
	ID int64
}

type Letter struct {
	Id       int       `db:"id"`
	Letter   string    `db:"letter"`
	Email    string    `db:"email"`
	Date     time.Time `db:"date"`
	IsActual bool      `db:"isactual"`
}
