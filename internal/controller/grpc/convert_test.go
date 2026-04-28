package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/ibeloyar/gophkeeper/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	gophkeeperv1 "github.com/ibeloyar/gophkeeper/proto/gophkeeper/v1"
)

func TestConvertSecretTypeToProto(t *testing.T) {
	tests := []struct {
		name       string
		secretType model.SecretType
		expected   gophkeeperv1.SecretType
	}{
		{
			name:       "LoginPassword",
			secretType: model.SecretTypeLoginPassword,
			expected:   gophkeeperv1.SecretType_LOGIN_PASSWORD,
		},
		{
			name:       "Text",
			secretType: model.SecretTypeText,
			expected:   gophkeeperv1.SecretType_TEXT,
		},
		{
			name:       "Binary",
			secretType: model.SecretTypeBinary,
			expected:   gophkeeperv1.SecretType_BINARY,
		},
		{
			name:       "Card",
			secretType: model.SecretTypeCard,
			expected:   gophkeeperv1.SecretType_CARD,
		},
		{
			name:       "Unknown",
			secretType: "unknown",
			expected:   gophkeeperv1.SecretType_UNSPECIFIED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertSecretTypeToProto(tt.secretType)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestConvertSecretTypeToDTO(t *testing.T) {
	tests := []struct {
		name      string
		protoType gophkeeperv1.SecretType
		expected  model.SecretType
	}{
		{
			name:      "LoginPassword",
			protoType: gophkeeperv1.SecretType_LOGIN_PASSWORD,
			expected:  model.SecretTypeLoginPassword,
		},
		{
			name:      "Text",
			protoType: gophkeeperv1.SecretType_TEXT,
			expected:  model.SecretTypeText,
		},
		{
			name:      "Binary",
			protoType: gophkeeperv1.SecretType_BINARY,
			expected:  model.SecretTypeBinary,
		},
		{
			name:      "Card",
			protoType: gophkeeperv1.SecretType_CARD,
			expected:  model.SecretTypeCard,
		},
		{
			name:      "Unspecified",
			protoType: gophkeeperv1.SecretType_UNSPECIFIED,
			expected:  "",
		},
		{
			name:      "Unknown",
			protoType: gophkeeperv1.SecretType(-1),
			expected:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertSecretTypeToDTO(tt.protoType)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestConvertGetSecretToProto(t *testing.T) {
	now := time.Now().UTC()
	nowFormatted := now.Format(time.RFC3339)

	tests := []struct {
		name     string
		input    *model.Secret
		expected *gophkeeperv1.GetSecretResponse
	}{
		{
			name: "LoginPassword",
			input: &model.Secret{
				ID:         1,
				UserID:     123,
				Title:      "login-pass-secret",
				SecretType: model.SecretTypeLoginPassword,
				Metadata:   "meta-login-pass",
				Login:      "testuser",
				Password:   "testpass",
				CreatedAt:  now,
				UpdatedAt:  now,
				TextData:   "",
				BinaryData: nil,
				CardNumber: "",
				CardExp:    "",
				CardHolder: "",
			},
			expected: &gophkeeperv1.GetSecretResponse{
				Id:         1,
				UserId:     123,
				Title:      "login-pass-secret",
				SecretType: gophkeeperv1.SecretType_LOGIN_PASSWORD,
				Metadata:   "meta-login-pass",
				Login:      "testuser",
				Password:   "testpass",
				CreatedAt:  nowFormatted,
				UpdatedAt:  nowFormatted,
				TextData:   "",
				BinaryData: nil,
				CardNumber: "",
				CardExp:    "",
				CardHolder: "",
			},
		},
		{
			name: "Text",
			input: &model.Secret{
				ID:         2,
				UserID:     123,
				Title:      "text-secret",
				SecretType: model.SecretTypeText,
				Metadata:   "meta-text",
				TextData:   "text-data",
				CreatedAt:  now,
				UpdatedAt:  now,
			},
			expected: &gophkeeperv1.GetSecretResponse{
				Id:         2,
				UserId:     123,
				Title:      "text-secret",
				SecretType: gophkeeperv1.SecretType_TEXT,
				Metadata:   "meta-text",
				TextData:   "text-data",
				CreatedAt:  nowFormatted,
				UpdatedAt:  nowFormatted,
				Login:      "",
				Password:   "",
				BinaryData: nil,
				CardNumber: "",
				CardExp:    "",
				CardHolder: "",
			},
		},
		{
			name: "Binary",
			input: &model.Secret{
				ID:         3,
				UserID:     123,
				Title:      "binary-secret",
				SecretType: model.SecretTypeBinary,
				Metadata:   "meta-binary",
				BinaryData: []byte{1, 2, 3},
				CreatedAt:  now,
				UpdatedAt:  now,
			},
			expected: &gophkeeperv1.GetSecretResponse{
				Id:         3,
				UserId:     123,
				Title:      "binary-secret",
				SecretType: gophkeeperv1.SecretType_BINARY,
				Metadata:   "meta-binary",
				BinaryData: []byte{1, 2, 3},
				CreatedAt:  nowFormatted,
				UpdatedAt:  nowFormatted,
				Login:      "",
				Password:   "",
				TextData:   "",
				CardNumber: "",
				CardExp:    "",
				CardHolder: "",
			},
		},
		{
			name: "Card",
			input: &model.Secret{
				ID:         4,
				UserID:     123,
				Title:      "card-secret",
				SecretType: model.SecretTypeCard,
				Metadata:   "meta-card",
				CardNumber: "1234567812345678",
				CardExp:    "12/25",
				CardHolder: "John Doe",
				CreatedAt:  now,
				UpdatedAt:  now,
			},
			expected: &gophkeeperv1.GetSecretResponse{
				Id:         4,
				UserId:     123,
				Title:      "card-secret",
				SecretType: gophkeeperv1.SecretType_CARD,
				Metadata:   "meta-card",
				CardNumber: "1234567812345678",
				CardExp:    "12/25",
				CardHolder: "John Doe",
				CreatedAt:  nowFormatted,
				UpdatedAt:  nowFormatted,
				Login:      "",
				Password:   "",
				TextData:   "",
				BinaryData: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertGetSecretToProto(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestConvertGetSecretsToProto(t *testing.T) {
	now := time.Now().UTC()

	secrets := []*model.Secret{
		{
			ID:         1,
			UserID:     123,
			Title:      "text-secret",
			SecretType: model.SecretTypeText,
			Metadata:   "meta-text",
			TextData:   "text-data",
			CreatedAt:  now,
			UpdatedAt:  now.Add(time.Second),
		},
		{
			ID:         2,
			UserID:     123,
			Title:      "binary-secret",
			SecretType: model.SecretTypeBinary,
			Metadata:   "meta-binary",
			BinaryData: []byte{1, 2, 3},
			CreatedAt:  now.Add(2 * time.Second),
			UpdatedAt:  now.Add(3 * time.Second),
		},
	}

	got := convertGetSecretsToProto(secrets)

	assert.Len(t, got, 2)

	// text-secret
	assert.Equal(t, int64(1), got[0].Id)
	assert.Equal(t, int64(123), got[0].UserId)
	assert.Equal(t, "text-secret", got[0].Title)
	assert.Equal(t, gophkeeperv1.SecretType_TEXT, got[0].SecretType)
	assert.Equal(t, "meta-text", got[0].Metadata)
	assert.Equal(t, "text-data", got[0].TextData)
	assert.Equal(t, now.Format(time.RFC3339), got[0].CreatedAt)
	assert.Equal(t, now.Add(time.Second).Format(time.RFC3339), got[0].UpdatedAt)

	// binary-secret
	assert.Equal(t, int64(2), got[1].Id)
	assert.Equal(t, int64(123), got[1].UserId)
	assert.Equal(t, "binary-secret", got[1].Title)
	assert.Equal(t, gophkeeperv1.SecretType_BINARY, got[1].SecretType)
	assert.Equal(t, "meta-binary", got[1].Metadata)
	assert.Equal(t, []byte{1, 2, 3}, got[1].BinaryData)
	assert.Equal(t, now.Add(2*time.Second).Format(time.RFC3339), got[1].CreatedAt)
	assert.Equal(t, now.Add(3*time.Second).Format(time.RFC3339), got[1].UpdatedAt)
}

func TestConvertCreateSecretToDTO_TokenPresent(t *testing.T) {
	input := &gophkeeperv1.CreateSecretRequest{
		Title:      "test-title",
		SecretType: gophkeeperv1.SecretType_TEXT,
		Metadata:   "meta",
		Login:      "test-login",
		Password:   "test-pass",
		TextData:   "test-text",
		BinaryData: []byte("test-binary"),
		CardExp:    "12/25",
		CardNumber: "1234567812345678",
		CardHolder: "Test Holder",
	}

	dto, err := convertCreateSecretToDTO(context.Background(), input)
	require.NoError(t, err)
	require.NotNil(t, dto)

	assert.Equal(t, "test-title", dto.Title)
	assert.Equal(t, model.SecretType("text"), dto.SecretType)

	assert.Equal(t, "meta", dto.Metadata)
	assert.Equal(t, "test-login", dto.Login)
	assert.Equal(t, "test-pass", dto.Password)
	assert.Equal(t, "test-text", dto.TextData)
	assert.Equal(t, []byte("test-binary"), dto.BinaryData)
	assert.Equal(t, "12/25", dto.CardExp)
	assert.Equal(t, "1234567812345678", dto.CardNumber)
	assert.Equal(t, "Test Holder", dto.CardHolder)
}
