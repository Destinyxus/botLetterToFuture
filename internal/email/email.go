package email

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"LetterToFuture/internal/store"
)

type Email struct {
	Store *store.Store
}

func NewEmail(store *store.Store) *Email {
	return &Email{
		Store: store,
	}
}

func (s *Email) SendEmail(email, letter string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	token := os.Getenv("SENDGRID_API_KEY")

	from := mail.NewEmail("botLetterToFuture", "vldmrbusiness@gmail.com")

	subject := "LetterFromPast"
	to := mail.NewEmail("Example User", email)
	plainTextContent := ""
	htmlContent := letter
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(token)
	response, err := client.Send(message)

	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}
