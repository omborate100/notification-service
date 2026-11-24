package service

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"notification-service/internal/model"
	"notification-service/internal/renderer"
	"notification-service/internal/repository"
	"notification-service/internal/mail"
)

type EmailService struct {
	templateRepo     *repository.TemplateRepository
	notificationRepo *repository.NotificationRepository
	renderer         *renderer.TemplateRenderer
	sender           mail.Sender
}

func NewEmailService(
	templateRepo *repository.TemplateRepository,
	notificationRepo *repository.NotificationRepository,
	renderer *renderer.TemplateRenderer,
	sender mail.Sender,
) *EmailService {

	return &EmailService{
		templateRepo:     templateRepo,
		notificationRepo: notificationRepo,
		renderer:         renderer,
		sender:           sender,
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
	err = s.sender.Send(
		req.To,
		subject,
		body,
	)

	if err != nil {
		log.Printf("[ERROR] Failed to send email to %s: %v\n", req.To, err)

		// Update notification as FAILED
		_ = s.notificationRepo.MarkFailed(
			ctx,
			notificationID,
			err.Error(),
		)

		return err
	}

	log.Println("Sent email to:", req.To)

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

func (s *EmailService) SendBulkEmail(
	ctx context.Context,
	req *model.BulkEmailRequest,
) (*model.BulkEmailResponse, error) {
	log.Printf("Starting bulk email send for template '%s' to %d recipients\n", req.TemplateKey, len(req.Recipients))

	// Fetch template once for all recipients
	emailTemplate, err := s.templateRepo.GetByTemplateKey(
		ctx,
		req.TemplateKey,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch template '%s': %w", req.TemplateKey, err)
	}

	total := len(req.Recipients)
	results := make([]model.BulkRecipientResult, total)

	concurrency := 5
	if total < concurrency {
		concurrency = total
	}
	if concurrency <= 0 {
		concurrency = 1
	}

	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup
	var mu sync.Mutex

	successfulCount := 0
	failedCount := 0

	for i, recipient := range req.Recipients {
		wg.Add(1)
		sem <- struct{}{}

		go func(idx int, r model.BulkRecipient) {
			defer wg.Done()
			defer func() { <-sem }()

			trimmedTo := strings.TrimSpace(r.To)
			if trimmedTo == "" {
				mu.Lock()
				failedCount++
				results[idx] = model.BulkRecipientResult{
					To:      r.To,
					Success: false,
					Error:   "recipient email is empty",
				}
				mu.Unlock()
				return
			}

			// Merge common variables with recipient variables
			mergedVars := make(map[string]interface{})
			for k, v := range req.Variables {
				mergedVars[k] = v
			}
			for k, v := range r.Variables {
				mergedVars[k] = v
			}

			// Render subject
			subject, err := s.renderer.Render(
				emailTemplate.SubjectTemplate,
				mergedVars,
			)
			if err != nil {
				mu.Lock()
				failedCount++
				results[idx] = model.BulkRecipientResult{
					To:      trimmedTo,
					Success: false,
					Error:   fmt.Sprintf("failed to render subject: %v", err),
				}
				mu.Unlock()
				return
			}

			// Render body
			body, err := s.renderer.Render(
				emailTemplate.BodyTemplate,
				mergedVars,
			)
			if err != nil {
				mu.Lock()
				failedCount++
				results[idx] = model.BulkRecipientResult{
					To:      trimmedTo,
					Success: false,
					Error:   fmt.Sprintf("failed to render body: %v", err),
				}
				mu.Unlock()
				return
			}

			notification := &model.EmailNotification{
				TemplateID:     emailTemplate.ID,
				RecipientEmail: trimmedTo,
				Subject:        subject,
				Body:           body,
				Variables:      mergedVars,
				Status:         model.StatusPending,
				Provider:       model.ProviderSMTP,
			}

			notificationID, err := s.notificationRepo.Create(
				ctx,
				notification,
			)
			if err != nil {
				log.Printf("[WARN] Failed to create notification log for %s: %v\n", trimmedTo, err)
			}

			// Send email
			sendErr := s.sender.Send(
				trimmedTo,
				subject,
				body,
			)
			if sendErr != nil {
				log.Printf("[ERROR] Bulk email failed to send to %s: %v\n", trimmedTo, sendErr)
				if notificationID > 0 {
					_ = s.notificationRepo.MarkFailed(
						ctx,
						notificationID,
						sendErr.Error(),
					)
				}
				mu.Lock()
				failedCount++
				results[idx] = model.BulkRecipientResult{
					To:      trimmedTo,
					Success: false,
					Error:   sendErr.Error(),
				}
				mu.Unlock()
				return
			}

			if notificationID > 0 {
				_ = s.notificationRepo.MarkSent(
					ctx,
					notificationID,
					"",
				)
			}

			mu.Lock()
			successfulCount++
			results[idx] = model.BulkRecipientResult{
				To:      trimmedTo,
				Success: true,
			}
			mu.Unlock()
		}(i, recipient)
	}

	wg.Wait()

	log.Printf("Bulk email send completed: %d total, %d succeeded, %d failed\n", total, successfulCount, failedCount)

	return &model.BulkEmailResponse{
		Total:      total,
		Successful: successfulCount,
		Failed:     failedCount,
		Results:    results,
	}, nil
}