package app

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ibeloyar/gophkeeper/internal/service"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	mockPG "github.com/ibeloyar/gophkeeper/internal/repository/pgstorage/mocks"
)

func TestBuildGRPCServer_Success(t *testing.T) {
	lg := zaptest.NewLogger(t).Sugar()
	defer lg.Sync()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)

	svc := service.New(
		lg,
		mockStorage,
		testUserPasswordCost,
		testSecretPasswordKey,
		testTokenExp,
		testTokenSecret,
	)

	_, err := buildGRPCServer(lg, svc, testTokenSecret)
	require.NoError(t, err)
}
