package mail

import (
	"notification-service/config"
	"notification-service/internal/mail/brevo"
	"notification-service/internal/mail/smtp"
)

func NewSender(cfg *config.Config) Sender {

	if cfg.AppEnv == "PROD" {
		return brevo.NewSender(cfg)
	}

	return smtp.NewSMTPSender(cfg)
}