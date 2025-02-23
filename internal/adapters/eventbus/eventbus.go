package eventbus

import (
	"github.com/ercancavusoglu/messaging/internal/ports"
	"sync"
)

type InMemoryEventBus struct {
	handlers map[string][]ports.EventHandler
	mu       sync.RWMutex
}

func NewInMemoryEventBus() *InMemoryEventBus {
	return &InMemoryEventBus{
		handlers: make(map[string][]ports.EventHandler),
	}
}

func (b *InMemoryEventBus) Publish(event ports.Event) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	handlers, exists := b.handlers[event.EventName()]
	if !exists {
		return nil
	}

	for _, handler := range handlers {
		if err := handler(event); err != nil {
			return err
		}
	}

	return nil
}

func (b *InMemoryEventBus) Subscribe(eventName string, handler ports.EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	handlers := b.handlers[eventName]
	handlers = append(handlers, handler)
	b.handlers[eventName] = handlers
}

func (b *InMemoryEventBus) Unsubscribe(eventName string, handler ports.EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	handlers := b.handlers[eventName]
	for i, h := range handlers {
		if &h == &handler {
			b.handlers[eventName] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}
}
