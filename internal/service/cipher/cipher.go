package cipher

import (
	"chat/internal/config"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

type Service struct {
	encryptionKey string
}

func NewService(cfg *config.Config) *Service {
	return &Service{
		encryptionKey: cfg.EncryptionKey,
	}
}

func (s *Service) Encrypt(plainText string) (string, error) {
	block, err := aes.NewCipher([]byte(s.encryptionKey))
	if err != nil {
		return "", err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}

	stream := cipher.NewCTR(block, iv)

	cipherText := make([]byte, len(plainText))
	stream.XORKeyStream(cipherText, []byte(plainText))

	// Вернем IV + ciphertext в base64
	result := append(iv, cipherText...)
	return base64.StdEncoding.EncodeToString(result), nil
}

func (s *Service) Decrypt(cipherText string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	if len(data) < aes.BlockSize {
		return "", fmt.Errorf("invalid ciphertext")
	}

	iv := data[:aes.BlockSize]
	cipherData := data[aes.BlockSize:]

	block, err := aes.NewCipher([]byte(s.encryptionKey))
	if err != nil {
		return "", err
	}

	stream := cipher.NewCTR(block, iv)

	plainText := make([]byte, len(cipherData))
	stream.XORKeyStream(plainText, cipherData)

	return string(plainText), nil
}
