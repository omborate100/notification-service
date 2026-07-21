package service

import (
	"context"
	"fmt"
	"log"
	"notification-service/internal/model"
	"notification-service/internal/renderer"
	"notification-service/internal/repository"
	"notification-service/internal/smtp"
)

type EmailService struct {
	templateRepo     *repository.TemplateRepository
	notificationRepo *repository.NotificationRepository
	renderer         *renderer.TemplateRenderer
	smtpSender       *smtp.SMTPSender
}

func NewEmailService(
	templateRepo *repository.TemplateRepository,
	notificationRepo *repository.NotificationRepository,
	renderer *renderer.TemplateRenderer,
	smtpSender *smtp.SMTPSender,
) *EmailService {

	return &EmailService{
		templateRepo:     templateRepo,
		notificationRepo: notificationRepo,
		renderer:         renderer,
		smtpSender:       smtpSender,
	}
}

func (s *EmailService) SendEmail(
	ctx context.Context,
	req *model.EmailRequest,
) error {
	log.Println("Starting to send email in service")
	// Fetch template
	emailTemplate, err := s.templateRepo.GetByTemplateKey(
		ctx,
		req.TemplateKey,
	)

	if err != nil {
		return err
	}

	// Render subject
	subject, err := s.renderer.Render(
		emailTemplate.SubjectTemplate,
		req.Variables,
	)

	if err != nil {
		return fmt.Errorf("failed to render subject: %w", err)
	}

	// Render body
	body, err := s.renderer.Render(
		emailTemplate.BodyTemplate,
		req.Variables,
	)

	if err != nil {
		return fmt.Errorf("failed to render body: %w", err)
	}

	notification := &model.EmailNotification{
		TemplateID:     emailTemplate.ID,
		RecipientEmail: req.To,
		Subject:        subject,
		Body:           body,
		Variables:      req.Variables,
		Status:         model.StatusPending,
		Provider:       model.ProviderSMTP,
	}

	// Create notification entry
	notificationID, err := s.notificationRepo.Create(
		ctx,
		notification,
	)
	log.Println("Created notification with ID:", notificationID)
	if err != nil {
		return err
	}

	// Send email
	err = s.smtpSender.Send(
		req.To,
		subject,
		body,
	)
	log.Println("Sent email to:", req.To)

	if err != nil {

		// Update notification as FAILED
		_ = s.notificationRepo.MarkFailed(
			ctx,
			notificationID,
			err.Error(),
		)

		return err
	}

	// Update notification as SENT
	err = s.notificationRepo.MarkSent(
		ctx,
		notificationID,
		"",
	)

	if err != nil {
		return err
	}

	return nil
}