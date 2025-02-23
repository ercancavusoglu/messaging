package eventbus

import (
	"sync"

	"github.com/ercancavusoglu/messaging/internal/domain/event"
)

type EventHandler func(event event.Event) error

type EventBus interface {
	Publish(event event.Event) error
	Subscribe(eventName string, handler EventHandler)
	Unsubscribe(eventName string, handler EventHandler)
}

type InMemoryEventBus struct {
	handlers map[string][]EventHandler
	mu       sync.RWMutex
}

func NewInMemoryEventBus() *InMemoryEventBus {
	return &InMemoryEventBus{
		handlers: make(map[string][]EventHandler),
	}
}

func (b *InMemoryEventBus) Publish(event event.Event) error {
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

func (b *InMemoryEventBus) Subscribe(eventName string, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	handlers := b.handlers[eventName]
	handlers = append(handlers, handler)
	b.handlers[eventName] = handlers
}

func (b *InMemoryEventBus) Unsubscribe(eventName string, handler EventHandler) {
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
