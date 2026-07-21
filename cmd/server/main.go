package main

import (
	"log"
	"net/http"

	"notification-service/config"
	"notification-service/internal/database"
	"notification-service/internal/handler"
	"notification-service/internal/renderer"
	"notification-service/internal/repository"
	"notification-service/internal/routes"
	"notification-service/internal/service"
	"notification-service/internal/smtp"
)

func main() {

	// Load configuration
	cfg := config.Load()
	
	// Connect Database
	database.Connect(cfg.DatabaseURL)
	defer database.Close()

	// Get DB Pool
	db := database.GetDB()

	// Repositories
	templateRepo := repository.NewTemplateRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)

	// Renderer
	templateRenderer := renderer.NewTemplateRenderer()

	// SMTP Sender
	smtpSender := smtp.NewSMTPSender(cfg)

	// Service
	emailService := service.NewEmailService(
		templateRepo,
		notificationRepo,
		templateRenderer,
		smtpSender,
	)

	// Handler
	emailHandler := handler.NewEmailHandler(emailService)

	// Router
	mux := http.NewServeMux()

	routes.RegisterRoutes(
		mux,
		emailHandler,
	)

	log.Printf("Server started on port %s", cfg.AppPort)

	err := http.ListenAndServe(
		":"+cfg.AppPort,
		mux,
	)

	if err != nil {
		log.Fatal(err)
	}
}