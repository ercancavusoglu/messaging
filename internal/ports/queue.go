package ports

type MessageQueue interface {
	Publish(queueName string, message []byte) error
	Consume(queueName string, handler func([]byte) error) error
	Close() error
}
