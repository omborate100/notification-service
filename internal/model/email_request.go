package model

type EmailRequest struct {
	TemplateKey string                 `json:"template_key"`
	To          string                 `json:"to"`
	Variables   map[string]interface{} `json:"variables"`
}