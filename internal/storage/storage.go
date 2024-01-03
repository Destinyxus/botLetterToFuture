package storage

import (
	"github.com/Destinyxus/botLetterToFuture/internal/bot_commander"
	"github.com/jmoiron/sqlx"
	"time"
)

type Storage struct {
	conn *sqlx.DB
}

func New(conn *sqlx.DB) (*Storage, error) {
	st := &Storage{
		conn: conn,
	}

	return st, nil
}

func (s *Storage) InsertLetter(letter, email string, date time.Time) error {
	query := `INSERT INTO letters(letter,email,date,isActual)
			  values ($1,$2,$3,$4)`

	_, err := s.conn.Exec(query, letter, email, date, true)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetLetter(date time.Time) ([]bot_commander.Letter, error) {
	query := `
			SELECT * FROM letters WHERE date = $1 and isActual is true

`

	if date.IsZero() {
		letters, err := s.getWithActualDate(query)
		if err != nil {
			return letters, err
		}

		return letters, nil
	}

	var letter []bot_commander.Letter

	if err := s.conn.Select(&letter, query, date); err != nil {
		return nil, err
	}

	setToFalse := `UPDATE letters SET isActual = false WHERE date = $1`

	if _, err := s.conn.Exec(setToFalse, date); err != nil {
		return nil, err
	}

	return letter, nil
}

func (s *Storage) getWithActualDate(query string) ([]bot_commander.Letter, error) {
	now := time.Now().Format("2006-01-02")

	currentDate, err := time.Parse("2006-01-02", now)
	if err != nil {
		return nil, err
	}

	var letter []bot_commander.Letter

	if err = s.conn.Select(&letter, query, currentDate); err != nil {
		return nil, err
	}

	if len(letter) != 0 {
		setToFalse := `UPDATE letters SET isActual = false WHERE date = $1`

		if _, err = s.conn.Exec(setToFalse, currentDate); err != nil {
			return nil, err
		}

		return letter, nil
	}

	return nil, nil
}
