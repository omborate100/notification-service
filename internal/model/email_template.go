package model

import "time"

type EmailTemplate struct {
	ID              int64
	TemplateKey     string
	TemplateName    string
	SubjectTemplate string
	BodyTemplate    string
	Description     *string
	IsActive        bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}