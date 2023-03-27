package pkg

import (
	"net/mail"
	"time"
)

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func DateValidation(s string) bool {
	_, err := time.Parse("2006-01-02", s)
	if err != nil {
		return false
	}
	return true
}
