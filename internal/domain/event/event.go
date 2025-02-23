package event

import (
	"encoding/json"
	"time"
)

type Event interface {
	EventName() string
	OccurredAt() time.Time
	GetAggregateID() string
}

type EventHandler func(event Event) error

type EventBus interface {
	Publish(event Event) error
	Subscribe(eventName string, handler EventHandler)
	Unsubscribe(eventName string, handler EventHandler)
}

type BaseEvent struct {
	Name        string    `json:"name"`
	OccurredOn  time.Time `json:"occurred_on"`
	AggregateID string    `json:"aggregate_id"`
}

func (e BaseEvent) EventName() string {
	return e.Name
}

func (e BaseEvent) OccurredAt() time.Time {
	return e.OccurredOn
}

func (e BaseEvent) GetAggregateID() string {
	return e.AggregateID
}

func NewBaseEvent(name string, aggregateID string) BaseEvent {
	return BaseEvent{
		Name:        name,
		OccurredOn:  time.Now(),
		AggregateID: aggregateID,
	}
}

type EventEnvelope struct {
	Name        string          `json:"name"`
	OccurredOn  time.Time       `json:"occurred_on"`
	AggregateID string          `json:"aggregate_id"`
	Data        json.RawMessage `json:"data"`
}

func (e *EventEnvelope) EventName() string {
	return e.Name
}

func (e *EventEnvelope) OccurredAt() time.Time {
	return e.OccurredOn
}

func (e *EventEnvelope) GetAggregateID() string {
	return e.AggregateID
}
