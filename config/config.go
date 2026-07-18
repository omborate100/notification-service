package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	DatabaseURL string

	SMTPHost string
	SMTPPort string

	SMTPUsername string
	SMTPPassword string

	FromEmail string
}

func Load() *Config {

	err := godotenv.Load()

	if err != nil {
		log.Println(".env file not found. Using system environment variables.")
	}

	return &Config{
		AppPort: os.Getenv("APP_PORT"),

		DatabaseURL: os.Getenv("DATABASE_URL"),

		SMTPHost: os.Getenv("SMTP_HOST"),
		SMTPPort: os.Getenv("SMTP_PORT"),

		SMTPUsername: os.Getenv("SMTP_USERNAME"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),

		FromEmail: os.Getenv("FROM_EMAIL"),
	}
}