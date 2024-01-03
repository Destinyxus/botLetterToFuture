package bot_commander

import (
	"errors"
	"net/mail"
	"strings"
	"time"
)

const (
	dateFormat     = "2006-01-02"
	fromConstraint = "2023-03-28"
	toConstraint   = "2025-03-28"
)

var ErrNotValidEmailOrDate = errors.New("not valid Email or Date")

func ValidateMessage(message string) (Letter, error) {
	s := strings.Split(message, ";")

	if len(s) != 3 {
		return Letter{}, ErrNotValidEmailOrDate
	}

	letter := s[0]
	email := s[1]
	date := s[2]

	if !ValidateEmail(email) || !DateValidation(date) {
		return Letter{}, ErrNotValidEmailOrDate
	}

	datee, err := time.Parse(dateFormat, date)
	if err != nil {
		return Letter{}, err
	}

	return Letter{Letter: letter, Email: email, Date: datee}, nil
}

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)

	return err == nil
}

func DateValidation(datee string) bool {
	date, err := time.Parse(dateFormat, datee)
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
