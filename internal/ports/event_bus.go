package ports

import (
	"time"
)

type EventHandler func(event Event) error

type Event interface {
	EventName() string
	OccurredAt() time.Time
	GetAggregateID() string
}

type EventBus interface {
	Publish(event Event) error
	Subscribe(eventName string, handler EventHandler)
	Unsubscribe(eventName string, handler EventHandler)
}
