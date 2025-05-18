package main

import (
	"chat/internal/app"
	"chat/internal/config"
	"chat/internal/service/cipher"
	"chat/internal/service/memory"
	"chat/internal/storage"
	"log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("config.NewConfig: %v", err)
	}

	storage, err := storage.NewStorage(cfg)
	if err != nil {
		log.Fatalf("storage.NewStorage: %v", err)
	}
	defer storage.Close()

	memory := memory.NewService(cfg)
	cipher := cipher.NewService(cfg)

	app, err := app.NewApp(cfg, storage, memory, cipher)
	if err != nil {
		log.Fatalf("app.NewApp: %v", err)
	}

	err = app.Run()
	if err != nil {
		log.Fatalf("app.Run: %v", err)
	}
}
