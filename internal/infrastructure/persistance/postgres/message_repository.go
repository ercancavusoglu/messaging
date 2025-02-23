package postgres

import (
	"database/sql"
	"fmt"

	"github.com/ercancavusoglu/messaging/internal/domain/message"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) UpdateStatus(id int64, status message.MessageStatus, messageID string) error {
	query := `
		UPDATE messages
		SET status = $1, message_id = $2
		WHERE id = $3
	`

	_, err := r.db.Exec(query, status, messageID, id)
	return err
}

func (r *MessageRepository) List() ([]*message.Message, error) {
	fmt.Println("List messages")
	query := `
		SELECT id, recipient, content, status, sent_at, message_id, created_at
		FROM messages
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	defer rows.Close()

	var messages []*message.Message
	for rows.Next() {
		var msg message.Message
		var sentAt sql.NullTime
		var messageID sql.NullString

		err := rows.Scan(
			&msg.ID,
			&msg.To,
			&msg.Content,
			&msg.Status,
			&sentAt,
			&messageID,
			&msg.CreatedAt,
		)
		if err != nil {
			fmt.Println("Error scanning row", err)
			return nil, err
		}

		if sentAt.Valid {
			msg.SentAt = &sentAt.Time
		}
		if messageID.Valid {
			msg.MessageID = &messageID.String
		}

		messages = append(messages, &msg)
	}

	return messages, nil
}

func (r *MessageRepository) GetPendingMessages(limit int) ([]*message.Message, error) {
	query := `
		SELECT id, recipient, content, status, sent_at, message_id, created_at
		FROM messages
		WHERE status = $1
		ORDER BY created_at ASC
		LIMIT $2
	`

	rows, err := r.db.Query(query, message.StatusPending, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*message.Message
	for rows.Next() {
		var msg message.Message
		var sentAt sql.NullTime
		var messageID sql.NullString

		err := rows.Scan(
			&msg.ID,
			&msg.To,
			&msg.Content,
			&msg.Status,
			&sentAt,
			&messageID,
			&msg.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if sentAt.Valid {
			msg.SentAt = &sentAt.Time
		}
		if messageID.Valid {
			msg.MessageID = &messageID.String
		}

		messages = append(messages, &msg)
	}

	return messages, nil
}
