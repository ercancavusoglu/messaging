package eventbus

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ercancavusoglu/messaging/internal/domain/event"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	exchangeName = "events"
	queueName    = "events.queue"
)

type RabbitMQEventBus struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	handlers map[string][]event.EventHandler
}

func NewRabbitMQEventBus(url string) (*RabbitMQEventBus, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %v", err)
	}

	err = ch.ExchangeDeclare(
		exchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %v", err)
	}

	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %v", err)
	}

	err = ch.QueueBind(
		queueName,
		"#",
		exchangeName,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to bind queue: %v", err)
	}

	bus := &RabbitMQEventBus{
		conn:     conn,
		channel:  ch,
		handlers: make(map[string][]event.EventHandler),
	}

	go bus.startConsumer()

	return bus, nil
}

func (b *RabbitMQEventBus) startConsumer() {
	msgs, err := b.channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Printf("[RabbitMQ] Failed to register consumer: %v\n", err)
		return
	}

	for d := range msgs {
		fmt.Printf("[RabbitMQ] Received message: %s\n", string(d.Body))

		var env event.EventEnvelope
		if err := json.Unmarshal(d.Body, &env); err != nil {
			fmt.Printf("[RabbitMQ] Failed to unmarshal event envelope: %v\n", err)
			d.Reject(true)
			continue
		}

		fmt.Printf("[RabbitMQ] Processing event: %s, AggregateID: %s\n", env.Name, env.AggregateID)

		handlers := b.handlers[env.Name]
		if len(handlers) == 0 {
			fmt.Printf("[RabbitMQ] No handlers found for event: %s\n", env.Name)
			d.Ack(false)
			continue
		}

		var handlerError error
		for _, handler := range handlers {
			if err := handler(&env); err != nil {
				fmt.Printf("[RabbitMQ] Handler failed for event %s: %v\n", env.Name, err)
				handlerError = err
				break
			}
		}

		if handlerError != nil {
			d.Reject(true)
			continue
		}

		d.Ack(false)
		fmt.Printf("[RabbitMQ] Successfully processed event: %s\n", env.Name)
	}
}

func (b *RabbitMQEventBus) Publish(event event.Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %v", err)
	}

	fmt.Printf("[RabbitMQ] Event data: %s\n", string(data))

	env := struct {
		Name        string          `json:"name"`
		OccurredOn  time.Time       `json:"occurred_on"`
		AggregateID string          `json:"aggregate_id"`
		Data        json.RawMessage `json:"data"`
	}{
		Name:        event.EventName(),
		OccurredOn:  event.OccurredAt(),
		AggregateID: event.GetAggregateID(),
		Data:        data,
	}

	body, err := json.Marshal(env)
	if err != nil {
		return fmt.Errorf("failed to marshal envelope: %v", err)
	}

	fmt.Printf("[RabbitMQ] Publishing event: %s\n", string(body))

	err = b.channel.Publish(
		exchangeName,
		event.EventName(),
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish event: %v", err)
	}

	fmt.Printf("[RabbitMQ] Successfully published event: %s\n", event.EventName())
	return nil
}

func (b *RabbitMQEventBus) Subscribe(eventName string, handler event.EventHandler) {
	b.handlers[eventName] = append(b.handlers[eventName], handler)
}

func (b *RabbitMQEventBus) Unsubscribe(eventName string, handler event.EventHandler) {
	handlers := b.handlers[eventName]
	for i, h := range handlers {
		if &h == &handler {
			b.handlers[eventName] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}
}

func (b *RabbitMQEventBus) Close() error {
	if err := b.channel.Close(); err != nil {
		return fmt.Errorf("failed to close channel: %v", err)
	}
	if err := b.conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %v", err)
	}
	return nil
}
