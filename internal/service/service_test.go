package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/ibeloyar/gophkeeper/internal/model"
	"github.com/ibeloyar/gophkeeper/internal/repository/pgstorage"
	"github.com/ibeloyar/gophkeeper/pkg/password"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	mockPG "github.com/ibeloyar/gophkeeper/internal/repository/pgstorage/mocks"
)

const (
	testUserPasswordCost  = 3
	testSecretPasswordKey = "key"
	testTokenSecret       = "secret"
	testTokenExp          = 5 * time.Minute
)

func newTestLogger(t *testing.T) *zap.SugaredLogger {
	logger := zaptest.NewLogger(t)
	return logger.Sugar()
}

func TestService_Register_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	svc := New(newTestLogger(t), mockStorage, testUserPasswordCost, testSecretPasswordKey, testTokenExp, testTokenSecret)

	input := &model.RegisterDTO{
		Login:    "testuser",
		Password: "testpass123",
	}

	mockStorage.EXPECT().
		CreateUser(context.Background(), gomock.Any()).
		Return(int64(123), nil).
		Times(1)

	token, err := svc.Register(context.Background(), input)

	assert.Nil(t, err)
	assert.NotEmpty(t, token)
}

func TestService_Register_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	svc := New(newTestLogger(t), mockStorage, testUserPasswordCost, testSecretPasswordKey, testTokenExp, testTokenSecret)

	input := &model.RegisterDTO{Login: "test", Password: "1"}
	ctx := context.Background()

	mockStorage.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)

	token, err := svc.Register(ctx, input)
	assert.Error(t, err)
	assert.Empty(t, token)
}

func TestService_Register_PasswordHashError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	svc := New(newTestLogger(t), mockStorage, 99, testSecretPasswordKey, testTokenExp, testTokenSecret)

	input := &model.RegisterDTO{Login: "testuser", Password: "pass"}
	ctx := context.Background()

	mockStorage.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)

	token, err := svc.Register(ctx, input)
	assert.ErrorIs(t, err, model.ErrServerInternal)
	assert.Empty(t, token)
}

func TestService_Register_StorageError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	svc := New(newTestLogger(t), mockStorage, testUserPasswordCost, testSecretPasswordKey, testTokenExp, testTokenSecret)

	input := &model.RegisterDTO{Login: "testuser", Password: "pass123"}
	ctx := context.Background()

	mockStorage.EXPECT().
		CreateUser(ctx, gomock.Any()).
		Return(int64(0), errors.New("db error")).
		Times(1)

	token, err := svc.Register(ctx, input)
	assert.ErrorIs(t, err, model.ErrServerInternal)
	assert.Empty(t, token)
}

func TestService_Register_UserExistsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	svc := New(newTestLogger(t), mockStorage, testUserPasswordCost, testSecretPasswordKey, testTokenExp, testTokenSecret)

	input := &model.RegisterDTO{Login: "testuser", Password: "pass123"}
	ctx := context.Background()

	mockStorage.EXPECT().
		CreateUser(ctx, gomock.Any()).
		Return(int64(0), fmt.Errorf("%s: duplicate key", pgstorage.ErrIsExistCode)).
		Times(1)

	token, err := svc.Register(ctx, input)
	assert.ErrorIs(t, err, model.ErrUserAlreadyExist)
	assert.Empty(t, token)
}

func TestService_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	svc := New(newTestLogger(t), mockStorage, testUserPasswordCost, testSecretPasswordKey, testTokenExp, testTokenSecret)

	input := &model.LoginDTO{Login: "testuser", Password: "testpass123"}
	ctx := context.Background()

	hash, _ := password.HashPassword(input.Password, testUserPasswordCost)
	expectedUser := &model.User{
		ID:           123,
		Login:        input.Login,
		PasswordHash: hash,
	}

	mockStorage.EXPECT().
		GetUserByLogin(ctx, input.Login).
		Return(expectedUser).
		Times(1)

	token, err := svc.Login(ctx, input)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestService_Login_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	svc := New(newTestLogger(t), mockStorage, testUserPasswordCost, testSecretPasswordKey, testTokenExp, testTokenSecret)

	input := &model.LoginDTO{Login: "te", Password: "12"}
	ctx := context.Background()

	mockStorage.EXPECT().GetUserByLogin(gomock.Any(), gomock.Any()).Times(0)

	_, err := svc.Login(ctx, input)
	assert.Error(t, err)
}

func TestService_Login_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	svc := New(newTestLogger(t), mockStorage, testUserPasswordCost, testSecretPasswordKey, testTokenExp, testTokenSecret)

	input := &model.LoginDTO{Login: "nonexistent", Password: "pass123"}
	ctx := context.Background()

	mockStorage.EXPECT().
		GetUserByLogin(ctx, input.Login).
		Return(nil).
		Times(1)

	token, err := svc.Login(ctx, input)
	assert.ErrorIs(t, err, model.ErrInvalidLoginOrPassword)
	assert.Empty(t, token)
}

