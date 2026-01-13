
package notification

import (
	"fmt"
	"net/smtp"
	"os"
)

// SendResetEmail sends a reset password email with the provided token.
func SendResetEmail(email, token string) error {
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

	// TODO: Replace with SendGrid/AWS SES or proper SMTP
	err := sendSMTP(email, subject, body)
	if err != nil {
		return fmt.Errorf("failed to send reset email: %w", err)
	}
	return nil
}

// sendSMTP is a low-level helper (can be swapped with provider later).
func sendSMTP(to, subject, body string) error {
	smtpHost := os.Getenv("SMTP_HOST") // e.g., smtp.gmail.com
	smtpPort := os.Getenv("SMTP_PORT") // e.g., 587
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body))

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{to}, msg)
}
