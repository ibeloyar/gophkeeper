package model

import (
	"time"
)

type SecretType string

const (
	SecretTypeLoginPassword SecretType = "login_password"
	SecretTypeText          SecretType = "text"
	SecretTypeBinary        SecretType = "binary"
	SecretTypeCard          SecretType = "card"
)

type Secret struct {
	ID         int64      `json:"id"`
	UserID     int64      `json:"user_id"`
	Title      string     `json:"title"`
	Metadata   string     `json:"metadata"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	SecretType SecretType `json:"secret_type"`

	// login_password
	Login    string `json:"login"`
	Password string `json:"password"`

	// text
	TextData string `json:"text_data"`

	// binary
	BinaryData []byte `json:"binary_data"`

	// card
	CardNumber string `json:"card_number"`
	CardExp    string `json:"card_exp"`
	CardHolder string `json:"card_holder"`
}

type CreateSecretDTO struct {
	UserID     int64      `json:"user_id"`
	Title      string     `json:"title"`
	SecretType SecretType `json:"secret_type"`

	// custom key-value data
	Metadata string `json:"metadata"`

	// login_password
	Login    string `json:"login"`
	Password string `json:"password"`

	// text
	TextData string `json:"text_data"`

	// binary
	BinaryData []byte `json:"binary_data"`

	// card
	CardNumber string `json:"card_number"`
	CardExp    string `json:"card_exp"`
	CardHolder string `json:"card_holder"`
}

type GetSecretDTO struct {
	UserID int64  `json:"user_id"`
	Title  string `json:"title"`
}
