package service

import (
	"testing"

	"github.com/ibeloyar/gophkeeper/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestValidateLogin(t *testing.T) {
	tests := []struct {
		name    string
		login   string
		wantErr bool
	}{
		{"too short", "", true},
		{"min length", string(make([]byte, minLoginLen)), false},
		{"valid", "user123", false},
		{"max length", string(make([]byte, maxLoginLen)), false},
		{"too long", string(make([]byte, maxLoginLen+1)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateLogin(tt.login)
			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, model.ErrInvalidLoginOrPassword)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name    string
		pass    string
		wantErr bool
	}{
		{"too short", "", true},
		{"min length", string(make([]byte, minPassLen)), false},
		{"valid", "pass123", false},
		{"max length", string(make([]byte, maxPassLen)), false},
		{"too long", string(make([]byte, maxPassLen+1)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePassword(tt.pass)
			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, model.ErrInvalidLoginOrPassword)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateLoginDTO(t *testing.T) {
	tests := []struct {
		name    string
		input   *model.LoginDTO
		wantErr bool
	}{
		{"empty login", &model.LoginDTO{Password: "pass123"}, true},
		{"short login", &model.LoginDTO{Login: "a", Password: "pass123"}, true},
		{"long login", &model.LoginDTO{Login: string(make([]byte, maxLoginLen+1)), Password: "pass123"}, true},
		{"short password", &model.LoginDTO{Login: "user", Password: "12"}, true},
		{"valid", &model.LoginDTO{Login: "user123", Password: "pass123"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateLoginDTO(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, model.ErrInvalidLoginOrPassword)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateRegisterDTO(t *testing.T) {
	tests := []struct {
		name    string
		input   *model.RegisterDTO
		wantErr bool
	}{
		{"empty login", &model.RegisterDTO{Password: "pass123"}, true},
		{"short login", &model.RegisterDTO{Login: "a", Password: "pass123"}, true},
		{"long login", &model.RegisterDTO{Login: string(make([]byte, maxLoginLen+1)), Password: "pass123"}, true},
		{"short password", &model.RegisterDTO{Login: "user", Password: "12"}, true},
		{"valid", &model.RegisterDTO{Login: "user123", Password: "pass123"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRegisterDTO(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, model.ErrInvalidLoginOrPassword)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateCreateSecretDTO(t *testing.T) {
	tests := []struct {
		name    string
		input   *model.CreateSecretDTO
		wantErr bool
		errCode string
	}{
		{"empty title", &model.CreateSecretDTO{}, true, ""},

		{"loginpass empty login", &model.CreateSecretDTO{Title: "title", SecretType: model.SecretTypeLoginPassword}, true, "ErrSecretLoginRequired"},
		{"loginpass long login", &model.CreateSecretDTO{Title: "title", SecretType: model.SecretTypeLoginPassword, Login: string(make([]byte, maxSecretLoginLen+1))}, true, "ErrSecretLoginMaxLen"},
		{"loginpass empty password", &model.CreateSecretDTO{Title: "title", SecretType: model.SecretTypeLoginPassword, Login: "login"}, true, "ErrSecretPasswordRequired"},
		{"loginpass long password", &model.CreateSecretDTO{Title: "title", SecretType: model.SecretTypeLoginPassword, Login: "login", Password: string(make([]byte, maxSecretPassLen+1))}, true, "ErrSecretPasswordMaxLen"},
		{"loginpass valid", &model.CreateSecretDTO{Title: "title", SecretType: model.SecretTypeLoginPassword, Login: "login", Password: "pass123"}, false, ""},

		{"text empty data", &model.CreateSecretDTO{Title: "title", SecretType: model.SecretTypeText}, true, "ErrSecretTextDataRequired"},
		{"text valid", &model.CreateSecretDTO{Title: "title", SecretType: model.SecretTypeText, TextData: "data"}, false, ""},

		{"binary empty data", &model.CreateSecretDTO{Title: "title", SecretType: model.SecretTypeBinary}, true, "ErrSecretBinaryRequired"},
		{"binary valid", &model.CreateSecretDTO{Title: "title", SecretType: model.SecretTypeBinary, BinaryData: []byte("data")}, false, ""},

		{"card empty number", &model.CreateSecretDTO{Title: "title", SecretType: model.SecretTypeCard}, true, "ErrSecretCardNumberRequired"},
		{"card long number", &model.CreateSecretDTO{Title: "title", SecretType: model.SecretTypeCard, CardNumber: string(make([]byte, 17))}, true, "ErrSecretCardNumberMaxLen"},
		{"card empty exp", &model.CreateSecretDTO{Title: "title", SecretType: model.SecretTypeCard, CardNumber: "1234"}, true, "ErrSecretCardExpRequired"},
		{"card invalid exp len", &model.CreateSecretDTO{Title: "title", SecretType: model.SecretTypeCard, CardNumber: "1234", CardExp: "12/34/56"}, true, "ErrSecretCardExpInvalid"},
		{"card invalid exp format", &model.CreateSecretDTO{Title: "title", SecretType: model.SecretTypeCard, CardNumber: "1234", CardExp: "1234"}, true, "ErrSecretCardExpInvalid"},
		{"card empty holder", &model.CreateSecretDTO{Title: "title", SecretType: model.SecretTypeCard, CardNumber: "1234", CardExp: "12/34"}, true, "ErrSecretCardHolderRequired"},
		{"card long holder", &model.CreateSecretDTO{Title: "title", SecretType: model.SecretTypeCard, CardNumber: "1234567890123456", CardExp: "12/34", CardHolder: string(make([]byte, 129))}, true, "ErrSecretCardHolderMaxLen"},
		{"card valid", &model.CreateSecretDTO{Title: "title", SecretType: model.SecretTypeCard, CardNumber: "1234567890123456", CardExp: "12/34", CardHolder: "John Doe"}, false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCreateSecretDTO(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
