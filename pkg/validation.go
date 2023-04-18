package pkg

import (
	"net/mail"
	"time"
)

const (
	dateFormat     = "2006-01-02"
	fromConstraint = "2023-03-28"
	toConstraint   = "2025-03-28"
)

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func DateValidation(s string) bool {
	date, err := time.Parse(dateFormat, s)
	if err != nil {
		return false
	}

	from, err := time.Parse(dateFormat, fromConstraint)
	if err != nil {
		return false
	}

	to, err := time.Parse(dateFormat, toConstraint)
	if err != nil {
		return false
	}

	now := time.Now().Format(dateFormat)
	currentDate, err := time.Parse(dateFormat, now)
	if err != nil {
		return false
	}

	if date.Before(from) || date.After(to) || date.Before(currentDate) {
		return false
	}
	return true
}
