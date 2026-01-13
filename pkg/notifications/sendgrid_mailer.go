// internal/notification/sendgrid_mailer.go
package notification

import (
	"fmt"
	// "os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	client *sendgrid.Client
	from   string
	
	frontendURL string

}

func NewSendGridMailer(apiKey, from, frontendURL string) *SendGridMailer {
	return &SendGridMailer{
		client:      sendgrid.NewSendClient(apiKey),
		from:        from,
		frontendURL: frontendURL,
	}
}


func (m *SendGridMailer) SendEmail(to, subject, body string) error {
	from := mail.NewEmail("ClaimEasy", m.from)
	toEmail := mail.NewEmail("", to)
	message := mail.NewSingleEmail(from, subject, toEmail, body, body)
	_, err := m.client.Send(message)
	return err
}

func (m *SendGridMailer) SendResetEmail(to, token string) error {
	link := fmt.Sprintf("%s/reset-password?token=%s", m.frontendURL, token)
	subject := "Password Reset Request"
	body := fmt.Sprintf("Click to reset: %s", link)
	return m.SendEmail(to, subject, body)
}
