package app

import (
	"chat/internal/config"
	"chat/internal/domain"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func (a *App) chatsHandler(w http.ResponseWriter, r *http.Request) {
	if !a.isAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	session, _ := a.memory.GetSession(r, "session-name")
	username := session.Values["username"].(string)

	user, err := a.storage.GetUserByUsername(username)
	if err != nil {
		log.Printf("chatsHandler: storage.GetUserByUsername: %v", err)
		http.Error(w, "Ошибка получения пользователя", http.StatusInternalServerError)
		return
	}

	chats, err := a.storage.GetChatsByUserID(user.ID)
	if err != nil {
		log.Printf("chatsHandler: storage.GetChatsByUserID: %v", err)
		http.Error(w, "Ошибка получения чатов", http.StatusInternalServerError)
		return
	}

	for i, chat := range chats {
		chats[i].UnreadMessageCount, err = a.storage.CountUnreadMessages(chat.ID, user.ID, chat.LastVisit)
		if err != nil {
			log.Printf("chatsHandler: storage.CountUnreadMessages: %v", err)
			chats[i].UnreadMessageCount = 0
		}

		lastMessage, err := a.storage.GetLastMessage(chat.ID)
		if err != nil {
			log.Printf("chatsHandler: storage.GetLastMessage: %v", err)
			chats[i].LastMessage = nil
		} else if lastMessage != nil {
			log.Printf("DEBUG: lastMessage.Content = %q", lastMessage.Content)

			if lastMessage.Content != "" {
				decryptedContent, err := a.cipher.Decrypt(lastMessage.Content)
				if err != nil {
					log.Printf("chatsHandler: cipher.Decrypt: %v", err)
					chats[i].LastMessage = nil
					continue
				}
				lastMessage.Content = decryptedContent
			}

			chats[i].LastMessage = lastMessage
		} else {
			chats[i].LastMessage = nil
		}
	}

	fullName := fmt.Sprintf("%s %s %s", user.Surname, user.Name, user.Patronymic)

	tmpl := template.Must(template.ParseFiles(filepath.Join(config.TemplatesDirPath, "chats.html")))
	err = tmpl.Execute(w, struct {
		FullName string
		Chats    []domain.UserChat
	}{
		FullName: fullName,
		Chats:    chats,
	})
	if err != nil {
		log.Printf("chatsHandler: tmpl.Execute: %v", err)
	}
}
