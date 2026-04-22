package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/ibeloyar/gophkeeper/internal/model"
	"github.com/ibeloyar/gophkeeper/internal/service"
	"github.com/ibeloyar/gophkeeper/pgk/auth"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	mockPG "github.com/ibeloyar/gophkeeper/internal/repository/pgstorage/mocks"
	gophkeeperv1 "github.com/ibeloyar/gophkeeper/proto/gophkeeper/v1"
)

const (
	testTokenSecret       = "test-token-secret"
	testUserPasswordCost  = 10
	testSecretPasswordKey = "test-secret-key"
	testTokenExp          = time.Hour * 24
)

func TestGetSecret_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	secret := &model.Secret{
		ID:         1,
		UserID:     1,
		Title:      "test-title",
		SecretType: model.SecretTypeText,
		TextData:   "test-data",
	}

	mockStorage.EXPECT().
		GetSecret(gomock.Any(), "test-title", int64(1)).
		Return(secret, nil)

	svc := service.New(
		zap.NewNop().Sugar(),
		mockStorage,
		testUserPasswordCost,
		testSecretPasswordKey,
		testTokenExp,
		testTokenSecret,
	)
	controller := New(zap.NewNop().Sugar(), svc)

	tokenInfo := model.TokenInfo{ID: 1, Login: "admin"}
	ctx := context.WithValue(context.Background(), auth.TokenDataContextKey, &tokenInfo)

	req := &gophkeeperv1.GetSecretRequest{Title: "test-title"}
	resp, err := controller.GetSecret(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
}
