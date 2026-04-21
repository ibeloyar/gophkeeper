package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/ibeloyar/gophkeeper/internal/model"
	"github.com/ibeloyar/gophkeeper/internal/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	mockPG "github.com/ibeloyar/gophkeeper/internal/repository/pgstorage/mocks"
)

const (
	testUserPasswordCost  = 3
	testSecretPasswordKey = "key"
	testTokenSecret       = "secret"
	testTokenExp          = 5 * time.Minute
)

func TestController_Register_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	svc := service.New(zap.NewNop().Sugar(), mockStorage, testUserPasswordCost, testSecretPasswordKey, testTokenExp, testTokenSecret)

	mockStorage.EXPECT().CreateUser(context.Background(), gomock.Any()).Return(int64(0), nil).Times(1)

	c := New(zap.NewNop().Sugar(), svc)

	router := http.NewServeMux()
	router.HandleFunc("/api/v1/register", c.Register)

	input := model.RegisterDTO{
		Login:    "testuser",
		Password: "testpass",
	}

	body, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/api/v1/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.NotEmpty(t, rr.Header().Get("Authorization"))
}

func TestController_Register_EmptyFields(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	svc := service.New(zap.NewNop().Sugar(), mockStorage, testUserPasswordCost, testSecretPasswordKey, testTokenExp, testTokenSecret)

	c := New(zap.NewNop().Sugar(), svc)

	router := http.NewServeMux()
	router.HandleFunc("/api/v1/register", c.Register)

	// empty login
	input1 := model.RegisterDTO{Login: "", Password: "testpass"}
	body1, _ := json.Marshal(input1)
	req1 := httptest.NewRequest("POST", "/api/v1/register", bytes.NewReader(body1))
	req1.Header.Set("Content-Type", "application/json")

	rr1 := httptest.NewRecorder()
	router.ServeHTTP(rr1, req1)

	assert.Equal(t, http.StatusBadRequest, rr1.Code)

	// empty password
	input2 := model.RegisterDTO{Login: "testuser", Password: ""}
	body2, _ := json.Marshal(input2)
	req2 := httptest.NewRequest("POST", "/api/v1/register", bytes.NewReader(body2))
	req2.Header.Set("Content-Type", "application/json")

	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	assert.Equal(t, http.StatusBadRequest, rr2.Code)

	// all empty
	input3 := model.RegisterDTO{Login: "", Password: ""}
	body3, _ := json.Marshal(input3)
	req3 := httptest.NewRequest("POST", "/api/v1/register", bytes.NewReader(body3))
	req3.Header.Set("Content-Type", "application/json")

	rr3 := httptest.NewRecorder()
	router.ServeHTTP(rr3, req3)

	assert.Equal(t, http.StatusBadRequest, rr3.Code)
}

func TestController_Register_InvalidLoginOrPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	svc := service.New(zap.NewNop().Sugar(), mockStorage, testUserPasswordCost, testSecretPasswordKey, testTokenExp, testTokenSecret)

	c := New(zap.NewNop().Sugar(), svc)

	router := http.NewServeMux()
	router.HandleFunc("/api/v1/register", c.Register)

	input := model.RegisterDTO{
		Login:    "d",
		Password: "d",
	}
	body, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/api/v1/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestController_Register_ServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	svc := service.New(zap.NewNop().Sugar(), mockStorage, testUserPasswordCost, testSecretPasswordKey, testTokenExp, testTokenSecret)

	mockStorage.EXPECT().CreateUser(context.Background(), gomock.Any()).Return(int64(0), model.ErrServerInternal)

	c := New(zap.NewNop().Sugar(), svc)

	router := http.NewServeMux()
	router.HandleFunc("/api/v1/register", c.Register)

	body, _ := json.Marshal(model.RegisterDTO{Login: "test", Password: "test"})
	req := httptest.NewRequest("POST", "/api/v1/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestController_Register_BodyReadError(t *testing.T) {
	core, logs := observer.New(zap.InfoLevel)
	logger := zap.New(core).Sugar()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	svc := service.New(logger, mockStorage, testUserPasswordCost, testSecretPasswordKey, testTokenExp, testTokenSecret)

	c := New(logger, svc)

	router := http.NewServeMux()
	router.HandleFunc("/api/v1/register", c.Register)

	req := httptest.NewRequest("POST", "/api/v1/register", bytes.NewReader([]byte("not-json")))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	assert.True(t, logs.Len() > 0)
	assert.Contains(t, logs.All()[0].Message, "failed to parse request body")
}

