package password

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordRequired = errors.New("password is required")
	ErrPasswordMaxLen64 = errors.New("password too long, max 64 characters")
	ErrPasswordGenerate = errors.New("password generate error")
	ErrDecryptFailed    = errors.New("decrypt failed")
)

func HashPassword(password string, passCost int) (string, error) {
	if len(password) < 1 {
		return "", ErrPasswordRequired
	}
	if len(password) > 64 {
		return "", ErrPasswordMaxLen64
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), passCost)
	if err != nil {
		return "", ErrPasswordGenerate
	}

	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func EncryptPassword(password string, keyStr string) (string, error) {
	if len(password) < 1 {
		return "", ErrDecryptFailed
	}
	if len(password) > 64 {
		return "", ErrDecryptFailed
	}

	key := []byte(keyStr)
	if len(key) != 32 { // AES‑256
		return "", errors.New("key must be 32 bytes long")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := gcm.Seal(nil, nonce, []byte(password), nil)
	out := append(nonce, cipherText...)

	return base64.StdEncoding.EncodeToString(out), nil
}

func DecryptPassword(cipherTextB64, keyStr string) (string, error) {
	key := []byte(keyStr)
	if len(key) != 32 {
		return "", errors.New("key must be 32 bytes long")
	}

	data, err := base64.StdEncoding.DecodeString(cipherTextB64)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", ErrDecryptFailed
	}

	nonce, cipherText := data[:nonceSize], data[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", ErrDecryptFailed
	}

	return string(plainText), nil
}
