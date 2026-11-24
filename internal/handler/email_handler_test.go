package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"notification-service/internal/handler"
	"notification-service/internal/model"
)

func TestSendBulkEmail_Validation(t *testing.T) {
	h := handler.NewEmailHandler(nil)

	// Test Method Not Allowed
	req := httptest.NewRequest(http.MethodGet, "/api/v1/email/bulk", nil)
	w := httptest.NewRecorder()
	h.SendBulkEmail(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}

	// Test Missing template_key
	invalidBody, _ := json.Marshal(model.BulkEmailRequest{
		TemplateKey: "",
		Recipients: []model.BulkRecipient{
			{To: "test@example.com"},
		},
	})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/email/bulk", bytes.NewBuffer(invalidBody))
	w = httptest.NewRecorder()
	h.SendBulkEmail(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d for empty template_key, got %d", http.StatusBadRequest, w.Code)
	}

	// Test Missing recipients
	noRecipientsBody, _ := json.Marshal(model.BulkEmailRequest{
		TemplateKey: "WELCOME",
		Recipients:  []model.BulkRecipient{},
	})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/email/bulk", bytes.NewBuffer(noRecipientsBody))
	w = httptest.NewRecorder()
	h.SendBulkEmail(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d for empty recipients, got %d", http.StatusBadRequest, w.Code)
	}
}
