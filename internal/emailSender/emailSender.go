package emailSender

import (
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Email struct {
	client    *sendgrid.Client
	emailFrom *mail.Email
}

func NewEmail(token, mailName, address string) *Email {
	return &Email{
		client:    sendgrid.NewSendClient(token),
		emailFrom: mail.NewEmail(mailName, address),
	}
}

func (s *Email) SendEmail(email, letter string) error {
	subject := "Письмо из прошлого"
	to := mail.NewEmail("", email)

	message := mail.NewSingleEmail(s.emailFrom, subject, to, "", letter)

	_, err := s.client.Send(message)
	if err != nil {
		return fmt.Errorf("sending mail to client: %w", err)
	}

	return nil
}
