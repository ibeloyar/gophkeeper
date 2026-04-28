package password

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword_EmptyPassword(t *testing.T) {
	hash, err := HashPassword("", 4)

	assert.ErrorIs(t, err, ErrPasswordRequired)
	assert.Empty(t, hash)
}

func TestHashPassword_TooLongPassword(t *testing.T) {
	longPassword := string(make([]byte, 65))

	hash, err := HashPassword(longPassword, 4)

	assert.ErrorIs(t, err, ErrPasswordMaxLen64)
	assert.Empty(t, hash)
}

func TestHashPassword_ValidPassword(t *testing.T) {
	password := "testpass123"

	hash, err := HashPassword(password, 4)
	assert.NoError(t, err)

	assert.Contains(t, hash, "$2a$")
	assert.Contains(t, hash, "04$")

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	assert.NoError(t, err)
}

func TestHashPassword_BcryptError(t *testing.T) {
	hash, err := HashPassword("testpass", 32)

	assert.ErrorIs(t, err, ErrPasswordGenerate)
	assert.Empty(t, hash)
}

func TestCheckPasswordHash_Valid(t *testing.T) {
	hash, err := HashPassword("testpass", 4)
	assert.NoError(t, err)

	result := CheckPasswordHash("testpass", hash)
	assert.True(t, result)
}

func TestCheckPasswordHash_Invalid(t *testing.T) {
	hash, _ := HashPassword("testpass", 4)

	result := CheckPasswordHash("wrongpass", hash)
	assert.False(t, result)
}

func TestHashPassword_Basic(t *testing.T) {
	tests := []struct {
		password string
		cost     int
	}{
		{"pass", 4},
		{"a" + string(make([]byte, 63)), 4}, // max len
	}

	for _, tt := range tests {
		t.Run(tt.password, func(t *testing.T) {
			hash, err := HashPassword(tt.password, tt.cost)
			assert.NoError(t, err)
			assert.NotEmpty(t, hash)
		})
	}
}

func TestEncryptPassword_EmptyPassword(t *testing.T) {
	_, err := EncryptPassword("", "12345678901234567890123456789012") // 32 bytes

	assert.ErrorIs(t, err, ErrDecryptFailed)
}

func TestEncryptPassword_TooLongPassword(t *testing.T) {
	longPassword := string(make([]byte, 65))
	_, err := EncryptPassword(longPassword, "12345678901234567890123456789012")

	assert.ErrorIs(t, err, ErrDecryptFailed)
}

func TestEncryptPassword_InvalidKey(t *testing.T) {
	_, err := EncryptPassword("testpass", "short")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "key must be 32 bytes long")
}

func TestEncryptPassword_Valid(t *testing.T) {
	key := "12345678901234567890123456789012" // exactly 32 bytes
	password := "testpass123"

	encrypted, err := EncryptPassword(password, key)
	assert.NoError(t, err)
	assert.NotEmpty(t, encrypted)

	// Verify it's base64
	_, decodeErr := base64.StdEncoding.DecodeString(encrypted)
	assert.NoError(t, decodeErr)
}

func TestDecryptPassword_InvalidKey(t *testing.T) {
	key := "12345678901234567890123456789012"
	encrypted, _ := EncryptPassword("testpass", key)

	_, err := DecryptPassword(encrypted, "short")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "key must be 32 bytes long")
}

func TestDecryptPassword_InvalidBase64(t *testing.T) {
	key := "12345678901234567890123456789012"
	_, err := DecryptPassword("invalidbase64!!!", key)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "illegal base64 data")
}

func TestDecryptPassword_ShortData(t *testing.T) {
	key := "12345678901234567890123456789012"
	shortData := "YQ==" // too short
	_, err := DecryptPassword(shortData, key)

	assert.ErrorIs(t, err, ErrDecryptFailed)
}

func TestDecryptPassword_WrongKey(t *testing.T) {
	key1 := "12345678901234567890123456789012"
	key2 := "another32bytekey123456789012345678901234567890"[:32]

	encrypted, _ := EncryptPassword("testpass", key1)
	_, err := DecryptPassword(encrypted, key2)

	assert.ErrorIs(t, err, ErrDecryptFailed)
}

func TestDecryptPassword_InvalidAuth(t *testing.T) {
	key := "12345678901234567890123456789012"
	// create a valid base64, but with corrupted data.
	data := make([]byte, 50)
	base64Data := base64.StdEncoding.EncodeToString(data)

	_, err := DecryptPassword(base64Data, key)

	assert.ErrorIs(t, err, ErrDecryptFailed)
}

func TestEncryptDecrypt_Roundtrip(t *testing.T) {
	key := "12345678901234567890123456789012"
	tests := []struct {
		password string
	}{
		{"testpass"},
		{"a"},
		{"a" + string(make([]byte, 63))}, // max len
	}

	for _, tt := range tests {
		t.Run(tt.password, func(t *testing.T) {
			encrypted, err := EncryptPassword(tt.password, key)
			assert.NoError(t, err)

			decrypted, err := DecryptPassword(encrypted, key)
			assert.NoError(t, err)

			assert.Equal(t, tt.password, decrypted)
		})
	}
}
