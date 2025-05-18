package storage

import (
	"chat/internal/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func connectionString(cfg *config.Config) string {
	return fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s sslmode=disable",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
		cfg.DB.Host,
	)
}

func NewStorage(cfg *config.Config) (*Storage, error) {
	db, err := sql.Open("postgres", connectionString(cfg))
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) Close() {
	err := s.db.Close()
	if err != nil {
		log.Printf("Error closing connection to the databse: %v", err)
	}
}
