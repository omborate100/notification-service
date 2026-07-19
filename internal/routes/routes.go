package routes

import (
	"net/http"

	"notification-service/internal/handler"
)

func RegisterRoutes(
	emailHandler *handler.EmailHandler,
) {

	http.HandleFunc(
		"/api/v1/email",
		emailHandler.SendEmail,
	)

	http.HandleFunc(
		"/health",
		func(w http.ResponseWriter, r *http.Request) {

			w.WriteHeader(http.StatusOK)

			w.Write([]byte("OK"))
		},
	)
}