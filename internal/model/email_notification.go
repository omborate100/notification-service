package model

import "time"

type EmailNotification struct {
	ID                int64
	TemplateID        int64
	RecipientEmail    string
	Subject           string
	Body              string
	Variables         []byte
	Status            string
	Provider          string
	ProviderMessageID *string
	ErrorMessage      *string
	SentAt            *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}