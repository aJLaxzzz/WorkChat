package app

import (
	"log"
	"net/http"
)

func (a *App) logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	session, _ := a.memory.GetSession(r, "session-name")
	username, ok := session.Values["username"].(string)

	if !ok {
		http.Error(w, "Пользователь не авторизован", http.StatusUnauthorized)
		return
	}

	err := a.storage.UpdateUserStatus(username, "offline")
	if err != nil {
		log.Printf("logoutHandler: storage.UpdateUserStatus: %v", err)
		http.Error(w, "Ошибка обновления статуса", http.StatusInternalServerError)
		return
	}

	delete(session.Values, "username")
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
