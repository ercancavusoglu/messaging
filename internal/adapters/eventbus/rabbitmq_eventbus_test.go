package eventbus

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/ercancavusoglu/messaging/internal/domain"
	"github.com/ercancavusoglu/messaging/internal/ports"
	"github.com/stretchr/testify/assert"
)

type MockEvent struct {
	name        string
	occurredOn  time.Time
	aggregateID string
	data        interface{}
}

func (e *MockEvent) EventName() string {
	return e.name
}

func (e *MockEvent) OccurredAt() time.Time {
	return e.occurredOn
}

func (e *MockEvent) GetAggregateID() string {
	return e.aggregateID
}

func TestRabbitMQEventBus_Subscribe(t *testing.T) {
	bus := &RabbitMQEventBus{
		handlers: make(map[string][]ports.EventHandler),
	}

	eventName := "test.event"
	handler := func(event ports.Event) error {
		return nil
	}

	bus.Subscribe(eventName, handler)

	assert.Len(t, bus.handlers[eventName], 1)
}

func TestRabbitMQEventBus_Unsubscribe(t *testing.T) {
	bus := &RabbitMQEventBus{
		handlers: make(map[string][]ports.EventHandler),
	}

	eventName := "test.event"
	handler := func(event ports.Event) error {
		return nil
	}

	bus.Subscribe(eventName, handler)
	assert.Len(t, bus.handlers[eventName], 1)

	bus.Unsubscribe(eventName, handler)
	assert.Len(t, bus.handlers[eventName], 0)
}

func TestRabbitMQEventBus_Integration(t *testing.T) {
	// RabbitMQ bağlantısı gerektiği için bu testi skip edelim
	t.Skip("Skipping integration test")

	bus, err := NewRabbitMQEventBus("amqp://guest:guest@localhost:5672/")
	if err != nil {
		t.Fatalf("Failed to create RabbitMQ event bus: %v", err)
	}
	defer bus.Close()

	eventName := "test.event"
	eventData := domain.Message{
		ID:        1,
		To:        "+905551234567",
		Content:   "Test message",
		Status:    domain.StatusPending,
		CreatedAt: time.Now(),
	}

	eventReceived := make(chan bool)
	handler := func(event ports.Event) error {
		var msg domain.Message
		data, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("failed to marshal event: %v", err)
		}

		if err := json.Unmarshal(data, &msg); err != nil {
			return fmt.Errorf("failed to unmarshal event data: %v", err)
		}

		assert.Equal(t, eventData.ID, msg.ID)
		assert.Equal(t, eventData.To, msg.To)
		assert.Equal(t, eventData.Content, msg.Content)
		assert.Equal(t, eventData.Status, msg.Status)

		eventReceived <- true
		return nil
	}

	bus.Subscribe(eventName, handler)

	mockEvent := &MockEvent{
		name:        eventName,
		occurredOn:  time.Now(),
		aggregateID: "1",
		data:        eventData,
	}

	err = bus.Publish(mockEvent)
	assert.NoError(t, err)

	select {
	case <-eventReceived:
		// Event başarıyla alındı
	case <-time.After(5 * time.Second):
		t.Error("Event timeout")
	}
}
