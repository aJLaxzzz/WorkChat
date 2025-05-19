package app

import (
	"chat/internal/utils"
	"log"
	"net/http"
)

func (a *App) deleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	messageID := r.FormValue("message_id")
	chatID := r.FormValue("chat_id")

	err := a.storage.DeleteMessage(messageID)
	if err != nil {
		log.Printf("deleteMessageHandler: storage.DeleteMessage: %v", err)
		http.Error(w, "Ошибка удаления сообщения", http.StatusInternalServerError)
		return
	}

	clients := a.memory.GetClientsByChatID(utils.Atoi(chatID))
	for _, client := range clients {
		err := client.Conn.WriteJSON(map[string]interface{}{
			"action": "delete",
			"id":     messageID,
		})
		if err != nil {
			log.Printf("deleteMessageHandler: client.Conn.WriteJSON: %v", err)
			client.Conn.Close()
			a.memory.DeleteClient(client)
		}
	}
}
