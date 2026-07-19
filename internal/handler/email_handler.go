package handler

import (
	"encoding/json"
	"net/http"

	"notification-service/internal/model"
	"notification-service/internal/service"
)

type EmailHandler struct {
	emailService *service.EmailService
}

func NewEmailHandler(emailService *service.EmailService) *EmailHandler {
	return &EmailHandler{
		emailService: emailService,
	}
}

func (h *EmailHandler) SendEmail(
	w http.ResponseWriter,
	r *http.Request,
) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var request model.EmailRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {

		http.Error(
			w,
			"Invalid request payload",
			http.StatusBadRequest,
		)

		return
	}

	if request.TemplateKey == "" {
		http.Error(
			w,
			"template_key is required",
			http.StatusBadRequest,
		)
		return
	}

	if request.To == "" {
		http.Error(
			w,
			"to is required",
			http.StatusBadRequest,
		)
		return
	}

	if request.Variables == nil {
		request.Variables = map[string]interface{}{}
	}

	err = h.emailService.SendEmail(
		r.Context(),
		&request,
	)

	if err != nil {

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Email sent successfully",
	})
}