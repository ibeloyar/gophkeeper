package service

import (
	"regexp"

	"github.com/ibeloyar/gophkeeper/internal/model"
)

const (
	minPassLen  = 4
	maxPassLen  = 64
	minLoginLen = 3
	maxLoginLen = 64

	maxSecretPassLen  = 64
	maxSecretLoginLen = 64
)

func validateLoginDTO(input *model.LoginDTO) error {
	if err := validateLogin(input.Login); err != nil {
		return err
	}

	if err := validatePassword(input.Password); err != nil {
		return err
	}

	return nil
}

func validateRegisterDTO(input *model.RegisterDTO) error {
	if err := validateLogin(input.Login); err != nil {
		return err
	}

	if err := validatePassword(input.Password); err != nil {
		return err
	}

	return nil
}

func validateLogin(login string) error {
	if len(login) < minLoginLen || len(login) > maxLoginLen {
		return model.ErrInvalidLoginOrPassword
	}

	return nil
}

func validatePassword(password string) error {
	if len(password) < minPassLen || len(password) > maxPassLen {
		return model.ErrInvalidLoginOrPassword
	}

	return nil
}

func validateCreateSecretDTO(input *model.CreateSecretDTO) error {
	if input.Title == "" {
		return model.ErrTitleIsRequired
	}

	if input.SecretType == model.SecretTypeLoginPassword {
		if input.Login == "" {
			return model.ErrSecretLoginRequired
		}
		if len(input.Login) > maxSecretLoginLen {
			return model.ErrSecretLoginMaxLen
		}
		if input.Password == "" {
			return model.ErrSecretPasswordRequired
		}
		if len(input.Password) > maxSecretPassLen {
			return model.ErrSecretPasswordMaxLen
		}
	}

	if input.SecretType == model.SecretTypeText {
		if input.TextData == "" {
			return model.ErrSecretTextDataRequired
		}
	}

	if input.SecretType == model.SecretTypeBinary {
		if len(input.BinaryData) == 0 {
			return model.ErrSecretBinaryRequired
		}
	}

	if input.SecretType == model.SecretTypeCard {
		if input.CardNumber == "" {
			return model.ErrSecretCardNumberRequired
		}
		if len(input.CardNumber) > 16 {
			return model.ErrSecretCardNumberMaxLen
		}
		if input.CardExp == "" {
			return model.ErrSecretCardExpRequired
		}
		if len(input.CardExp) != 5 || !regexp.MustCompile(`^\d{2}/\d{2}$`).MatchString(input.CardExp) {
			return model.ErrSecretCardExpInvalid
		}
		if input.CardHolder == "" {
			return model.ErrSecretCardHolderRequired
		}
		if len(input.CardHolder) > 128 {
			return model.ErrSecretCardHolderMaxLen
		}
	}

	return nil
}
