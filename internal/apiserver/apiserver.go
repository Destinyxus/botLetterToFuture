package apiserver

import (
	"log"

	"LetterToFuture/internal/store"
)

type APIServer struct {
	Logger *log.Logger
	Store  *store.Store
}

func NewAPIServer() *APIServer {
	return &APIServer{
		Logger: &log.Logger{},
		Store:  store.NewStore(),
	}
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
