package storage

import "chat/internal/domain"

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
