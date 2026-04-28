package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/ibeloyar/gophkeeper/internal/model"
	"github.com/ibeloyar/gophkeeper/internal/repository/pgstorage"
	"github.com/ibeloyar/gophkeeper/internal/service"
	"github.com/ibeloyar/gophkeeper/pkg/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestController_Register_UserAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	svc := service.New(zap.NewNop().Sugar(), mockStorage, testUserPasswordCost, testSecretPasswordKey, testTokenExp, testTokenSecret)

	mockStorage.EXPECT().CreateUser(context.Background(), gomock.Any()).Return(int64(0), errors.New(pgstorage.ErrIsExistCode)).Times(1)

	c := New(zap.NewNop().Sugar(), svc)

	router := http.NewServeMux()
	router.HandleFunc("/api/v1/register", c.Register)

	body, _ := json.Marshal(model.RegisterDTO{Login: "test", Password: "test"})
	req := httptest.NewRequest("POST", "/api/v1/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusConflict, rr.Code)
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

func TestGetSecret_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	mockStorage.EXPECT().
		GetSecret(context.Background(), "test-title", int64(1)).
		Return(&model.Secret{
			ID:         1,
			UserID:     1,
			Title:      "test-title",
			SecretType: model.SecretTypeText,
			TextData:   "test-text",
		}, nil)

	svc := service.New(
		zap.NewNop().Sugar(),
		mockStorage,
		testUserPasswordCost,
		testSecretPasswordKey,
		testTokenExp,
		testTokenSecret,
	)
	controller := New(zap.NewNop().Sugar(), svc)

	handler := http.HandlerFunc(controller.GetSecret)

	authMW := auth.AuthBearerMiddlewareInit[model.TokenInfo](testTokenSecret)
	authHandler := authMW(handler)

	body := model.GetSecretBody{Title: "test-title"}
	bodyBytes, _ := json.Marshal(body)

	tokenInfo := model.TokenInfo{
		ID:    1,
		Login: "admin",
	}
	token, _ := auth.GenerateBearerToken[model.TokenInfo](
		tokenInfo,
		testTokenExp,
		testTokenSecret,
	)

	req := httptest.NewRequest("POST", "/api/v1/get-secret", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	var resp model.Secret
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, "test-title", resp.Title)
	assert.Equal(t, int64(1), resp.UserID)
	assert.Equal(t, "test-text", resp.TextData)
}

