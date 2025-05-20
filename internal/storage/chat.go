package storage

import (
	"chat/internal/domain"
	"time"
)

func (s *Storage) GetChatByID(chatID int) (*domain.Chat, error) {
	rows, err := s.db.Query(
		"SELECT id, name, is_private, creator_id, created_at FROM chats WHERE id = $1",
		chatID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chat domain.Chat
	if rows.Next() {
		if err := rows.Scan(
			&chat.ID,
			&chat.Name,
			&chat.IsPrivate,
			&chat.CreatorID,
			&chat.CreatedAt,
		); err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}
	return &chat, nil
}

func (s *Storage) GetMessagesByChatID(chatID int) ([]domain.Message, error) {
	messageRows, err := s.db.Query(`
		SELECT m.id, m.user_id, m.content, m.created_at, 
		       u.surname || ' ' || u.name || ' ' || u.patronymic AS full_name, 
		       m.file_name, m.file_content 
		FROM messages m 
		JOIN users u ON m.user_id = u.id 
		WHERE chat_id = $1 
		ORDER BY m.created_at`, chatID)
	if err != nil {
		return nil, err
	}
	defer messageRows.Close()

	var messages []domain.Message
	for messageRows.Next() {
		var message domain.Message
		if err := messageRows.Scan(
			&message.ID,
			&message.UserID,
			&message.Content,
			&message.CreatedAt,
			&message.FullName,
			&message.File.Name,
			&message.File.Data,
		); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (s *Storage) GetChatMembersByChatID(chatID int) ([]domain.User, error) {
	membersRows, err := s.db.Query(
		`SELECT u.id, u.username, u.surname, u.name, u.patronymic, u.status, u.last_active
		 FROM chat_users cu
		 JOIN users u ON cu.user_id = u.id
		 WHERE cu.chat_id = $1`, chatID)

	if err != nil {
		return nil, err
	}
	defer membersRows.Close()

	var members []domain.User
	for membersRows.Next() {
		var member domain.User
		if err := membersRows.Scan(
			&member.ID,
			&member.Username,
			&member.Surname,
			&member.Name,
			&member.Patronymic,
			&member.Status,
			&member.LastActive,
		); err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	return members, nil
}

func (s *Storage) GetChatsByUserID(userID int) ([]domain.UserChat, error) {
	rows, err := s.db.Query(`
	SELECT c.id, 
	       CASE 
	           WHEN c.is_private THEN 
	               (SELECT surname || ' ' || name || ' ' || patronymic 
	                FROM users 
	                WHERE id != $1 AND id IN (SELECT user_id FROM chat_users WHERE chat_id = c.id))
	           ELSE 
	               c.name
	       END AS name,
	       c.is_private,
		   cu.last_chat_visit
	FROM chats c
	JOIN chat_users cu ON c.id = cu.chat_id
	WHERE cu.user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []domain.UserChat
	for rows.Next() {
		var chat domain.UserChat
		if err := rows.Scan(&chat.ID, &chat.Name, &chat.IsPrivate, &chat.LastVisit); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}

	return chats, nil
}

func (s *Storage) CountUnreadMessages(chatID int, userID int, timepoint time.Time) (int, error) {
	rows, err := s.db.Query(
		"SELECT count(*) FROM messages WHERE chat_id=$1 AND user_id!=$2 AND created_at > $3",
		chatID, userID, timepoint,
	)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var count int
	if rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return 0, err
		}
	} else {
		return 0, nil
	}
	return count, nil
}

func (s *Storage) InsertChat(chat domain.Chat) (int, error) {
	err := s.db.QueryRow(
		"INSERT INTO chats (name, is_private, creator_id) VALUES ($1, $2, $3) RETURNING id",
		chat.Name, chat.IsPrivate, chat.CreatorID,
	).Scan(&chat.ID)
	if err != nil {
		return 0, err
	}
	return chat.ID, nil
}

func (s *Storage) AddUserToChat(chatID int, userID int) error {
	_, err := s.db.Exec("INSERT INTO chat_users (chat_id, user_id) VALUES ($1, $2)", chatID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) UpdateLastChatVisitTime(chatID int, userID int) error {
	_, err := s.db.Exec("UPDATE chat_users SET last_chat_visit=NOW() WHERE chat_id=$1 AND user_id=$2", chatID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetChatIDByUserIDs(firstID int, secondID int) (int, error) {
	var chatID int
	err := s.db.QueryRow(
		`
			SELECT c.id FROM chats c
			JOIN chat_users cu1 ON c.id = cu1.chat_id
			JOIN chat_users cu2 ON c.id = cu2.chat_id
			WHERE cu1.user_id = $1 AND cu2.user_id = $2 AND c.is_private = true
		`,
		firstID, secondID,
	).Scan(&chatID)
	if err != nil {
		return 0, err
	}
	return chatID, nil
}
