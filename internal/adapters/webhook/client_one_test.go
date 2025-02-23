package webhook

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ercancavusoglu/messaging/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestClient_SendMessage_Success(t *testing.T) {
	// Test sunucusu oluştur
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Request kontrolü
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer test-api-key", r.Header.Get("Authorization"))

		// Request body kontrolü
		var requestBody map[string]string
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		assert.NoError(t, err)
		assert.Equal(t, "+905551234567", requestBody["to"])
		assert.Equal(t, "Test message", requestBody["content"])

		// Response
		response := domain.WebhookResponse{
			MessageID: "msg_123",
			Message:   "Message sent successfully",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Client oluştur
	client := NewClient(server.URL, "test-api-key")

	// Test
	response, err := client.SendMessage("+905551234567", "Test message")
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "msg_123", response.MessageID)
	assert.Equal(t, "Message sent successfully", response.Message)
}

func TestClient_SendMessage_Error(t *testing.T) {
	// Test sunucusu oluştur
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	// Client oluştur
	client := NewClient(server.URL, "test-api-key")

	// Test
	response, err := client.SendMessage("+905551234567", "Test message")
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "unexpected status code: 500")
}

func TestClient_SendMessage_InvalidResponse(t *testing.T) {
	// Test sunucusu oluştur
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	// Client oluştur
	client := NewClient(server.URL, "test-api-key")

	// Test
	response, err := client.SendMessage("+905551234567", "Test message")
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "failed to decode response")
}