func TestService_Login_WrongPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	svc := New(newTestLogger(t), mockStorage, testUserPasswordCost, testSecretPasswordKey, testTokenExp, testTokenSecret)

	input := &model.LoginDTO{Login: "testuser", Password: "wrongpass"}
	ctx := context.Background()

	wrongHash, _ := password.HashPassword("totallywrong", testUserPasswordCost)
	user := &model.User{
		ID:           123,
		Login:        input.Login,
		PasswordHash: wrongHash,
	}

	mockStorage.EXPECT().
		GetUserByLogin(ctx, input.Login).
		Return(user).
		Times(1)

	token, err := svc.Login(ctx, input)
	assert.ErrorIs(t, err, model.ErrInvalidLoginOrPassword)
	assert.Empty(t, token)
}

func TestService_CreateSecret_SuccessText(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	svc := New(newTestLogger(t), mockStorage, testUserPasswordCost, testSecretPasswordKey, testTokenExp, testTokenSecret)

	input := &model.CreateSecretDTO{
		Title:      "title",
		UserID:     1,
		SecretType: model.SecretTypeText,
		TextData:   "data",
		Metadata:   "meta",
	}
	ctx := context.Background()

	mockStorage.EXPECT().
		CreateSecret(ctx, gomock.Any()).
		Return(int64(1), nil)

	err := svc.CreateSecret(ctx, input)
	assert.NoError(t, err)
}

func TestService_CreateSecret_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	svc := New(newTestLogger(t), mockStorage, testUserPasswordCost, testSecretPasswordKey, testTokenExp, testTokenSecret)

	input := &model.CreateSecretDTO{}
	ctx := context.Background()

	mockStorage.EXPECT().CreateSecret(gomock.Any(), gomock.Any()).Times(0)

	err := svc.CreateSecret(ctx, input)
	assert.Error(t, err)
}

func TestService_CreateSecret_StorageError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	svc := New(newTestLogger(t), mockStorage, testUserPasswordCost, testSecretPasswordKey, testTokenExp, testTokenSecret)

	input := &model.CreateSecretDTO{
		Title:      "title",
		UserID:     1,
		SecretType: model.SecretTypeText,
		TextData:   "data",
		Metadata:   "meta",
	}
	ctx := context.Background()

	mockStorage.EXPECT().
		CreateSecret(ctx, gomock.Any()).
		Return(int64(0), errors.New("db error"))

	err := svc.CreateSecret(ctx, input)
	assert.Error(t, err)
}

func TestService_CreateSecret_EncryptPasswordError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	svc := New(newTestLogger(t), mockStorage, testUserPasswordCost, "", testTokenExp, testTokenSecret)

	input := &model.CreateSecretDTO{
		Title:      "title",
		UserID:     1,
		SecretType: model.SecretTypeLoginPassword,
		Login:      "login",
		Password:   "secretpass",
		Metadata:   "meta",
	}
	ctx := context.Background()

	mockStorage.EXPECT().CreateSecret(gomock.Any(), gomock.Any()).Times(0)

	err := svc.CreateSecret(ctx, input)
	assert.ErrorIs(t, err, model.ErrServerInternal)
}

func TestService_GetSecret_SuccessText(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	svc := New(newTestLogger(t), mockStorage, testUserPasswordCost, testSecretPasswordKey, testTokenExp, testTokenSecret)

	title := "secret1"
	userID := int64(1)
	ctx := context.Background()

	expectedSecret := &model.Secret{
		ID:         123,
		Title:      title,
		UserID:     userID,
		SecretType: model.SecretTypeText,
		TextData:   "plain data",
		Metadata:   "meta",
	}

	mockStorage.EXPECT().
		GetSecret(ctx, title, userID).
		Return(expectedSecret, nil)

	secret, err := svc.GetSecret(ctx, title, userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedSecret, secret)
}

func TestService_GetSecret_EmptyTitle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	svc := New(newTestLogger(t), mockStorage, testUserPasswordCost, testSecretPasswordKey, testTokenExp, testTokenSecret)

	ctx := context.Background()

	mockStorage.EXPECT().GetSecret(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	secret, err := svc.GetSecret(ctx, "", 1)
	assert.ErrorIs(t, err, model.ErrTitleIsRequired)
	assert.Nil(t, secret)
}

func TestService_GetSecret_StorageError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	svc := New(newTestLogger(t), mockStorage, testUserPasswordCost, testSecretPasswordKey, testTokenExp, testTokenSecret)

	title := "secret1"
	userID := int64(1)
	ctx := context.Background()

	mockStorage.EXPECT().
		GetSecret(ctx, title, userID).
		Return(nil, errors.New("not found"))

	secret, err := svc.GetSecret(ctx, title, userID)
	assert.Error(t, err)
	assert.Nil(t, secret)
}

