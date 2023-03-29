package store

import (
	"database/sql"
	"fmt"
	"log"
	"time"

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
	fmt.Println("opened")
	s.db = db

	return nil
}

func (s *Store) CreateAccountTable() error {
	query := `create table if not exists letters (
    			id bigserial primary key,
    			email varchar(100) not null,
    			text_date date not null,
    			encrypted_letter varchar not null,
    			sent boolean default false
    )`

	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}
	fmt.Println("created")
	return err
}

func (s *Store) GetLetter() ([]*model.Model, error) {

	currentDate := time.Now().Format("2006-01-02")
	row, err := s.db.Query("SELECT email, encrypted_letter FROM letters WHERE text_date = $1 AND sent = false", currentDate)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	var letters []*model.Model
	model1 := &model.Model{}

	for row.Next() {
		err = row.Scan(&model1.Email, &model1.EncryptedLetter)
		if err != nil {
			log.Println(err)
			continue
		}
		letters = append(letters, model1)
	}

	return letters, nil
}

func (s *Store) IsSent(email string) error {
	_, err := s.db.Exec("UPDATE letters SET sent=true WHERE email = $1", email)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) CreateALetter(m *model.Model) error {
	query := fmt.Sprintf("insert into letters (email,text_date,encrypted_letter) values ('%s','%s','%s')", m.Email, m.Date, m.EncryptedLetter)

	if _, err := s.db.Exec(query); err != nil {
		return err
	}

	return nil

}
