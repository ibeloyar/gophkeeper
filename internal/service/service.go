package service

import (
	"context"
	"strings"
	"time"

	"github.com/ibeloyar/gophkeeper/internal/model"
	"github.com/ibeloyar/gophkeeper/internal/repository/pgstorage"
	"github.com/ibeloyar/gophkeeper/pgk/auth"
	"github.com/ibeloyar/gophkeeper/pgk/password"
	"go.uber.org/zap"
)

type Storage interface {
	CreateUser(ctx context.Context, user model.User) (int64, error)
	GetUserByLogin(ctx context.Context, login string) *model.User

	CreateSecret(ctx context.Context, secret *model.CreateSecretDTO) (int64, error)
	GetSecret(ctx context.Context, title string, userID int64) (*model.Secret, error)
	GetSecrets(ctx context.Context, userID int64) ([]*model.Secret, error)
	DeleteSecret(ctx context.Context, title string, userID int64) error
}

type Service struct {
	lg      *zap.SugaredLogger
	storage Storage

	userPasswordCost  int
	secretPasswordKey string
	tokenSecret       string
	tokenExp          time.Duration
}

// New creates Service with all dependencies injected. Configures security parameters.
func New(lg *zap.SugaredLogger, storage Storage, userPasswordCost int, secretPasswordKey string, tokenExp time.Duration, tokenSecret string) *Service {
	return &Service{
		lg:                lg,
		storage:           storage,
		userPasswordCost:  userPasswordCost,
		secretPasswordKey: secretPasswordKey,
		tokenExp:          tokenExp,
		tokenSecret:       tokenSecret,
	}
}

// Register validates input, hashes password with BCrypt, creates user, generates JWT.
// Maps PostgreSQL unique constraint violation to ErrUserAlreadyExist.
// Returns Bearer token for successful registration.
func (s *Service) Register(ctx context.Context, input *model.RegisterDTO) (string, error) {
	if err := validateRegisterDTO(input); err != nil {
		return "", err
	}

	passwordHash, err := password.HashPassword(input.Password, s.userPasswordCost)
	if err != nil {
		return "", model.ErrServerInternal
	}

	userID, err := s.storage.CreateUser(ctx, model.User{
		Login:        input.Login,
		PasswordHash: passwordHash,
	})
	if err != nil {
		if strings.Contains(err.Error(), pgstorage.ErrIsExistCode) {
			return "", model.ErrUserAlreadyExist
		}
		return "", model.ErrServerInternal
	}

	token, err := auth.GenerateBearerToken(model.TokenInfo{
		ID:    userID,
		Login: input.Login,
	}, s.tokenExp, s.tokenSecret)
	if err != nil {
		return "", model.ErrServerInternal
	}

	return token, nil
}

// Login validates input, retrieves user, verifies BCrypt password hash, generates JWT.
// Returns ErrInvalidLoginOrPassword for wrong credentials or missing user.
func (s *Service) Login(ctx context.Context, input *model.LoginDTO) (string, error) {
	if err := validateLoginDTO(input); err != nil {
		return "", err
	}

	user := s.storage.GetUserByLogin(ctx, input.Login)
	if user == nil {
		return "", model.ErrInvalidLoginOrPassword
	}

	if !password.CheckPasswordHash(input.Password, user.PasswordHash) {
		return "", model.ErrInvalidLoginOrPassword
	}

	token, err := auth.GenerateBearerToken(model.TokenInfo{
		ID:    user.ID,
		Login: user.Login,
	}, s.tokenExp, s.tokenSecret)
	if err != nil {
		return "", model.ErrServerInternal
	}

	return token, nil
}

// CreateSecret validates input, encrypts password for LoginPassword secrets using AES-256-GCM.
// Delegates storage with encrypted data. Propagates validation/storage errors.
func (s *Service) CreateSecret(ctx context.Context, input *model.CreateSecretDTO) error {
	if err := validateCreateSecretDTO(input); err != nil {
		return err
	}

	if input.SecretType == model.SecretTypeLoginPassword {
		passwordHash, err := password.EncryptPassword(input.Password, s.secretPasswordKey)
		if err != nil {
			return model.ErrServerInternal
		}
		input.Password = passwordHash
	}

	if _, err := s.storage.CreateSecret(ctx, input); err != nil {
		return err
	}

	return nil
}

// GetSecret fetches secret from storage, decrypts password for LoginPassword type.
// Returns ErrTitleIsRequired for empty title, ErrSecretNotFound if not owned by user.
func (s *Service) GetSecret(ctx context.Context, title string, userID int64) (*model.Secret, error) {
	if title == "" {
		return nil, model.ErrTitleIsRequired
	}

	secret, err := s.storage.GetSecret(ctx, title, userID)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, model.ErrSecretNotFound
	}

	if secret.SecretType == model.SecretTypeLoginPassword {
		decryptedPassword, err := password.DecryptPassword(secret.Password, s.secretPasswordKey)
		if err != nil {
			return nil, model.ErrServerInternal
		}
		secret.Password = decryptedPassword
	}

	return secret, nil
}

// GetSecrets lists all user secrets with password decryption for LoginPassword type.
// Returns decrypted plaintext passwords in results.
func (s *Service) GetSecrets(ctx context.Context, userID int64) ([]*model.Secret, error) {
	secrets, err := s.storage.GetSecrets(ctx, userID)
	if err != nil {
		return nil, err
	}

	for _, secret := range secrets {
		if secret.SecretType == model.SecretTypeLoginPassword {
			decryptedPassword, err := password.DecryptPassword(secret.Password, s.secretPasswordKey)
			if err != nil {
				return nil, model.ErrServerInternal
			}
			secret.Password = decryptedPassword
		}
	}

	return secrets, nil
}

// DeleteSecret validates title, delegates to storage layer.
// Propagates ErrTitleIsRequired and storage errors (including ErrSecretNotFound).
func (s *Service) DeleteSecret(ctx context.Context, title string, userID int64) error {
	if title == "" {
		return model.ErrTitleIsRequired
	}

	return s.storage.DeleteSecret(ctx, title, userID)
}
