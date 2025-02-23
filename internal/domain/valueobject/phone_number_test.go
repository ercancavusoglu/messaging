package valueobject

import (
	"testing"
)

func TestNewPhoneNumber(t *testing.T) {
	tests := []struct {
		name    string
		number  string
		wantErr error
	}{
		{
			name:    "Valid phone number",
			number:  "+90555123456",
			wantErr: nil,
		},
		{
			name:    "Invalid length - too short",
			number:  "+90555",
			wantErr: ErrInvalidPhoneNumberLength,
		},
		{
			name:    "Invalid length - too long",
			number:  "+905551234567890",
			wantErr: ErrInvalidPhoneNumberLength,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pn, err := NewPhoneNumber(tt.number)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("NewPhoneNumber() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewPhoneNumber() unexpected error = %v", err)
				return
			}

			if pn.String() != tt.number {
				t.Errorf("PhoneNumber.String() = %v, want %v", pn.String(), tt.number)
			}
		})
	}
}

func TestPhoneNumber_String(t *testing.T) {
	number := "+90555123456"
	pn, err := NewPhoneNumber(number)
	if err != nil {
		t.Fatalf("Failed to create PhoneNumber: %v", err)
	}

	if pn.String() != number {
		t.Errorf("PhoneNumber.String() = %v, want %v", pn.String(), number)
	}
}
