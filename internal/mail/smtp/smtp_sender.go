package smtp

import (
	"notification-service/config"

	"gopkg.in/mail.v2"
	"log"
)

type SMTPSender struct {
	host      string
	port      int
	username  string
	password  string
	fromEmail string
	ReplyToEmail string
}

func NewSMTPSender(cfg *config.Config) *SMTPSender {
	log.Printf("Using SMTP as email sender")
	port := 587

	if cfg.SMTPPort == "465" {
		port = 465
	}

	return &SMTPSender{
		host:      cfg.SMTPHost,
		port:      port,
		username:  cfg.SMTPUsername,
		password:  cfg.SMTPPassword,
		fromEmail: cfg.FromEmail,
		ReplyToEmail: cfg.ReplyToEmail,
	}
}

func (s *SMTPSender) Send(
	to string,
	subject string,
	body string,
) error {

	message := mail.NewMessage()

	message.SetHeader("From", s.fromEmail)
	message.SetHeader("To", to)
	message.SetHeader("Reply-To", s.ReplyToEmail)
	message.SetHeader("Subject", subject)

	message.SetBody("text/html", body)

	dialer := mail.NewDialer(
		s.host,
		s.port,
		s.username,
		s.password,
	)

	return dialer.DialAndSend(message)
}