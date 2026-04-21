package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/ibeloyar/gophkeeper/internal/model"
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

func New(lg *zap.SugaredLogger, s Service) *Controller {
	return &Controller{
		lg:      lg,
		service: s,
	}
}

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
