package storage

import (
	"chat/internal/domain"
	"database/sql"
	"fmt"
)

func (s *Storage) DeleteMessage(messageID string) error {
	_, err := s.db.Exec("DELETE FROM messages WHERE id = $1", messageID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) UpdateMessageContent(messageID string, content string) error {
	_, err := s.db.Exec("UPDATE messages SET content = $1 WHERE id = $2", content, messageID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetUsernameByMessageID(messageID int) (string, error) {
	var username string
	err := s.db.QueryRow(
		"SELECT u.username FROM messages m JOIN users u ON m.user_id = u.id WHERE m.id = $1", messageID,
	).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}

func (s *Storage) InsertMessage(message domain.Message) (int, error) {
	err := s.db.QueryRow(
		"INSERT INTO messages (chat_id, user_id, content, file_name, file_content) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		message.ChatID, message.UserID, message.Content, message.File.Name, message.File.Data,
	).Scan(&message.ID)
	if err != nil {
		return 0, err
	}
	return message.ID, nil
}

func (s *Storage) GetLastMessage(chatID int) (*domain.Message, error) {
	var message domain.Message
	err := s.db.QueryRow(
		`SELECT id, chat_id, user_id, content, created_at 
		FROM messages 
		WHERE chat_id = $1 
		ORDER BY created_at DESC 
		LIMIT 1`,
		chatID,
	).Scan(&message.ID, &message.ChatID, &message.UserID, &message.Content, &message.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("GetLastMessage: %w", err)
	}

	message.Username, err = s.GetUsernameByMessageID(message.ID)
	if err != nil {
		return nil, fmt.Errorf("GetLastMessage: %w", err)
	}

	return &message, nil
}

func (s *Storage) GetMessageByID(messageID string, message *domain.Message) error {
	err := s.db.QueryRow(
		"SELECT id, chat_id, user_id, content, file_name, file_content FROM messages WHERE id = $1",
		messageID,
	).Scan(&message.ID, &message.ChatID, &message.UserID, &message.Content, &message.File.Name, &message.File.Data)
	return err
}
