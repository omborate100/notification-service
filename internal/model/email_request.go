package model

type EmailRequest struct {
	TemplateKey string                 `json:"template_key"`
	To          string                 `json:"to"`
	Variables   map[string]interface{} `json:"variables"`
}

type BulkRecipient struct {
	To        string                 `json:"to"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type BulkEmailRequest struct {
	TemplateKey string                 `json:"template_key"`
	Variables   map[string]interface{} `json:"variables,omitempty"`
	Recipients  []BulkRecipient        `json:"recipients"`
}

type BulkRecipientResult struct {
	To      string `json:"to"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type BulkEmailResponse struct {
	Total      int                   `json:"total"`
	Successful int                   `json:"successful"`
	Failed     int                   `json:"failed"`
	Results    []BulkRecipientResult `json:"results"`
}
