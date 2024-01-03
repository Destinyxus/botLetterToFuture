package emailSender

import (
	"net/smtp"
)

type EmailClient struct {
	auth        smtp.Auth
	clientEmail string
	smtpAddr    string
}

func New(token, clientEmail, host, smtpAddr string) *EmailClient {
	return &EmailClient{
		auth:        smtp.PlainAuth("", clientEmail, token, host),
		clientEmail: clientEmail,
		smtpAddr:    smtpAddr,
	}
}

func (s *EmailClient) SendEmail(email, letter string) error {
	if err := smtp.SendMail(s.smtpAddr,
		s.auth,
		s.clientEmail,
		[]string{email},
		[]byte(letter)); err != nil {

		return err
	}

	return nil
}
