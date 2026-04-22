package app

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/ibeloyar/gophkeeper/internal/config"
	"github.com/ibeloyar/gophkeeper/internal/service"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	mockPG "github.com/ibeloyar/gophkeeper/internal/repository/pgstorage/mocks"
)

const (
	testTokenSecret       = "test-token-secret"
	testUserPasswordCost  = 10
	testSecretPasswordKey = "test-secret-key"
	testTokenExp          = time.Hour * 24
)

func TestBuildHTTPServer_Success(t *testing.T) {
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

	cfg, _ := config.Read("../../config/server/config.yaml")

	_, err := buildHTTPServer(lg, svc, cfg)
	require.NoError(t, err)
}
