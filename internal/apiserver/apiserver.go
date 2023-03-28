package apiserver

import (
	"fmt"
	"log"
	"time"

	"LetterToFuture/internal/email"
	"LetterToFuture/internal/store"
)

type APIServer struct {
	Logger *log.Logger
	Store  *store.Store
	Email  *email.Email
}

func NewAPIServer() *APIServer {
	s := &APIServer{
		Logger: &log.Logger{},
		Store:  store.NewStore(),
	}
	s.Email = email.NewEmail(s.Store)
	return s
}

func (s *APIServer) Start() error {
	s.configureStore()
	go s.configureEmail()
	return nil
}

func (s *APIServer) configureStore() error {

	if err := s.Store.Open(); err != nil {
		return err
	}

	if err := s.Store.CreateAccountTable(); err != nil {
		return err

	}

	return nil
}

func (s *APIServer) configureEmail() {
	ticker := time.NewTicker(1 * time.Minute)

	for {
		select {
		case <-ticker.C:
			letters, err := s.Store.GetLetter()
			if err != nil {
				log.Println(err)
				continue
			}
			for _, letter := range letters {
				fmt.Println(letter.Letter)
				s.Email.SendEmail(letter.Email, letter.Letter)
				if err := s.Store.IsSent(letter.Email); err != nil {
					fmt.Errorf("error")
				}

			}
		}
	}
}
