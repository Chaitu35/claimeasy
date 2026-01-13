// internal/pkg/notification/smtp_mailer.go
package notification

import (
	"fmt"
	"net/smtp"
	"os"
)

type SMTPMailer struct {
	host string
	port string
	user string
	pass string
	from string
}

func NewSMTPMailer() *SMTPMailer {
	return &SMTPMailer{
		host: os.Getenv("SMTP_HOST"), // e.g. smtp.gmail.com
		port: os.Getenv("SMTP_PORT"), // e.g. 587
		user: os.Getenv("SMTP_USER"),
		pass: os.Getenv("SMTP_PASS"),
		from: os.Getenv("SMTP_FROM"),
	}
}

func (m *SMTPMailer) SendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", m.user, m.pass, m.host)

	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body))

	return smtp.SendMail(m.host+":"+m.port, auth, m.from, []string{to}, msg)
}

func (m *SMTPMailer) SendResetEmail(to, token string) error {
	resetLink := fmt.Sprintf("%s/reset-password?token=%s",
		os.Getenv("FRONTEND_URL"), token)

	subject := "Password Reset Request"
	body := fmt.Sprintf(`
Hello,

We received a request to reset your password.
Click the link below to reset it:

%s

If you did not request this, please ignore this email.

Thanks,
The ClaimEasy Team
`, resetLink)

	return m.SendEmail(to, subject, body)
}