func TestController_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	svc := service.New(zap.NewNop().Sugar(), mockStorage, testUserPasswordCost, testSecretPasswordKey, testTokenExp, testTokenSecret)

	mockStorage.EXPECT().GetUserByLogin(context.Background(), "admin").Return(&model.User{
		ID:           1,
		Login:        "admin",
		PasswordHash: "$2a$10$e.Q6kFnSA591Gxi4tfx/LuyS7.NjEpFRLDvrnmuqHNILxgfHOpdvi",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})

	c := New(zap.NewNop().Sugar(), svc)

	router := http.NewServeMux()
	router.HandleFunc("/api/v1/login", c.Login)

	input := model.LoginDTO{
		Login:    "admin",
		Password: "admin",
	}
	body, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.NotEmpty(t, rr.Header().Get("Authorization"))
}

func TestController_Login_BadLoginOrPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	svc := service.New(zap.NewNop().Sugar(), mockStorage, testUserPasswordCost, testSecretPasswordKey, testTokenExp, testTokenSecret)

	mockStorage.EXPECT().GetUserByLogin(context.Background(), "baduser").Return(nil)

	c := New(zap.NewNop().Sugar(), svc)

	router := http.NewServeMux()
	router.HandleFunc("/api/v1/login", c.Login)

	input := model.LoginDTO{
		Login:    "baduser",
		Password: "badpass",
	}
	body, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "invalid login or password\n", rr.Body.String())
}

func TestController_Login_BodyReadError(t *testing.T) {
	core, logs := observer.New(zap.InfoLevel)
	logger := zap.New(core).Sugar()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	svc := service.New(logger, mockStorage, testUserPasswordCost, testSecretPasswordKey, testTokenExp, testTokenSecret)

	c := New(logger, svc)

	router := http.NewServeMux()
	router.HandleFunc("/api/v1/login", c.Login)

	req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewReader([]byte("not-json")))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	assert.True(t, logs.Len() > 0)
	assert.Contains(t, logs.All()[0].Message, "failed to parse request body")
}

func TestController_Login_EmptyFields(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	svc := service.New(zap.NewNop().Sugar(), mockStorage, testUserPasswordCost, testSecretPasswordKey, testTokenExp, testTokenSecret)

	c := New(zap.NewNop().Sugar(), svc)

	router := http.NewServeMux()
	router.HandleFunc("/api/v1/login", c.Login)

	// empty login
	input1 := model.LoginDTO{Login: "", Password: "testpass"}
	body1, _ := json.Marshal(input1)
	req1 := httptest.NewRequest("POST", "/api/v1/login", bytes.NewReader(body1))
	req1.Header.Set("Content-Type", "application/json")

	rr1 := httptest.NewRecorder()
	router.ServeHTTP(rr1, req1)

	assert.Equal(t, http.StatusBadRequest, rr1.Code)

	// empty password
	input2 := model.LoginDTO{Login: "testuser", Password: ""}
	body2, _ := json.Marshal(input2)
	req2 := httptest.NewRequest("POST", "/api/v1/login", bytes.NewReader(body2))
	req2.Header.Set("Content-Type", "application/json")

	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	assert.Equal(t, http.StatusBadRequest, rr2.Code)
}
