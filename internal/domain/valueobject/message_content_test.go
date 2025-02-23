package valueobject

import (
	"testing"
)

func TestNewMessageContent(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		wantErr     error
		wantContent string
	}{
		{
			name:        "Valid content",
			content:     "Hello World",
			wantErr:     nil,
			wantContent: "Hello World",
		},
		{
			name:        "Empty content",
			content:     "",
			wantErr:     ErrEmptyContent,
			wantContent: "",
		},
		{
			name:        "Content too long",
			content:     "This is a very long message that exceeds the maximum length",
			wantErr:     ErrContentTooLong,
			wantContent: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc, err := NewMessageContent(tt.content)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("NewMessageContent() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewMessageContent() unexpected error = %v", err)
				return
			}

			if mc.String() != tt.wantContent {
				t.Errorf("MessageContent.String() = %v, want %v", mc.String(), tt.wantContent)
			}
		})
	}
}

func TestMessageContent_String(t *testing.T) {
	content := "Test message"
	mc, err := NewMessageContent(content)
	if err != nil {
		t.Fatalf("Failed to create MessageContent: %v", err)
	}

	if mc.String() != content {
		t.Errorf("MessageContent.String() = %v, want %v", mc.String(), content)
	}
}
