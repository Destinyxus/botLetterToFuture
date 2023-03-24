package store

import (
	"database/sql"
	"fmt"

	"LetterToFuture/internal/model"

	_ "github.com/lib/pq"
)

type Store struct {
	db *sql.DB
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) Open() error {
	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=future_letter password=10120001 sslmode=disable")
	if err != nil {
		return err

	}

	if err := db.Ping(); err != nil {
		return err
	}
	fmt.Println("openned")
	s.db = db
	return nil
}

func (s *Store) CreateAccountTable() error {
	query := `create table if not exists letters (
    			id bigserial primary key,
    			email varchar(100) unique not null,
    			text_date date not null,
    			letter varchar not null
    )`

	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}
	fmt.Println("gleb")
	return err
}

func (s *Store) CreateALetter(m *model.Model) error {
	query := fmt.Sprintf("insert into letters (email,text_date,letter) values ('%s','%s','%s')", m.Email, m.Date, m.Letter)

	if _, err := s.db.Exec(query); err != nil {
		return err
	}

	return nil

}
