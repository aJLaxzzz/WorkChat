package app

import (
	"chat/internal/config"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
)

func (a *App) loginHandler(w http.ResponseWriter, r *http.Request) {
	if a.isAuthenticated(r) {
		http.Redirect(w, r, "/chats", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		user, err := a.storage.GetUserByUsername(username)
		if err != nil {
			log.Printf("loginHandler: storage.GetUserByUsername: %v", err)
			http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			log.Printf("loginHandler: bcrypt.CompareHashAndPassword: %v", err)
			http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
			return
		}

		session, _ := a.memory.GetSession(r, "session-name")
		session.Values["username"] = username
		err = session.Save(r, w)
		if err != nil {
			log.Printf("loginHandler: session.Save: %v", err)
			http.Error(w, "Ошибка сохранения сессии", http.StatusInternalServerError)
			return
		}

		err = a.storage.UpdateUserStatus(username, "online")
		if err != nil {
			log.Printf("loginHandler: storage.UpdateUserStatus: %v", err)
			http.Error(w, "Ошибка обновления статуса", http.StatusInternalServerError)
			return
		}

		if username == "admin" {
			http.Redirect(w, r, "/register", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/chats", http.StatusSeeOther)
		}
		return
	}

	tmpl := template.Must(template.ParseFiles(filepath.Join(config.TemplatesDirPath, "login.html")))
	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("loginHandler: tmpl.Execute: %v", err)
	}
}
