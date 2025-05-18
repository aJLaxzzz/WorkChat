package app

import (
	"chat/internal/utils"
	"log"
	"net/http"
	"strconv"
)

func (a *App) editMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	messageID := r.FormValue("message_id")
	newContent := r.FormValue("content")
	chatID := r.FormValue("chat_id")

	id, err := strconv.Atoi(messageID)
	if err != nil {
		log.Printf("editMessageHandler: strconv.Atoi: %v", err)
		http.Error(w, "Неверный идентификатор сообщения", http.StatusBadRequest)
		return
	}

	// Шифруем новое содержимое сообщения
	encryptedContent, err := a.cipher.Encrypt(newContent)
	if err != nil {
		log.Printf("editMessageHandler: cipher.Encrypt: %v", err)
		http.Error(w, "Ошибка шифрования сообщения: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = a.storage.UpdateMessageContent(messageID, encryptedContent)
	if err != nil {
		log.Printf("editMessageHandler: storage.UpdateMessageContent: %v", err)
		http.Error(w, "Ошибка редактирования сообщения: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Editing message with ID:", id)
	username, err := a.storage.GetUsernameByMessageID(id)
	if err != nil {
		log.Printf("editMessageHandler: storage.GetUsernameByMessageID: %v", err)
		http.Error(w, "Ошибка получения имени пользователя: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем уведомление всем клиентам (с расшифровкой для отображения)
	clients := a.memory.GetClientsByChatID(utils.Atoi(chatID))
	for _, client := range clients {
		err := client.Conn.WriteJSON(map[string]interface{}{
			"action":   "edit",
			"id":       id,
			"content":  newContent, // уже расшифрованное содержимое
			"Username": username,
		})
		if err != nil {
			log.Printf("editMessageHandler: client.Conn.WriteJSON: %v", err)
			client.Conn.Close()
			a.memory.DeleteClient(client)
		}
	}
}
