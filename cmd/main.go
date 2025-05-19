package main

import (
	"chat/internal/app"
	"chat/internal/config"
	"chat/internal/domain"
	"chat/internal/service/cipher"
	"chat/internal/service/memory"
	"chat/internal/storage"
	"log"

	"golang.org/x/crypto/bcrypt"
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("bcrypt.GenerateFromPassword: %v", err)
	}

	exists, err := storage.UserExists("admin")
	if err != nil {
		log.Fatalf("checking admin user existence: %v", err)
	}

	if !exists {
		err = storage.InsertUser(domain.User{
			Username:   "admin",
			Name:       "Администратор",
			Surname:    "Системы",
			Patronymic: "",
			Password:   string(hashedPassword),
		})
		if err != nil {
			log.Fatalf("creating admin user: %v", err)
		}
		log.Println("Admin user created")
	}

	app, err := app.NewApp(cfg, storage, memory, cipher)
	if err != nil {
		log.Fatalf("app.NewApp: %v", err)
	}

	err = app.Run()
	if err != nil {
		log.Fatalf("app.Run: %v", err)
	}
}
