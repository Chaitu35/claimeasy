package notification

import "github.com/chaitu35/claimeasy/configs"

// ProvideMailer chooses the appropriate Mailer implementation
func ProvideMailer(cfg *config.Config) Mailer {
	if cfg.SendGridAPIKey != "" && cfg.SendGridFrom != "" && cfg.FrontendURL != "" {
		// Production or staging → SendGrid
		return NewSendGridMailer(cfg.SendGridAPIKey, cfg.SendGridFrom, cfg.FrontendURL)
	}
	// Development → ConsoleMailer
	return &ConsoleMailer{}
}
