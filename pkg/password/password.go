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

// HashPassword generates a BCrypt hash for the given password using specified cost factor.
// Password length must be 1-64 characters. Cost range: 4-31 (default: 10).
// Returns BCrypt hash string or error.
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

// CheckPasswordHash verifies if plaintext password matches the BCrypt hash.
// Returns true if passwords match, false otherwise.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// EncryptPassword encrypts password using AES-256-GCM authenticated encryption.
// Requires exactly 32-byte key (AES-256). Generates random nonce.
// Returns base64-encoded nonce+ciphertext for safe storage/transmission.
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

// DecryptPassword decrypts base64-encoded AES-256-GCM ciphertext back to plaintext.
// Extracts nonce from first gcm.NonceSize() bytes. Validates integrity.
// Returns original password or error if decryption fails.
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
