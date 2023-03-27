package apiserver

import (
	"log"

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
