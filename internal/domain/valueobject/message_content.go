package valueobject

import (
	"errors"
	"fmt"
)

type MessageContent struct {
	value string
}

var (
	MessageContentMaxLength = 20
	ErrEmptyContent         = errors.New("message content cannot be empty")
	ErrContentTooLong       = fmt.Errorf("message content cannot be longer than %d characters", MessageContentMaxLength)
)

func NewMessageContent(content string) (*MessageContent, error) {
	if content == "" {
		return nil, ErrEmptyContent
	}

	if len(content) >= MessageContentMaxLength {
		return nil, ErrContentTooLong
	}

	return &MessageContent{
		value: content,
	}, nil
}

func (mc *MessageContent) String() string {
	return mc.value
}