func TestService_GetSecret_DecryptPasswordError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	svc := New(newTestLogger(t), mockStorage, testUserPasswordCost, "invalidkey", testTokenExp, testTokenSecret)

	title := "secret1"
	userID := int64(1)
	ctx := context.Background()

	secret := &model.Secret{
		ID:         123,
		Title:      title,
		UserID:     userID,
		SecretType: model.SecretTypeLoginPassword,
		Password:   "invalidencrypted",
	}

	mockStorage.EXPECT().
		GetSecret(ctx, title, userID).
		Return(secret, nil)

	result, err := svc.GetSecret(ctx, title, userID)
	assert.ErrorIs(t, err, model.ErrServerInternal)
	assert.Nil(t, result)
}

func TestService_GetSecrets_SuccessEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	svc := New(newTestLogger(t), mockStorage, testUserPasswordCost, testSecretPasswordKey, testTokenExp, testTokenSecret)

	userID := int64(1)
	ctx := context.Background()

	mockStorage.EXPECT().
		GetSecrets(ctx, userID).
		Return([]*model.Secret{}, nil)

	secrets, err := svc.GetSecrets(ctx, userID)
	assert.NoError(t, err)
	assert.Empty(t, secrets)
}

func TestService_GetSecrets_SuccessTextOnly(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	svc := New(newTestLogger(t), mockStorage, testUserPasswordCost, testSecretPasswordKey, testTokenExp, testTokenSecret)

	userID := int64(1)
	ctx := context.Background()

	textSecret := &model.Secret{
		ID:         1,
		Title:      "text secret",
		SecretType: model.SecretTypeText,
		TextData:   "plain",
	}

	mockStorage.EXPECT().
		GetSecrets(ctx, userID).
		Return([]*model.Secret{textSecret}, nil)

	secrets, err := svc.GetSecrets(ctx, userID)
	assert.NoError(t, err)
	assert.Len(t, secrets, 1)
	assert.Equal(t, textSecret, secrets[0])
}

func TestService_GetSecrets_StorageError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	svc := New(newTestLogger(t), mockStorage, testUserPasswordCost, testSecretPasswordKey, testTokenExp, testTokenSecret)

	userID := int64(1)
	ctx := context.Background()

	mockStorage.EXPECT().
		GetSecrets(ctx, userID).
		Return(nil, errors.New("db error"))

	secrets, err := svc.GetSecrets(ctx, userID)
	assert.Error(t, err)
	assert.Nil(t, secrets)
}

func TestService_GetSecrets_DecryptError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	svc := New(newTestLogger(t), mockStorage, testUserPasswordCost, "invalidkey", testTokenExp, testTokenSecret)

	userID := int64(1)
	ctx := context.Background()

	badSecret := &model.Secret{
		SecretType: model.SecretTypeLoginPassword,
		Password:   "invalidencrypted",
	}

	mockStorage.EXPECT().
		GetSecrets(ctx, userID).
		Return([]*model.Secret{badSecret}, nil)

	secrets, err := svc.GetSecrets(ctx, userID)
	assert.ErrorIs(t, err, model.ErrServerInternal)
	assert.Nil(t, secrets)
}

func TestService_DeleteSecret_EmptyTitle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	svc := New(newTestLogger(t), mockStorage, testUserPasswordCost, testSecretPasswordKey, testTokenExp, testTokenSecret)

	ctx := context.Background()

	mockStorage.EXPECT().DeleteSecret(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	err := svc.DeleteSecret(ctx, "", 1)
	assert.ErrorIs(t, err, model.ErrTitleIsRequired)
}

func TestService_DeleteSecret_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	svc := New(newTestLogger(t), mockStorage, testUserPasswordCost, testSecretPasswordKey, testTokenExp, testTokenSecret)

	title := "secret1"
	userID := int64(1)
	ctx := context.Background()

	mockStorage.EXPECT().
		DeleteSecret(ctx, title, userID).
		Return(nil)

	err := svc.DeleteSecret(ctx, title, userID)
	assert.NoError(t, err)
}

func TestService_DeleteSecret_StorageError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	svc := New(newTestLogger(t), mockStorage, testUserPasswordCost, testSecretPasswordKey, testTokenExp, testTokenSecret)

	title := "secret1"
	userID := int64(1)
	ctx := context.Background()

	mockStorage.EXPECT().
		DeleteSecret(ctx, title, userID).
		Return(errors.New("not found"))

	err := svc.DeleteSecret(ctx, title, userID)
	assert.Error(t, err)
}
