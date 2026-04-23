package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/ibeloyar/gophkeeper/internal/model"
	"github.com/ibeloyar/gophkeeper/pkg/auth"
	"go.uber.org/zap"
)

type Service interface {
	Register(ctx context.Context, input *model.RegisterDTO) (string, error)
	Login(ctx context.Context, input *model.LoginDTO) (string, error)

	CreateSecret(ctx context.Context, input *model.CreateSecretDTO) error
	GetSecret(ctx context.Context, title string, userID int64) (*model.Secret, error)
	GetSecrets(ctx context.Context, userID int64) ([]*model.Secret, error)
	DeleteSecret(ctx context.Context, title string, userID int64) error
}

type Controller struct {
	service Service
	lg      *zap.SugaredLogger
}

// New creates Controller with logger and service dependencies injected.
func New(lg *zap.SugaredLogger, s Service) *Controller {
	return &Controller{
		lg:      lg,
		service: s,
	}
}

// Register handles HTTP POST /register. Parses RegisterDTO from JSON body.
// Maps business errors to HTTP 400/409/500. Returns Bearer token in Authorization header.
func (c *Controller) Register(w http.ResponseWriter, r *http.Request) {
	body, err := readBody[model.RegisterDTO](r)
	if err != nil {
		c.lg.Errorf("failed to parse request body: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	bearerToken, err := c.service.Register(context.Background(), &body)
	if err != nil {
		if errors.Is(err, model.ErrInvalidLoginOrPassword) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, model.ErrUserAlreadyExist) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		if errors.Is(err, model.ErrServerInternal) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Authorization", bearerToken)
	w.WriteHeader(http.StatusOK)
}

// Login handles HTTP POST /login. Parses LoginDTO from JSON body.
// Maps business errors to HTTP 400/500. Returns Bearer token in Authorization header.
func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	body, err := readBody[model.LoginDTO](r)
	if err != nil {
		c.lg.Errorf("failed to parse request body: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	bearerToken, err := c.service.Login(context.Background(), &body)
	if err != nil {
		if errors.Is(err, model.ErrInvalidLoginOrPassword) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, model.ErrServerInternal) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Authorization", bearerToken)
	w.WriteHeader(http.StatusOK)
}

// GetSecret handles HTTP GET/POST /secret. Requires JWT auth. Parses title from body.
// Returns 404 for non-existent secrets. JSON response with secret details.
func (c *Controller) GetSecret(w http.ResponseWriter, r *http.Request) {
	tokenInfo := auth.GetTokenInfo[model.TokenInfo](r)

	body, err := readBody[model.GetSecretBody](r)
	if err != nil {
		c.lg.Errorf("failed to parse request body: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	secret, err := c.service.GetSecret(context.Background(), body.Title, tokenInfo.ID)
	if err != nil {
		if errors.Is(err, model.ErrSecretNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		c.lg.Errorf("GetSecret error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writeJSON(w, c.lg, secret, http.StatusOK)
}

// GetSecrets handles HTTP GET /secrets. Requires JWT auth. Returns JSON array of user's secrets.
func (c *Controller) GetSecrets(w http.ResponseWriter, r *http.Request) {
	tokenInfo := auth.GetTokenInfo[model.TokenInfo](r)

	secrets, err := c.service.GetSecrets(context.Background(), tokenInfo.ID)
	if err != nil {
		c.lg.Errorf("GetSecret error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writeJSON(w, c.lg, secrets, http.StatusOK)
}

// CreateSecret handles HTTP POST /secret. Requires JWT auth. Parses CreateSecretDTO from body.
// Maps validation errors to HTTP 400. Sets UserID from JWT token.
func (c *Controller) CreateSecret(w http.ResponseWriter, r *http.Request) {
	body, err := readBody[model.CreateSecretDTO](r)
	if err != nil {
		c.lg.Errorf("failed to parse request body: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	tokenInfo := auth.GetTokenInfo[model.TokenInfo](r)
	body.UserID = tokenInfo.ID

	if err := c.service.CreateSecret(context.Background(), &body); err != nil {
		if errors.Is(err, model.ErrTitleIsRequired) ||
			errors.Is(err, model.ErrSecretLoginRequired) ||
			errors.Is(err, model.ErrSecretLoginMaxLen) ||
			errors.Is(err, model.ErrSecretPasswordRequired) ||
			errors.Is(err, model.ErrSecretPasswordMaxLen) ||
			errors.Is(err, model.ErrSecretTextDataRequired) ||
			errors.Is(err, model.ErrSecretBinaryRequired) ||
			errors.Is(err, model.ErrSecretCardNumberRequired) ||
			errors.Is(err, model.ErrSecretCardNumberMaxLen) ||
			errors.Is(err, model.ErrSecretCardExpRequired) ||
			errors.Is(err, model.ErrSecretCardExpInvalid) ||
			errors.Is(err, model.ErrSecretCardHolderRequired) ||
			errors.Is(err, model.ErrSecretCardHolderMaxLen) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		c.lg.Errorf("CreateSecret error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteSecret handles HTTP DELETE /secret. Requires JWT auth. Parses title from body.
// Returns 404 for non-existent secrets belonging to user.
func (c *Controller) DeleteSecret(w http.ResponseWriter, r *http.Request) {
	tokenInfo := auth.GetTokenInfo[model.TokenInfo](r)

	body, err := readBody[model.DeleteSecretBody](r)
	if err != nil {
		c.lg.Errorf("failed to parse request body: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := c.service.DeleteSecret(context.Background(), body.Title, tokenInfo.ID); err != nil {
		if errors.Is(err, model.ErrSecretNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		c.lg.Errorf("DeleteSecret error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
