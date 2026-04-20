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
	CreateUser(user model.User) (int64, error)
	GetUserByLogin(login string) *model.User

	CreateSecret(secret *model.CreateSecretDTO) (int64, error)
	GetSecret(title string, userID int64) (*model.Secret, error)
	GetSecrets(userID int64) ([]*model.Secret, error)
	DeleteSecret(title string, userID int64) error
}

type Service struct {
	lg      *zap.SugaredLogger
	storage Storage

	userPasswordCost  int
	secretPasswordKey string
	tokenSecret       string
	tokenExp          time.Duration
}

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

func (s *Service) Shutdown(ctx context.Context) error {
	return nil
}

func (s *Service) Register(_ context.Context, input *model.RegisterDTO) (string, error) {
	if err := validateRegisterDTO(input); err != nil {
		return "", err
	}

	passwordHash, err := password.HashPassword(input.Password, s.userPasswordCost)
	if err != nil {
		return "", model.ErrServerInternal
	}

	userID, err := s.storage.CreateUser(model.User{
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

func (s *Service) Login(_ context.Context, input *model.LoginDTO) (string, error) {
	if err := validateLoginDTO(input); err != nil {
		return "", err
	}

	user := s.storage.GetUserByLogin(input.Login)
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

	if _, err := s.storage.CreateSecret(input); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetSecret(ctx context.Context, title string, userID int64) (*model.Secret, error) {
	if title == "" {
		return nil, model.ErrTitleIsRequired
	}

	secret, err := s.storage.GetSecret(title, userID)
	if err != nil {
		return nil, err
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

func (s *Service) GetSecrets(ctx context.Context, userID int64) ([]*model.Secret, error) {
	secrets, err := s.storage.GetSecrets(userID)
	if err != nil {
		return nil, err
	}

	return secrets, nil
}

func (s *Service) DeleteSecret(ctx context.Context, title string, userID int64) error {
	if title == "" {
		return model.ErrTitleIsRequired
	}

	return s.storage.DeleteSecret(title, userID)
}
