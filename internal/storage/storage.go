package storage

import (
	"time"

	"github.com/Destinyxus/botLetterToFuture/internal/bot_commander"
	"github.com/jmoiron/sqlx"
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
	SELECT * 
	FROM letters 
	WHERE date = $1 
		AND isActual is true`

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

func (s *Storage) GetActualDates() (map[time.Time]struct{}, error) {
	query := `SELECT DISTINCT date 
			  FROM letters 
			  WHERE date >= CURRENT_DATE`

	rows, err := s.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dates := make([]time.Time, 0)

	for rows.Next() {
		var date time.Time

		if err = rows.Scan(&date); err != nil {
			return nil, err
		}

		dates = append(dates, date)
	}

	dateIndexes := make(map[time.Time]struct{})

	for _, d := range dates {
		dateIndexes[d] = struct{}{}
	}

	return dateIndexes, nil
}
