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

func FuzzMessage(f *testing.F) {
	f.Add("test message", "testuser", 1)
	f.Add("", "", 0)
	f.Add("very long message", "verylongusername", 999)

	f.Fuzz(func(t *testing.T, content, username string, chatID int) {
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

		message := domain.Message{
			Content:  content,
			Username: username,
			ChatID:   chatID,
		}
		body, _ := json.Marshal(message)

		req := httptest.NewRequest("POST", "/ws/chat/1", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		router := app.GetRouter()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Message sending failed with status: %d", w.Code)
		}
	})
}
