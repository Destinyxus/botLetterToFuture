package apiserver

import (
	"log"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/Destinyxus/botLetterToFuture/internal/email"
	"github.com/Destinyxus/botLetterToFuture/internal/encryptedLetter"
	"github.com/Destinyxus/botLetterToFuture/internal/store"
	"github.com/Destinyxus/botLetterToFuture/pkg/config"
)

type APIServer struct {
	Logger *logrus.Logger
	Store  *store.Store
	Email  *email.Email
	Config *config.Config
	log    *logrus.Entry
}

func NewAPIServer(config *config.Config) *APIServer {
	logger := logrus.New()
	logEntry := logrus.NewEntry(logger).WithFields(logrus.Fields{
		"package": "apiServer",
	})
	return &APIServer{
		log:    logEntry,
		Logger: logEntry.Logger,
		Store:  store.NewStore(logEntry),
		Config: config,
		Email:  email.NewEmail(store.NewStore(logEntry)),
	}
}

func (s *APIServer) Start() error {
	err := s.configureLogger()
	if err != nil {
		return err
	}
	err = s.configureStore()
	if err != nil {
		s.log.Error("Configuring store error")
	}
	go s.isDateActual()
	s.log.Info("Bot has been started!")
	return nil
}

func (s *APIServer) configureLogger() error {
	s.Logger.SetLevel(logrus.DebugLevel)
	file, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	s.Logger.SetOutput(file)
	return nil
}
func (s *APIServer) configureStore() error {
	if err := s.Store.Open(s.Config); err != nil {
		return err
	}
	return nil
}

func (s *APIServer) isDateActual() {
	ticker := time.NewTicker(1 * time.Minute)

	for {
		select {
		case <-ticker.C:
			letters, err := s.Store.GetLetter()
			if err != nil {
				log.Println(err)
				continue
			}
			s.sendEmail(letters)
		}
	}
}

func (s *APIServer) sendEmail(letters []*store.Letter) error {
	for _, letter := range letters {
		enc := encryptedLetter.NewEncrypter()
		decrypt, err := enc.Decrypt(letter.EncryptedLetter)
		if err != nil {
			return err
		}

		s.Email.SendEmail(letter.Email, decrypt, s.Config)
		if err := s.Store.IsSent(letter.ID); err != nil {
			return err
		}

	}
	return nil
}
