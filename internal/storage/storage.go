package storage

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type Storage struct {
	conn *pgx.Conn
}

func New(conn *pgx.Conn) (*Storage, error) {
	st := &Storage{
		conn: conn,
	}

	return st, nil
}

func (s *Storage) InsertLetter(letter, email, date string) error {
	query := `INSERT INTO letters(letter,email,date,isActual)
			  values ($1,$2,$3,$4)`

	_, err := s.conn.Exec(context.Background(), query, letter, email, date, true)
	if err != nil {
		return err
	}

	return nil
}

//func (s *Storage) GetLetter
