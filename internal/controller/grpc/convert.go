package grpc

import (
	"context"
	"time"

	"github.com/ibeloyar/gophkeeper/internal/model"
	"github.com/ibeloyar/gophkeeper/pgk/auth"

	gophkeeperv1 "github.com/ibeloyar/gophkeeper/proto/gophkeeper/v1"
)

func convertSecretTypeToProto(secretType model.SecretType) gophkeeperv1.SecretType {
	switch secretType {
	case model.SecretTypeLoginPassword:
		return gophkeeperv1.SecretType_LOGIN_PASSWORD
	case model.SecretTypeText:
		return gophkeeperv1.SecretType_TEXT
	case model.SecretTypeBinary:
		return gophkeeperv1.SecretType_BINARY
	case model.SecretTypeCard:
		return gophkeeperv1.SecretType_CARD
	default:
		return gophkeeperv1.SecretType_UNSPECIFIED
	}
}

func convertSecretTypeToDTO(secretType gophkeeperv1.SecretType) model.SecretType {
	switch secretType {
	case gophkeeperv1.SecretType_LOGIN_PASSWORD:
		return model.SecretTypeLoginPassword
	case gophkeeperv1.SecretType_TEXT:
		return model.SecretTypeText
	case gophkeeperv1.SecretType_BINARY:
		return model.SecretTypeBinary
	case gophkeeperv1.SecretType_CARD:
		return model.SecretTypeCard
	default:
		return ""
	}
}

// convertCreateSecretToDTO transforms gRPC CreateSecretRequest to model.CreateSecretDTO.
// Extracts userID from JWT context. Populates all secret fields based on SecretType.
// Returns populated DTO with user context for service layer processing.
func convertCreateSecretToDTO(ctx context.Context, input *gophkeeperv1.CreateSecretRequest) (*model.CreateSecretDTO, error) {
	userID := int64(0)

	tokenInfo := auth.GetTokenInfoFromContext[model.TokenInfo](ctx)
	if tokenInfo != nil {
		userID = tokenInfo.ID
	}

	return &model.CreateSecretDTO{
		UserID:     userID,
		Title:      input.Title,
		SecretType: convertSecretTypeToDTO(input.SecretType),

		Metadata: input.Metadata,

		Login:    input.Login,
		Password: input.Password,

		TextData: input.TextData,

		BinaryData: input.BinaryData,

		CardExp:    input.CardExp,
		CardNumber: input.CardNumber,
		CardHolder: input.CardHolder,
	}, nil
}

// convertGetSecretToProto transforms model.Secret to gRPC GetSecretResponse.
// Formats timestamps as RFC3339 strings. Converts SecretType enum.
// Used in GetSecret RPC response.
func convertGetSecretToProto(input *model.Secret) *gophkeeperv1.GetSecretResponse {
	return &gophkeeperv1.GetSecretResponse{
		Id:         input.ID,
		UserId:     input.UserID,
		Title:      input.Title,
		UpdatedAt:  input.UpdatedAt.Format(time.RFC3339),
		CreatedAt:  input.CreatedAt.Format(time.RFC3339),
		SecretType: convertSecretTypeToProto(input.SecretType),

		Metadata: input.Metadata,

		Login:    input.Login,
		Password: input.Password,

		TextData: input.TextData,

		BinaryData: input.BinaryData,

		CardExp:    input.CardExp,
		CardNumber: input.CardNumber,
		CardHolder: input.CardHolder,
	}
}

// convertGetSecretsToProto transforms []*model.Secret slice to []*gophkeeperv1.Secret slice.
// Applies convertGetSecretToProto logic to each secret. Used in GetSecrets RPC
func convertGetSecretsToProto(secrets []*model.Secret) []*gophkeeperv1.Secret {
	result := make([]*gophkeeperv1.Secret, 0)

	for _, secret := range secrets {
		result = append(result, &gophkeeperv1.Secret{
			Id:         secret.ID,
			UserId:     secret.UserID,
			Title:      secret.Title,
			UpdatedAt:  secret.UpdatedAt.Format(time.RFC3339),
			CreatedAt:  secret.CreatedAt.Format(time.RFC3339),
			SecretType: convertSecretTypeToProto(secret.SecretType),
			Metadata:   secret.Metadata,
			Login:      secret.Login,
			Password:   secret.Password,
			TextData:   secret.TextData,
			BinaryData: secret.BinaryData,
			CardExp:    secret.CardExp,
			CardNumber: secret.CardNumber,
			CardHolder: secret.CardHolder,
		})
	}

	return result
}
