package grpc

import (
	"context"
	"errors"
	"runtime/debug"

	"github.com/ibeloyar/gophkeeper/internal/model"
	"github.com/ibeloyar/gophkeeper/pgk/auth"
	gophkeeperv1 "github.com/ibeloyar/gophkeeper/proto/gophkeeper/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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
	gophkeeperv1.UnimplementedGophkeeperServer

	logger  *zap.SugaredLogger
	service Service
}

func New(lg *zap.SugaredLogger, service Service) *Controller {
	return &Controller{
		logger:  lg,
		service: service,
	}
}

func (c *Controller) HandlePanic(p any) error {
	if p != nil {
		c.logger.Errorf("%v\n%s", p, string(debug.Stack()))
	}
	return model.ErrServerInternal
}

func (c *Controller) Register(
	ctx context.Context,
	req *gophkeeperv1.RegisterRequest,
) (*gophkeeperv1.RegisterResponse, error) {
	token, err := c.service.Register(ctx, &model.RegisterDTO{
		Login:    req.Login,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, model.ErrInvalidLoginOrPassword) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, model.ErrUserAlreadyExist) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		if errors.Is(err, model.ErrServerInternal) {
			return nil, status.Error(codes.Internal, err.Error())
		}
		return nil, err
	}

	if err := grpc.SendHeader(ctx, metadata.Pairs("token", token)); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &gophkeeperv1.RegisterResponse{}, nil
}

func (c *Controller) Login(
	ctx context.Context,
	req *gophkeeperv1.LoginRequest,
) (*gophkeeperv1.LoginResponse, error) {
	token, err := c.service.Login(ctx, &model.LoginDTO{
		Login:    req.Login,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, model.ErrInvalidLoginOrPassword) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, model.ErrServerInternal) {
			return nil, status.Error(codes.Internal, err.Error())
		}
		return nil, err
	}

	if err := grpc.SendHeader(ctx, metadata.Pairs("token", token)); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &gophkeeperv1.LoginResponse{}, nil
}

func (c *Controller) CreateSecret(
	ctx context.Context,
	req *gophkeeperv1.CreateSecretRequest,
) (*gophkeeperv1.CreateSecretResponse, error) {
	dto, err := convertCreateSecretToDTO(ctx, req)
	if err != nil {
		return nil, err
	}

	if err = c.service.CreateSecret(ctx, dto); err != nil {
		return nil, err
	}

	return &gophkeeperv1.CreateSecretResponse{}, nil
}

func (c *Controller) GetSecret(
	ctx context.Context,
	req *gophkeeperv1.GetSecretRequest,
) (*gophkeeperv1.GetSecretResponse, error) {
	tokenInfo := auth.GetTokenInfoFromContext[model.TokenInfo](ctx)

	secret, err := c.service.GetSecret(ctx, req.Title, tokenInfo.ID)
	if secret == nil && err == nil {
		return nil, status.Error(codes.NotFound, model.ErrSecretNotFound.Error())
	}
	if err != nil {
		return nil, err
	}

	return convertGetSecretToProto(secret), nil
}

func (c *Controller) GetSecrets(
	ctx context.Context,
	req *gophkeeperv1.GetSecretsRequest,
) (*gophkeeperv1.GetSecretsResponse, error) {
	tokenInfo := auth.GetTokenInfoFromContext[model.TokenInfo](ctx)

	secrets, err := c.service.GetSecrets(ctx, tokenInfo.ID)
	if err != nil {
		return nil, err
	}

	return &gophkeeperv1.GetSecretsResponse{
		Secrets: convertGetSecretsToProto(secrets),
	}, nil
}

func (c *Controller) DeleteSecret(
	ctx context.Context,
	req *gophkeeperv1.DeleteSecretRequest,
) (*gophkeeperv1.DeleteSecretResponse, error) {
	tokenInfo := auth.GetTokenInfoFromContext[model.TokenInfo](ctx)

	if err := c.service.DeleteSecret(ctx, req.Title, tokenInfo.ID); err != nil {
		if errors.Is(err, model.ErrSecretNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}

	return &gophkeeperv1.DeleteSecretResponse{}, nil
}
