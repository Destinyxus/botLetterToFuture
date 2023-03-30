package apiserver

import (
	"fmt"
	"log"
	"time"

	"github.com/Destinyxus/botLetterToFuture/internal/email"
	"github.com/Destinyxus/botLetterToFuture/internal/encryptedLetter"
	"github.com/Destinyxus/botLetterToFuture/internal/store"
	"github.com/Destinyxus/botLetterToFuture/pkg/config"
)

type APIServer struct {
	Logger *log.Logger
	Store  *store.Store
	Email  *email.Email
	Config *config.Config
}

func NewAPIServer(config *config.Config) *APIServer {
	s := &APIServer{
		Logger: &log.Logger{},
		Store:  store.NewStore(),
		Config: config,
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

	if err := s.Store.Open(s.Config); err != nil {
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
				enc := encryptedLetter.NewEncrypter()
				decrypt, err := enc.Decrypt(letter.EncryptedLetter)
				if err != nil {
					return
				}

				s.Email.SendEmail(letter.Email, decrypt, s.Config)
				if err := s.Store.IsSent(letter.Email); err != nil {
					fmt.Errorf("error")
				}

			}
		}
	}
}
