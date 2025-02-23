package messaging

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testMessage struct {
	Content string `json:"content"`
}

func TestNewRabbitMQ_InvalidURL(t *testing.T) {
	// Geçersiz URL ile bağlantı denemesi
	_, err := NewRabbitMQ("amqp://invalid:5672")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to connect to RabbitMQ")
}

func TestRabbitMQ_Integration(t *testing.T) {
	// RabbitMQ bağlantısı gerektiği için bu testi skip edelim
	t.Skip("Skipping integration test")

	// RabbitMQ bağlantısı
	rabbit, err := NewRabbitMQ("amqp://guest:guest@localhost:5672/")
	assert.NoError(t, err)
	defer rabbit.Close()

	// Test mesajı
	testMsg := testMessage{Content: "test message"}
	msgBytes, err := json.Marshal(testMsg)
	assert.NoError(t, err)

	// Test kuyruğu
	queueName := "test_queue"

	// Mesaj alındı kanalı
	messageReceived := make(chan bool)

	// Consumer handler
	handler := func(msg []byte) error {
		var receivedMsg testMessage
		err := json.Unmarshal(msg, &receivedMsg)
		assert.NoError(t, err)
		assert.Equal(t, testMsg.Content, receivedMsg.Content)
		messageReceived <- true
		return nil
	}

	// Consumer başlat
	err = rabbit.Consume(queueName, handler)
	assert.NoError(t, err)

	// Mesaj gönder
	err = rabbit.Publish(queueName, msgBytes)
	assert.NoError(t, err)

	// Mesajın alınmasını bekle
	select {
	case <-messageReceived:
		// Mesaj başarıyla alındı
	case <-time.After(5 * time.Second):
		t.Error("Message not received within timeout")
	}
}

func TestRabbitMQ_Close(t *testing.T) {
	// RabbitMQ bağlantısı gerektiği için bu testi skip edelim
	t.Skip("Skipping close test")

	rabbit, err := NewRabbitMQ("amqp://guest:guest@localhost:5672/")
	assert.NoError(t, err)

	// Bağlantıyı kapat
	err = rabbit.Close()
	assert.NoError(t, err)

	// Kapalı bağlantı üzerinden işlem yapmayı dene
	err = rabbit.Publish("test_queue", []byte("test"))
	assert.Error(t, err)
}
