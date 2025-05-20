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

func FuzzFileUpload(f *testing.F) {
	f.Add("testfile.txt", "test content")
	f.Add("", "")
	f.Add("verylongfilename.txt", "very long file content")

	f.Fuzz(func(t *testing.T, filename, content string) {
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

		file := domain.File{
			Name: filename,
			Data: content,
		}
		body, _ := json.Marshal(file)

		req := httptest.NewRequest("POST", "/files", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		router := app.GetRouter()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("File upload failed with status: %d", w.Code)
		}
	})
}
