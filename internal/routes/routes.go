package routes

import (
	"net/http"

	"notification-service/internal/handler"
)

func RegisterRoutes(
	mux *http.ServeMux,
	emailHandler *handler.EmailHandler,
) {

	mux.HandleFunc(
		"/api/v1/email",
		emailHandler.SendEmail,
	)

	mux.HandleFunc(
		"/health",
		func(w http.ResponseWriter, r *http.Request) {

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"UP"}`))
		},
	)
}