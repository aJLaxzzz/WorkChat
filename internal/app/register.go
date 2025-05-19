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
	session, _ := a.memory.GetSession(r, "session-name")
	username, ok := session.Values["username"].(string)
	if !ok || username != "admin" {
		http.Error(w, "Доступ запрещён", http.StatusForbidden)
		return
	}

	var success bool

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		name := r.FormValue("name")
		surname := r.FormValue("surname")
		patronymic := r.FormValue("patronymic")
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

		err = a.storage.UpdateUserStatus(username, "online")
		if err != nil {
			log.Printf("registerHandler: storage.UpdateUserStatus: %v", err)
			http.Error(w, "Ошибка обновления статуса", http.StatusInternalServerError)
			return
		}

		success = true
	}

	tmpl := template.Must(template.ParseFiles(filepath.Join(config.TemplatesDirPath, "register.html")))
	err := tmpl.Execute(w, map[string]any{
		"Success": success,
	})
	if err != nil {
		log.Printf("registerHandler: tmpl.Execute: %v", err)
	}
}
