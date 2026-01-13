package notification

import "fmt"

type ConsoleMailer struct{}

func (c *ConsoleMailer) SendResetEmail(to, token string) error {
	fmt.Printf("[ConsoleMailer] Sending reset link to %s with token %s\n", to, token)
	return nil
}

func (c *ConsoleMailer) SendEmail(to, subject, body string) error {
	fmt.Printf("[ConsoleMailer] Email to: %s\nSubject: %s\nBody: %s\n", to, subject, body)
	return nil
}
