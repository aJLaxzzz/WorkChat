package app

import (
	"chat/internal/config"
	"chat/internal/domain"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
)

type Storage interface {
	GetChatByID(chatID int) (*domain.Chat, error)
	GetMessagesByChatID(chatID int) ([]domain.Message, error)
	GetChatMembersByChatID(chatID int) ([]domain.User, error)
	GetUserIDByUsername(username string) (int, error)
	GetUserByUsername(username string) (domain.User, error)
	GetChatsByUserID(userID int) ([]domain.UserChat, error)
	GetAllOtherUsers(username string) ([]domain.User, error)
	InsertChat(chat domain.Chat) (int, error)
	AddUserToChat(chatID int, userID int) error
	GetUserByID(id int) (domain.User, error)
	GetChatIDByUserIDs(firstID int, secondID int) (int, error)
	DeleteMessage(messageID string) error
	UpdateMessageContent(messageID string, content string) error
	GetUsernameByMessageID(messageID int) (string, error)
	UpdateUserStatus(username string, status string) error
	InsertUser(user domain.User) error
	InsertMessage(message domain.Message) (int, error)
	UpdateLastChatVisitTime(chatID int, userID int) error
	CountUnreadMessages(chatID int, userID int, timepoint time.Time) (int, error)
	GetLastMessage(chatID int) (*domain.Message, error)
}

type Memory interface {
	GetSession(r *http.Request, name string) (*sessions.Session, error)
	GetClientsByChatID(chatID int) []domain.Client
	DeleteClient(client domain.Client)
	AddClient(client domain.Client)
}

type Cipher interface {
	Encrypt(plainText string) (string, error)
	Decrypt(cipherText string) (string, error)
}

type App struct {
	cfg      *config.Config
	router   *mux.Router
	upgrader websocket.Upgrader
	storage  Storage
	memory   Memory
	cipher   Cipher
}

func NewApp(cfg *config.Config, storage Storage, memory Memory, cipher Cipher) (*App, error) {
	r := mux.NewRouter()
	app := App{
		cfg:    cfg,
		router: r,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		storage: storage,
		memory:  memory,
		cipher:  cipher,
	}

	r.HandleFunc("/", app.loginHandler).Methods("GET", "POST")

	r.HandleFunc("/register", app.registerHandler).Methods("GET", "POST")
	r.HandleFunc("/login", app.loginHandler).Methods("GET", "POST")
	r.HandleFunc("/chats", app.chatsHandler).Methods("GET")
	r.HandleFunc("/chat/{id:[0-9]+}", app.chatHandler).Methods("GET")
	r.HandleFunc("/ws/chat/{id:[0-9]+}", app.wsChatHandler) // Обработчик WebSocket
	r.HandleFunc("/logout", app.logoutHandler).Methods("POST")

	r.HandleFunc("/create_private_chat", app.createPrivateChatHandler).Methods("GET", "POST")
	r.HandleFunc("/create_group_chat", app.createGroupChatHandler).Methods("GET", "POST")

	r.HandleFunc("/edit-message", app.editMessageHandler).Methods("POST")
	r.HandleFunc("/delete-message", app.deleteMessageHandler).Methods("POST")

	return &app, nil
}

func (a *App) Run() error {
	server := http.Server{
		Addr:    fmt.Sprintf("%s:%s", a.cfg.Server.Host, a.cfg.Server.Port),
		Handler: a.router,
	}
	log.Printf("Starting server on %s", server.Addr)
	return server.ListenAndServe()
}

func (a *App) isAuthenticated(r *http.Request) bool {
	session, _ := a.memory.GetSession(r, "session-name")
	_, ok := session.Values["username"].(string)
	return ok
}
