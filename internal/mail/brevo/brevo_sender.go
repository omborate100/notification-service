package brevo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/mail"

	"notification-service/config"
)

type Sender struct {
	apiKey string

	fromName  string
	fromEmail string

	replyToName  string
	replyToEmail string
}

func NewSender(cfg *config.Config) *Sender {

	log.Printf("Using Brevo as email sender")

	fromAddr, err := mail.ParseAddress(cfg.FromEmail)
	if err != nil {
		panic(err)
	}

	replyAddr, err := mail.ParseAddress(cfg.ReplyToEmail)
	if err != nil {
		panic(err)
	}

	return &Sender{
		apiKey: cfg.BrevoAPIKey,

		fromName:  fromAddr.Name,
		fromEmail: fromAddr.Address,

		replyToName:  replyAddr.Name,
		replyToEmail: replyAddr.Address,
	}
}

type request struct {
	Sender sender `json:"sender"`

	ReplyTo sender `json:"replyTo,omitempty"`

	To []receiver `json:"to"`

	Subject string `json:"subject"`

	HTMLContent string `json:"htmlContent"`
}

type sender struct {
	Name string `json:"name"`

	Email string `json:"email"`
}

type receiver struct {
	Email string `json:"email"`
}

func (s *Sender) Send(
	to string,
	subject string,
	body string,
) error {

	reqBody := request{

		Sender: sender{
			Name:  s.fromName,
			Email: s.fromEmail,
		},

		ReplyTo: sender{
			Name:  s.replyToName,
			Email: s.replyToEmail,
		},

		Subject:     subject,
		HTMLContent: body,
	}

	reqBody.To = append(
		reqBody.To,
		receiver{
			Email: to,
		},
	)

	data, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.brevo.com/v3/smtp/email",
		bytes.NewBuffer(data),
	)

	if err != nil {
		return err
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("api-key", s.apiKey)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("brevo failed with status %d", resp.StatusCode)
	}

	return nil
}