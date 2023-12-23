package email

//
//import (
//	"fmt"
//	"log"
//
//	"github.com/Destinyxus/botLetterToFuture/internal/config"
//	"github.com/sendgrid/sendgrid-go"
//	"github.com/sendgrid/sendgrid-go/helpers/mail"
//
//	"github.com/Destinyxus/botLetterToFuture/internal/store"
//)
//
//type Email struct {
//	Store *store.Store
//}
//
//func NewEmail(store *store.Store) *Email {
//	return &Email{
//		Store: store,
//	}
//}
//
//func (s *Email) SendEmail(email, letter string, cfg *config.Config) {
//
//	from := mail.NewEmail("botLetterToFuture", "lettertofuturebot@gmail.com")
//
//	subject := "Письмо из прошлого"
//	to := mail.NewEmail("", email)
//	plainTextContent := ""
//	htmlContent := letter
//	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
//	client := sendgrid.NewSendClient("fsaf")
//	response, err := client.Send(message)
//
//	if err != nil {
//		log.Println(err)
//	} else {
//		fmt.Println(response.StatusCode)
//		fmt.Println(response.Body)
//		fmt.Println(response.Headers)
//	}
//}
