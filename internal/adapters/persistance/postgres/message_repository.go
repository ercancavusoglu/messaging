package postgres

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ercancavusoglu/messaging/internal/domain"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) UpdateStatus(id int64, status domain.MessageStatus, messageID string, provider string) error {
	log.Printf("[MessageRepository] Updating message status [id: %d, status: %s, messageID: %s, provider: %s]", id, status, messageID, provider)
	query := `
		UPDATE messages 
		SET message_status = $1::varchar, message_id = $2, provider = $3, sent_at = CASE WHEN $1::varchar = 'sent' THEN NOW() ELSE sent_at END
		WHERE id = $4
	`
	result, err := r.db.Exec(query, status, messageID, provider, id)
	if err != nil {
		return fmt.Errorf("failed to update message status: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no message found with id: %d", id)
	}

	log.Printf("[MessageRepository] Message status updated successfully [id: %d]", id)
	return nil
}

func (r *MessageRepository) GetPendingMessages() ([]*domain.Message, error) {
	query := `
		SELECT id, recipient, content, message_status, message_id, provider, created_at, sent_at
		FROM messages
		WHERE message_status = $1
		ORDER BY created_at ASC
		LIMIT $2
	`

	rows, err := r.db.Query(query, domain.StatusPending)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending messages: %v", err)
	}
	defer rows.Close()

	var messages []*domain.Message
	for rows.Next() {
		msg := &domain.Message{}
		var sentAt sql.NullTime
		var messageID sql.NullString
		var provider sql.NullString

		err := rows.Scan(
			&msg.ID,
			&msg.To,
			&msg.Content,
			&msg.Status,
			&messageID,
			&provider,
			&msg.CreatedAt,
			&sentAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %v", err)
		}

		if messageID.Valid {
			msg.MessageID = messageID.String
		}
		if provider.Valid {
			msg.Provider = provider.String
		}
		if sentAt.Valid {
			msg.SentAt = &sentAt.Time
		}

		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating messages: %v", err)
	}

	return messages, nil
}

func (r *MessageRepository) GetByStatus(status domain.MessageStatus) ([]*domain.Message, error) {
	query := `
		SELECT id, recipient, content, message_status, message_id, provider, created_at, sent_at
		FROM messages
		WHERE message_status = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(query, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages by status: %v", err)
	}
	defer rows.Close()

	var messages []*domain.Message
	for rows.Next() {
		msg := &domain.Message{}
		var sentAt sql.NullTime
		var messageID sql.NullString
		var provider sql.NullString

		err := rows.Scan(
			&msg.ID,
			&msg.To,
			&msg.Content,
			&msg.Status,
			&messageID,
			&provider,
			&msg.CreatedAt,
			&sentAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %v", err)
		}

		if messageID.Valid {
			msg.MessageID = messageID.String
		}
		if provider.Valid {
			msg.Provider = provider.String
		}
		if sentAt.Valid {
			msg.SentAt = &sentAt.Time
		}

		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating messages: %v", err)
	}

	return messages, nil
}
