package app

import (
	"chat/internal/config"
	"chat/internal/domain"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
)

func (a *App) registerHandler(w http.ResponseWriter, r *http.Request) {
	if a.isAuthenticated(r) {
		http.Redirect(w, r, "/chats", http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		name := r.FormValue("name")
		surname := r.FormValue("surname")
		patronymic := r.FormValue("patronymic")
		// Хешируем пароль
		password := r.FormValue("password")

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("registerHandler: bcrypt.GenerateFromPassword: %v", err)
			http.Error(w, "Ошибка регистрации", http.StatusInternalServerError)
			return
		}

		user := domain.User{
			Username:   username,
			Name:       name,
			Surname:    surname,
			Patronymic: patronymic,
			Password:   string(hashedPassword),
		}
		err = a.storage.InsertUser(user)
		if err != nil {
			log.Printf("registerHandler: storage.InsertUser: %v", err)
			http.Error(w, "Ошибка регистрации", http.StatusInternalServerError)
			return
		}

		session, _ := a.memory.GetSession(r, "session-name")
		session.Values["username"] = username
		err = session.Save(r, w)
		if err != nil {
			log.Printf("registerHandler: session.Save: %v", err)
			http.Error(w, "Ошибка сохранения сессии", http.StatusInternalServerError)
			return
		}

		err = a.storage.UpdateUserStatus(username, "online")
		if err != nil {
			log.Printf("registerHandler: storage.UpdateUserStatus: %v", err)
			http.Error(w, "Ошибка обновления статуса", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/chats", http.StatusSeeOther)
		return
	}
	tmpl := template.Must(template.ParseFiles(filepath.Join(config.TemplatesDirPath, "register.html")))
	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("registerHandler: tmpl.Execute: %v", err)
	}
}
