package notification


// Mailer interface
type Mailer interface {
	SendResetEmail(to, token string) error
	SendEmail(to, subject, body string) error
}

// SendGridMailer implements Mailer






