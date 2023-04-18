package store

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Destinyxus/botLetterToFuture/internal/model"
	"github.com/Destinyxus/botLetterToFuture/pkg/config"

	_ "github.com/lib/pq"
)

type Letter struct {
	ID              int    `gorm:"primaryKey;autoIncrement"`
	Email           string `gorm:"not null"`
	Date            string `gorm:"not null"`
	EncryptedLetter string `gorm:"not null"`
	Sent            bool   `gorm:"default:false"`
}

type Store struct {
	db *gorm.DB
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) Open(cfg *config.Config) error {
	db, err := gorm.Open(postgres.Open(cfg.StoreURL), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to db")
	}
	s.db = db
	err = s.db.AutoMigrate(&Letter{})
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) CreateALetter(m *model.Model) error {
	letter := &Letter{
		Email:           m.Email,
		Date:            m.Date,
		EncryptedLetter: m.EncryptedLetter,
	}
	if err := s.db.Create(letter).Error; err != nil {
		return err

	}
	return nil
}

func (s *Store) GetLetter() ([]*Letter, error) {
	currentDate := time.Now().Format("2006-01-02")
	var letters []*Letter
	if err := s.db.Where("date = ? AND sent = ?", currentDate, false).Find(&letters).Error; err != nil {
		return nil, err
	}
	return letters, nil

}

func (s *Store) IsSent(id int) error {
	letter := &Letter{}
	result := s.db.Model(letter).Where("id = ?", id).Update("sent", true)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
