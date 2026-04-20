package model

import (
	"errors"
)

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return e.Message
}

const (
	ErrInternalServerMessage         = "internal server error"
	ErrInvalidLoginOrPasswordMessage = "invalid login or password"
	ErrUserAlreadyExistMessage       = "user already exists"
)

var (
	ErrServerInternal           = errors.New(ErrInternalServerMessage)
	ErrUserAlreadyExist         = errors.New(ErrUserAlreadyExistMessage)
	ErrInvalidLoginOrPassword   = errors.New(ErrInvalidLoginOrPasswordMessage)
	ErrListeningToLocalAddress  = errors.New("listening to local address error")
	ErrStartingGrpcServer       = errors.New("starting grpc server error")
	ErrConfigFileReading        = errors.New("config file read error")
	ErrTitleIsRequired          = errors.New("title is required")
	ErrSecretLoginRequired      = errors.New("login is required for secret type login_password")
	ErrSecretPasswordRequired   = errors.New("password is required for secret type login_password")
	ErrSecretLoginMaxLen        = errors.New("max length login for secret type login_password is 64")
	ErrSecretPasswordMaxLen     = errors.New("max length password for secret type login_password is 64")
	ErrSecretTextDataRequired   = errors.New("text_data is required for secret type text")
	ErrSecretBinaryRequired     = errors.New("binary_data is required for secret type binary")
	ErrSecretCardNumberRequired = errors.New("card_number is required for secret type card")
	ErrSecretCardNumberMaxLen   = errors.New("max length card_number for secret type card is 16")
	ErrSecretCardExpRequired    = errors.New("card_exp is required for secret type card")
	ErrSecretCardExpInvalid     = errors.New("invalid card_exp for secret type card")
	ErrSecretCardHolderRequired = errors.New("card_holder is required for secret type card")
	ErrSecretCardHolderMaxLen   = errors.New("max length card_holder for secret type card is 128")
	ErrSecretNotFound           = errors.New("secret not found")
)