func TestGetSecret_InvalidBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	svc := service.New(
		zap.NewNop().Sugar(),
		mockStorage,
		testUserPasswordCost,
		testSecretPasswordKey,
		testTokenExp,
		testTokenSecret,
	)
	controller := New(zap.NewNop().Sugar(), svc)

	handler := http.HandlerFunc(controller.GetSecret)
	authMW := auth.AuthBearerMiddlewareInit[model.TokenInfo](testTokenSecret)
	authHandler := authMW(handler)

	req := httptest.NewRequest("POST", "/api/v1/get-secret", bytes.NewReader([]byte(`{???}`)))
	req.Header.Set("Content-Type", "application/json")

	tokenInfo := model.TokenInfo{
		ID:    1,
		Login: "admin",
	}
	token, _ := auth.GenerateBearerToken[model.TokenInfo](tokenInfo, testTokenExp, testTokenSecret)
	req.Header.Set("Authorization", token)

	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetSecret_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	mockStorage.EXPECT().
		GetSecret(context.Background(), "not-found", int64(1)).
		Return(nil, model.ErrSecretNotFound)

	svc := service.New(
		zap.NewNop().Sugar(),
		mockStorage,
		testUserPasswordCost,
		testSecretPasswordKey,
		testTokenExp,
		testTokenSecret,
	)
	controller := New(zap.NewNop().Sugar(), svc)

	handler := http.HandlerFunc(controller.GetSecret)
	authMW := auth.AuthBearerMiddlewareInit[model.TokenInfo](testTokenSecret)
	authHandler := authMW(handler)

	body := model.GetSecretBody{Title: "not-found"}
	bodyBytes, _ := json.Marshal(body)

	tokenInfo := model.TokenInfo{
		ID:    1,
		Login: "admin",
	}
	token, _ := auth.GenerateBearerToken[model.TokenInfo](
		tokenInfo,
		testTokenExp,
		testTokenSecret,
	)

	req := httptest.NewRequest("POST", "/api/v1/get-secret", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetSecret_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	mockStorage.EXPECT().
		GetSecret(context.Background(), "test-title", int64(1)).
		Return(nil, errors.New("internal storage error"))

	svc := service.New(
		zap.NewNop().Sugar(),
		mockStorage,
		testUserPasswordCost,
		testSecretPasswordKey,
		testTokenExp,
		testTokenSecret,
	)
	controller := New(zap.NewNop().Sugar(), svc)

	handler := http.HandlerFunc(controller.GetSecret)
	authMW := auth.AuthBearerMiddlewareInit[model.TokenInfo](testTokenSecret)
	authHandler := authMW(handler)

	body := model.GetSecretBody{Title: "test-title"}
	bodyBytes, _ := json.Marshal(body)

	tokenInfo := model.TokenInfo{
		ID:    1,
		Login: "admin",
	}
	token, _ := auth.GenerateBearerToken[model.TokenInfo](
		tokenInfo,
		testTokenExp,
		testTokenSecret,
	)

	req := httptest.NewRequest("POST", "/api/v1/get-secret", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetSecrets_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	expectedSecrets := []*model.Secret{
		{
			ID:         1,
			UserID:     1,
			Title:      "test-secret-1",
			SecretType: model.SecretTypeText,
			TextData:   "data1",
		},
		{
			ID:         2,
			UserID:     1,
			Title:      "test-secret-2",
			SecretType: model.SecretTypeText,
			TextData:   "data2",
		},
	}
	mockStorage.EXPECT().
		GetSecrets(context.Background(), int64(1)).
		Return(expectedSecrets, nil)

	svc := service.New(
		zap.NewNop().Sugar(),
		mockStorage,
		testUserPasswordCost,
		testSecretPasswordKey,
		testTokenExp,
		testTokenSecret,
	)
	controller := New(zap.NewNop().Sugar(), svc)

	handler := http.HandlerFunc(controller.GetSecrets)
	authMW := auth.AuthBearerMiddlewareInit[model.TokenInfo](testTokenSecret)
	authHandler := authMW(handler)

	tokenInfo := model.TokenInfo{
		ID:    1,
		Login: "admin",
	}
	token, _ := auth.GenerateBearerToken[model.TokenInfo](tokenInfo, testTokenExp, testTokenSecret)

	req := httptest.NewRequest("GET", "/api/v1/get-secrets", nil)
	req.Header.Set("Authorization", token)

	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	var resp []*model.Secret
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.Equal(t, "test-secret-1", resp[0].Title)
	assert.Equal(t, "test-secret-2", resp[1].Title)
}

func TestGetSecrets_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	mockStorage.EXPECT().
		GetSecrets(context.Background(), int64(1)).
		Return(nil, errors.New("internal storage error"))

	svc := service.New(
		zap.NewNop().Sugar(),
		mockStorage,
		testUserPasswordCost,
		testSecretPasswordKey,
		testTokenExp,
		testTokenSecret,
	)
	controller := New(zap.NewNop().Sugar(), svc)

	handler := http.HandlerFunc(controller.GetSecrets)
	authMW := auth.AuthBearerMiddlewareInit[model.TokenInfo](testTokenSecret)
	authHandler := authMW(handler)

	tokenInfo := model.TokenInfo{
		ID:    1,
		Login: "admin",
	}
	token, _ := auth.GenerateBearerToken[model.TokenInfo](tokenInfo, testTokenExp, testTokenSecret)

	req := httptest.NewRequest("GET", "/api/v1/get-secrets", nil)
	req.Header.Set("Authorization", token)

	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateSecret_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	mockStorage.EXPECT().
		CreateSecret(gomock.Any(), gomock.Any()).
		Return(int64(1), nil)

	svc := service.New(
		zap.NewNop().Sugar(),
		mockStorage,
		testUserPasswordCost,
		testSecretPasswordKey,
		testTokenExp,
		testTokenSecret,
	)
	controller := New(zap.NewNop().Sugar(), svc)

	handler := http.HandlerFunc(controller.CreateSecret)
	authMW := auth.AuthBearerMiddlewareInit[model.TokenInfo](testTokenSecret)
	authHandler := authMW(handler)

	body := model.CreateSecretDTO{
		Title:      "test-secret",
		SecretType: model.SecretTypeText,
		TextData:   "test-data",
	}
	bodyBytes, _ := json.Marshal(body)

	tokenInfo := model.TokenInfo{
		ID:    1,
		Login: "admin",
	}
	token, _ := auth.GenerateBearerToken[model.TokenInfo](tokenInfo, testTokenExp, testTokenSecret)

	req := httptest.NewRequest("POST", "/api/v1/secrets", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestCreateSecret_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	svc := service.New(
		zap.NewNop().Sugar(),
		mockStorage,
		testUserPasswordCost,
		testSecretPasswordKey,
		testTokenExp,
		testTokenSecret,
	)
	controller := New(zap.NewNop().Sugar(), svc)

	handler := http.HandlerFunc(controller.CreateSecret)
	authMW := auth.AuthBearerMiddlewareInit[model.TokenInfo](testTokenSecret)
	authHandler := authMW(handler)

	bodyBytes := []byte(`{invalid json}`)

	tokenInfo := model.TokenInfo{
		ID:    1,
		Login: "admin",
	}
	token, _ := auth.GenerateBearerToken[model.TokenInfo](tokenInfo, testTokenExp, testTokenSecret)

	req := httptest.NewRequest("POST", "/api/v1/secrets", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateSecret_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	mockStorage.EXPECT().
		CreateSecret(gomock.Any(), gomock.Any()).
		Return(int64(0), errors.New("internal storage error"))

	svc := service.New(
		zap.NewNop().Sugar(),
		mockStorage,
		testUserPasswordCost,
		testSecretPasswordKey,
		testTokenExp,
		testTokenSecret,
	)
	controller := New(zap.NewNop().Sugar(), svc)

	handler := http.HandlerFunc(controller.CreateSecret)
	authMW := auth.AuthBearerMiddlewareInit[model.TokenInfo](testTokenSecret)
	authHandler := authMW(handler)

	body := model.CreateSecretDTO{
		Title:      "test-secret",
		SecretType: model.SecretTypeText,
		TextData:   "test-data",
	}
	bodyBytes, _ := json.Marshal(body)

	tokenInfo := model.TokenInfo{
		ID:    1,
		Login: "admin",
	}
	token, _ := auth.GenerateBearerToken[model.TokenInfo](tokenInfo, testTokenExp, testTokenSecret)

	req := httptest.NewRequest("POST", "/api/v1/secrets", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteSecret_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	mockStorage.EXPECT().
		DeleteSecret(context.Background(), "test-title", int64(1)).
		Return(nil)

	svc := service.New(
		zap.NewNop().Sugar(),
		mockStorage,
		testUserPasswordCost,
		testSecretPasswordKey,
		testTokenExp,
		testTokenSecret,
	)
	controller := New(zap.NewNop().Sugar(), svc)

	handler := http.HandlerFunc(controller.DeleteSecret)
	authMW := auth.AuthBearerMiddlewareInit[model.TokenInfo](testTokenSecret)
	authHandler := authMW(handler)

	body := model.DeleteSecretBody{Title: "test-title"}
	bodyBytes, _ := json.Marshal(body)

	tokenInfo := model.TokenInfo{
		ID:    1,
		Login: "admin",
	}
	token, _ := auth.GenerateBearerToken[model.TokenInfo](tokenInfo, testTokenExp, testTokenSecret)

	req := httptest.NewRequest("DELETE", "/api/v1/secret", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestDeleteSecret_InvalidBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	svc := service.New(
		zap.NewNop().Sugar(),
		mockStorage,
		testUserPasswordCost,
		testSecretPasswordKey,
		testTokenExp,
		testTokenSecret,
	)
	controller := New(zap.NewNop().Sugar(), svc)

	handler := http.HandlerFunc(controller.DeleteSecret)
	authMW := auth.AuthBearerMiddlewareInit[model.TokenInfo](testTokenSecret)
	authHandler := authMW(handler)

	req := httptest.NewRequest("DELETE", "/api/v1/secret", bytes.NewReader([]byte(`{invalid json}`)))
	req.Header.Set("Content-Type", "application/json")

	tokenInfo := model.TokenInfo{
		ID:    1,
		Login: "admin",
	}
	token, _ := auth.GenerateBearerToken[model.TokenInfo](tokenInfo, testTokenExp, testTokenSecret)
	req.Header.Set("Authorization", token)

	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteSecret_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	mockStorage.EXPECT().
		DeleteSecret(context.Background(), "not-found", int64(1)).
		Return(model.ErrSecretNotFound)

	svc := service.New(
		zap.NewNop().Sugar(),
		mockStorage,
		testUserPasswordCost,
		testSecretPasswordKey,
		testTokenExp,
		testTokenSecret,
	)
	controller := New(zap.NewNop().Sugar(), svc)

	handler := http.HandlerFunc(controller.DeleteSecret)
	authMW := auth.AuthBearerMiddlewareInit[model.TokenInfo](testTokenSecret)
	authHandler := authMW(handler)

	body := model.DeleteSecretBody{Title: "not-found"}
	bodyBytes, _ := json.Marshal(body)

	tokenInfo := model.TokenInfo{
		ID:    1,
		Login: "admin",
	}
	token, _ := auth.GenerateBearerToken[model.TokenInfo](tokenInfo, testTokenExp, testTokenSecret)

	req := httptest.NewRequest("DELETE", "/api/v1/secret", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteSecret_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockPG.NewMockStorage(ctrl)
	mockStorage.EXPECT().
		DeleteSecret(context.Background(), "test-title", int64(1)).
		Return(errors.New("internal storage error"))

	svc := service.New(
		zap.NewNop().Sugar(),
		mockStorage,
		testUserPasswordCost,
		testSecretPasswordKey,
		testTokenExp,
		testTokenSecret,
	)
	controller := New(zap.NewNop().Sugar(), svc)

	handler := http.HandlerFunc(controller.DeleteSecret)
	authMW := auth.AuthBearerMiddlewareInit[model.TokenInfo](testTokenSecret)
	authHandler := authMW(handler)

	body := model.DeleteSecretBody{Title: "test-title"}
	bodyBytes, _ := json.Marshal(body)

	tokenInfo := model.TokenInfo{
		ID:    1,
		Login: "admin",
	}
	token, _ := auth.GenerateBearerToken[model.TokenInfo](tokenInfo, testTokenExp, testTokenSecret)

	req := httptest.NewRequest("DELETE", "/api/v1/secret", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	w := httptest.NewRecorder()
	authHandler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
