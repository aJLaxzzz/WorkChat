package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"chat/internal/app"
	"chat/internal/config"
	"chat/internal/domain"
	"chat/internal/service/cipher"
	"chat/internal/service/memory"
	"chat/internal/storage"
)

func FuzzAuth(f *testing.F) {
	f.Add("testuser", "password123")
	f.Add("", "")
	f.Add("verylongusername", "verylongpassword")

	f.Fuzz(func(t *testing.T, username, password string) {
		os.Setenv(config.ConfigPathEnvKey, "../../config.yaml")
		cfg, err := config.NewConfig()
		if err != nil {
			t.Fatalf("Failed to create config: %v", err)
		}
		storage, err := storage.NewStorage(cfg)
		if err != nil {
			t.Fatalf("Failed to create storage: %v", err)
		}
		memoryService := memory.NewService(cfg)
		cipherService := cipher.NewService(cfg)
		app, err := app.NewApp(cfg, storage, memoryService, cipherService)
		if err != nil {
			t.Fatalf("Failed to create app: %v", err)
		}

		user := domain.User{
			Username: username,
			Password: password,
		}
		body, _ := json.Marshal(user)

		req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		router := app.GetRouter()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			req = httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
			w = httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Login failed after successful registration")
			}
		}
	})
}
