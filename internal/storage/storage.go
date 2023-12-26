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

	query := `CREATE TABLE IF NOT EXISTS letters (
    id SERIAL PRIMARY KEY,
    letter TEXT NOT NULL,
    letter_date TIMESTAMP NOT NULL,
    email VARCHAR(255) NOT NULL,
    isActual BOOLEAN NOT NULL
)`

	if _, err := st.conn.Exec(context.Background(), query); err != nil {
		return st, err
	}

	return st, nil
}

func (s *Storage) InsertLetter() error {
	return nil
}
