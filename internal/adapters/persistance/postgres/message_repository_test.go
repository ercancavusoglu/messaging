package postgres

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ercancavusoglu/messaging/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestMessageRepository_UpdateStatus(t *testing.T) {
	// Mock DB oluştur
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMessageRepository(db)

	// Test verileri
	id := int64(1)
	status := domain.StatusSent
	messageID := "msg_123"
	provider := "client_one"

	// Mock beklentileri
	mock.ExpectExec("UPDATE messages").
		WithArgs(status, messageID, provider, id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Test
	err = repo.UpdateStatus(id, status, messageID, provider)
	assert.NoError(t, err)

	// Mock beklentilerinin karşılandığını kontrol et
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_UpdateStatus_NoRowsAffected(t *testing.T) {
	// Mock DB oluştur
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMessageRepository(db)

	// Test verileri
	id := int64(1)
	status := domain.StatusSent
	messageID := "msg_123"
	provider := "client_one"

	// Mock beklentileri
	mock.ExpectExec("UPDATE messages").
		WithArgs(status, messageID, provider, id).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// Test
	err = repo.UpdateStatus(id, status, messageID, provider)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no message found with id")

	// Mock beklentilerinin karşılandığını kontrol et
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_GetPendingMessages(t *testing.T) {
	// Mock DB oluştur
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMessageRepository(db)

	// Test verileri
	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "recipient", "content", "message_status", "message_id", "provider", "created_at", "sent_at"}).
		AddRow(1, "+905551234567", "Test message 1", domain.StatusPending, sql.NullString{}, sql.NullString{}, now, sql.NullTime{}).
		AddRow(2, "+905551234568", "Test message 2", domain.StatusPending, sql.NullString{}, sql.NullString{}, now, sql.NullTime{})

	// Mock beklentileri
	mock.ExpectQuery("SELECT (.+) FROM messages").
		WithArgs(domain.StatusPending).
		WillReturnRows(rows)

	// Test
	messages, err := repo.GetPendingMessages()
	assert.NoError(t, err)
	assert.Len(t, messages, 2)

	// İlk mesajı kontrol et
	assert.Equal(t, int64(1), messages[0].ID)
	assert.Equal(t, "+905551234567", messages[0].To)
	assert.Equal(t, "Test message 1", messages[0].Content)
	assert.Equal(t, domain.StatusPending, messages[0].Status)
	assert.Empty(t, messages[0].MessageID)
	assert.Empty(t, messages[0].Provider)
	assert.Equal(t, now, messages[0].CreatedAt)
	assert.Nil(t, messages[0].SentAt)

	// Mock beklentilerinin karşılandığını kontrol et
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_GetByStatus(t *testing.T) {
	// Mock DB oluştur
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMessageRepository(db)

	// Test verileri
	now := time.Now()
	sentAt := now.Add(time.Hour)
	rows := sqlmock.NewRows([]string{"id", "recipient", "content", "message_status", "message_id", "provider", "created_at", "sent_at"}).
		AddRow(1, "+905551234567", "Test message 1", domain.StatusSent, "msg_123", "client_one", now, sentAt).
		AddRow(2, "+905551234568", "Test message 2", domain.StatusSent, "msg_124", "client_two", now, sentAt)

	// Mock beklentileri
	mock.ExpectQuery("SELECT (.+) FROM messages").
		WithArgs(domain.StatusSent).
		WillReturnRows(rows)

	// Test
	messages, err := repo.GetByStatus(domain.StatusSent)
	assert.NoError(t, err)
	assert.Len(t, messages, 2)

	// İlk mesajı kontrol et
	assert.Equal(t, int64(1), messages[0].ID)
	assert.Equal(t, "+905551234567", messages[0].To)
	assert.Equal(t, "Test message 1", messages[0].Content)
	assert.Equal(t, domain.StatusSent, messages[0].Status)
	assert.Equal(t, "msg_123", messages[0].MessageID)
	assert.Equal(t, "client_one", messages[0].Provider)
	assert.Equal(t, now, messages[0].CreatedAt)
	assert.Equal(t, sentAt, *messages[0].SentAt)

	// Mock beklentilerinin karşılandığını kontrol et
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMessageRepository_GetByStatus_NoRows(t *testing.T) {
	// Mock DB oluştur
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewMessageRepository(db)

	// Mock beklentileri
	mock.ExpectQuery("SELECT (.+) FROM messages").
		WithArgs(domain.StatusSent).
		WillReturnRows(sqlmock.NewRows([]string{"id", "recipient", "content", "message_status", "message_id", "provider", "created_at", "sent_at"}))

	// Test
	messages, err := repo.GetByStatus(domain.StatusSent)
	assert.NoError(t, err)
	assert.Empty(t, messages)

	// Mock beklentilerinin karşılandığını kontrol et
	assert.NoError(t, mock.ExpectationsWereMet())
}
