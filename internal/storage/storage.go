package storage

import (
	"github.com/jackc/pgx/v5"
)

type Storage struct {
	conn *pgx.Conn
}

func New(conn *pgx.Conn) *Storage {
	return &Storage{
		conn: conn,
	}
}

func (s *Storage) InsertLetter() error {
	return nil
}
