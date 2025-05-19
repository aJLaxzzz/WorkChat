package domain

import (
	"time"

	"github.com/gorilla/websocket"
)

type User struct {
	ID         int
	Username   string
	Name       string
	Surname    string
	Patronymic string
	Password   string
	Status     string
	LastActive time.Time
}

type Chat struct {
	ID        int
	Name      string
	IsPrivate bool
	CreatorID int
	CreatedAt time.Time
}

type UserChat struct {
	Chat
	LastVisit          time.Time
	UnreadMessageCount int
	LastMessage        *Message
}

type File struct {
	Name string
	Data string
}

type Message struct {
	ID        int
	ChatID    int
	UserID    int
	Content   string
	CreatedAt time.Time
	Username  string
	File      File
	FullName  string
}

type Client struct {
	Conn   *websocket.Conn
	UserID int
	ChatID int
}
