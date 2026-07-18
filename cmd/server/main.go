package main

import (
	"log"
	"net/http"

	"notification-service/config"
	"notification-service/internal/database"
)

func main() {

	cfg := config.Load()

	database.Connect(cfg.DatabaseURL)

	defer database.Close()

	http.HandleFunc("/health", healthHandler)

	log.Println("Server started on port :", cfg.AppPort)

	err := http.ListenAndServe(":"+cfg.AppPort, nil)

	if err != nil {
		log.Fatal(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	w.Write([]byte(`{"status":"UP"}`))
}