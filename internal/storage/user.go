package storage

import (
	"chat/internal/domain"
)

func (s *Storage) GetUserIDByUsername(username string) (int, error) {
	var userID int
	err := s.db.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, err
}

func (s *Storage) GetUserByUsername(username string) (domain.User, error) {
	var user domain.User
	err := s.db.QueryRow("SELECT id, username, name, surname, patronymic, password FROM users WHERE username = $1", username).
		Scan(&user.ID, &user.Username, &user.Name, &user.Surname, &user.Patronymic, &user.Password)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (s *Storage) GetUserByID(id int) (domain.User, error) {
	var user domain.User
	err := s.db.QueryRow("SELECT id, username, name, surname, patronymic FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Username, &user.Name, &user.Surname, &user.Patronymic)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (s *Storage) GetAllOtherUsers(username string) ([]domain.User, error) {
	rows, err := s.db.Query(`
	SELECT id, name, surname, patronymic 
	FROM users 
	WHERE id != (SELECT id FROM users WHERE username = $1)`, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Surname, &user.Patronymic); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (s *Storage) UpdateUserStatus(username string, status string) error {
	_, err := s.db.Exec("UPDATE users SET status = $1, last_active = NOW() WHERE username = $2", status, username)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) InsertUser(user domain.User) error {
	_, err := s.db.Exec(
		"INSERT INTO users (username, name, surname, patronymic, password, status) VALUES ($1, $2, $3, $4, $5, 'offline')",
		user.Username, user.Name, user.Surname, user.Patronymic, user.Password,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) UserExists(username string) (bool, error) {
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)", username).Scan(&exists)
	return exists, err
}
