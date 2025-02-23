package postgres

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ercancavusoglu/messaging/internal/domain/message"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(msg *message.Message) error {
	query := `
		INSERT INTO messages (recipient, content, message_status, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := r.db.QueryRow(query, msg.To, msg.Content, msg.Status, msg.CreatedAt).Scan(&msg.ID)
	if err != nil {
		return fmt.Errorf("failed to create message: %v", err)
	}

	return nil
}

func (r *MessageRepository) GetByID(id int64) (*message.Message, error) {
	query := `
		SELECT id, recipient, content, message_status, message_id, created_at, sent_at
		FROM messages
		WHERE id = $1
	`

	msg := &message.Message{}
	var sentAt sql.NullTime
	var messageID sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&msg.ID,
		&msg.To,
		&msg.Content,
		&msg.Status,
		&messageID,
		&msg.CreatedAt,
		&sentAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get message: %v", err)
	}

	if messageID.Valid {
		msg.MessageID = messageID.String
	}
	if sentAt.Valid {
		msg.SentAt = &sentAt.Time
	}

	return msg, nil
}

func (r *MessageRepository) UpdateStatus(id int64, status message.MessageStatus, messageID string) error {
	log.Printf("[MessageRepository] Updating message status [id: %d, status: %s, messageID: %s]", id, status, messageID)
	query := `
		UPDATE messages 
		SET message_status = $1::varchar, message_id = $2, sent_at = CASE WHEN $1::varchar = 'sent' THEN NOW() ELSE sent_at END
		WHERE id = $3
	`
	result, err := r.db.Exec(query, status, messageID, id)
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

func (r *MessageRepository) GetPendingMessages(limit int) ([]*message.Message, error) {
	query := `
		SELECT id, recipient, content, message_status, message_id, created_at, sent_at
		FROM messages
		WHERE message_status = $1
		ORDER BY created_at ASC
		LIMIT $2
	`

	rows, err := r.db.Query(query, message.StatusPending, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending messages: %v", err)
	}
	defer rows.Close()

	var messages []*message.Message
	for rows.Next() {
		msg := &message.Message{}
		var sentAt sql.NullTime
		var messageID sql.NullString

		err := rows.Scan(
			&msg.ID,
			&msg.To,
			&msg.Content,
			&msg.Status,
			&messageID,
			&msg.CreatedAt,
			&sentAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %v", err)
		}

		if messageID.Valid {
			msg.MessageID = messageID.String
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

func (r *MessageRepository) GetByStatus(status message.MessageStatus) ([]*message.Message, error) {
	query := `
		SELECT id, recipient, content, message_status, message_id, created_at, sent_at
		FROM messages
		WHERE message_status = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(query, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages by status: %v", err)
	}
	defer rows.Close()

	var messages []*message.Message
	for rows.Next() {
		msg := &message.Message{}
		var sentAt sql.NullTime
		var messageID sql.NullString

		err := rows.Scan(
			&msg.ID,
			&msg.To,
			&msg.Content,
			&msg.Status,
			&messageID,
			&msg.CreatedAt,
			&sentAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %v", err)
		}

		if messageID.Valid {
			msg.MessageID = messageID.String
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

func (r *MessageRepository) List() ([]*message.Message, error) {
	query := `
		SELECT id, recipient, content, message_status, message_id, created_at, sent_at
		FROM messages
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list messages: %v", err)
	}
	defer rows.Close()

	var messages []*message.Message
	for rows.Next() {
		msg := &message.Message{}
		var sentAt sql.NullTime
		var messageID sql.NullString

		err := rows.Scan(
			&msg.ID,
			&msg.To,
			&msg.Content,
			&msg.Status,
			&messageID,
			&msg.CreatedAt,
			&sentAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %v", err)
		}

		if messageID.Valid {
			msg.MessageID = messageID.String
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
