package app

import (
	"chat/internal/config"
	"chat/internal/domain"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

func (a *App) createGroupChatHandler(w http.ResponseWriter, r *http.Request) {
	if !a.isAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	session, _ := a.memory.GetSession(r, "session-name")
	username := session.Values["username"].(string)

	currentUserID, err := a.storage.GetUserIDByUsername(username)
	if err != nil {
		log.Printf("createGroupChatHandler: storage.GetUserIDByUsername: %v", err)
		http.Error(w, "Ошибка получения пользователя", http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		chatName := r.FormValue("chat_name")
		isPrivate := false

		chat := domain.Chat{
			Name:      chatName,
			IsPrivate: isPrivate,
			CreatorID: currentUserID,
		}
		chat.ID, err = a.storage.InsertChat(chat)
		if err != nil {
			log.Printf("createGroupChatHandler: storage.InsertChat: %v", err)
			http.Error(w, "Ошибка создания чата", http.StatusInternalServerError)
			return
		}

		err = a.storage.AddUserToChat(chat.ID, currentUserID)
		if err != nil {
			log.Printf("createGroupChatHandler: storage.AddUserToChat: %v", err)
			http.Error(w, "Ошибка добавления пользователя в чат", http.StatusInternalServerError)
			return
		}

		userIDs := r.Form["user_ids"] // Получаем массив ID пользователей
		for _, userIDToAddStr := range userIDs {
			userIDToAdd, err := strconv.Atoi(userIDToAddStr)
			if err != nil {
				log.Printf("createGroupChatHandler: strconv.Atoi: %v", err)
				http.Error(w, "Ошибка получения ID пользователя", http.StatusBadRequest)
				return
			}
			err = a.storage.AddUserToChat(chat.ID, userIDToAdd)
			if err != nil {
				log.Printf("createGroupChatHandler: storage.AddUserToChat: %v", err)
				http.Error(w, "Ошибка добавления пользователя в чат", http.StatusInternalServerError)
				return
			}
		}

		http.Redirect(w, r, "/chats", http.StatusSeeOther)
		return
	}

	users, err := a.storage.GetAllOtherUsers(username)
	if err != nil {
		log.Printf("createGroupChatHandler: storage.GetAllOtherUsers: %v", err)
		http.Error(w, "Ошибка получения пользователей", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles(filepath.Join(config.TemplatesDirPath, "create_group_chat.html")))
	err = tmpl.Execute(w, struct {
		Users []domain.User
	}{Users: users})
	if err != nil {
		log.Printf("createGroupChatHandler: tmpl.Execute: %v", err)
	}
}
