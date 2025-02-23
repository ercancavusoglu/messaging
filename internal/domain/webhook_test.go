package domain

import (
	"encoding/json"
	"testing"
)

func TestWebhookResponse_JSON(t *testing.T) {
	response := WebhookResponse{
		MessageID: "msg_123",
		Message:   "Message sent successfully",
	}

	// JSON marshal testi
	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Errorf("Failed to marshal WebhookResponse: %v", err)
	}

	expectedJSON := `{"messageId":"msg_123","message":"Message sent successfully"}`
	if string(jsonData) != expectedJSON {
		t.Errorf("Expected JSON to be %s, got %s", expectedJSON, string(jsonData))
	}

	// JSON unmarshal testi
	var unmarshaledResponse WebhookResponse
	err = json.Unmarshal(jsonData, &unmarshaledResponse)
	if err != nil {
		t.Errorf("Failed to unmarshal WebhookResponse: %v", err)
	}

	if unmarshaledResponse.MessageID != response.MessageID {
		t.Errorf("Expected MessageID to be %s, got %s", response.MessageID, unmarshaledResponse.MessageID)
	}

	if unmarshaledResponse.Message != response.Message {
		t.Errorf("Expected Message to be %s, got %s", response.Message, unmarshaledResponse.Message)
	}
}
