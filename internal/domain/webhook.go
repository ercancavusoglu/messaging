package domain

type WebhookResponse struct {
	MessageID string `json:"messageId"`
	Message   string `json:"message"`
}
