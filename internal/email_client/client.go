package emailclient

import (
	"net/smtp"
	"strings"

	commander "github.com/Destinyxus/botLetterToFuture/internal/bot_commander"
)

type emailClient struct {
	auth        smtp.Auth
	clientEmail string
	smtpAddr    string
}

func New(token, clientEmail, host, smtpAddr string) commander.EmailClient {
	return &emailClient{
		auth:        smtp.PlainAuth("", clientEmail, strings.Join(splitIntoChunks(token, 4), " "), host),
		clientEmail: clientEmail,
		smtpAddr:    smtpAddr,
	}
}

func (s *emailClient) SendEmail(email, letter string) error {
	if err := smtp.SendMail(s.smtpAddr,
		s.auth,
		s.clientEmail,
		[]string{email},
		[]byte(letter)); err != nil {

		return err
	}

	return nil
}

func splitIntoChunks(s string, chunkSize int) []string {
	var chunks []string

	for i := 0; i < len(s); i += chunkSize {
		end := i + chunkSize
		if end > len(s) {
			end = len(s)
		}

		chunks = append(chunks, s[i:end])
	}

	return chunks
}
