package valueobject

import (
	"errors"
)

var (
	PhoneNumberLength           = 12
	ErrInvalidPhoneNumberLength = errors.New("invalid phone number length")
)

type PhoneNumber struct {
	number string
}

func NewPhoneNumber(number string) (*PhoneNumber, error) {
	if len(number) != PhoneNumberLength {
		return nil, ErrInvalidPhoneNumberLength
	}

	return &PhoneNumber{number: number}, nil
}

func (p PhoneNumber) String() string {
	return p.number
}
