package grpc

import (
	"context"
	"errors"
	"runtime/debug"

	"github.com/ibeloyar/gophkeeper/internal/model"
	"github.com/ibeloyar/gophkeeper/pkg/auth"
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

// Controller implements gophkeeperv1.GophkeeperServer with business logic delegation.
// Embeds UnimplementedGophkeeperServer and handles panic recovery, error mapping.
type Controller struct {
	gophkeeperv1.UnimplementedGophkeeperServer

	logger  *zap.SugaredLogger
	service Service
}

// New creates Controller instance with logger and service dependencies.
func New(lg *zap.SugaredLogger, service Service) *Controller {
	return &Controller{
		logger:  lg,
		service: service,
	}
}

// HandlePanic recovers from panics, logs stack trace, and returns standard internal error.
// Called by gRPC unary interceptor for global panic protection.
func (c *Controller) HandlePanic(p any) error {
	if p != nil {
		c.logger.Errorf("%v\n%s", p, string(debug.Stack()))
	}
	return model.ErrServerInternal
}

// Register handles gRPC RegisterRequest. Creates user account via service layer.
// Maps business errors to gRPC codes: InvalidArgument, AlreadyExists, Internal.
// Sends JWT token via gRPC headers for client retrieval.
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

// Login handles gRPC LoginRequest. Authenticates user via service layer.
// Maps business errors to gRPC codes. Sends JWT token via headers.
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

// CreateSecret handles gRPC CreateSecretRequest. Converts proto to DTO and delegates to service.
// Maps validation errors to gRPC InvalidArgument, others to Internal.
func (c *Controller) CreateSecret(
	ctx context.Context,
	req *gophkeeperv1.CreateSecretRequest,
) (*gophkeeperv1.CreateSecretResponse, error) {
	dto, err := convertCreateSecretToDTO(ctx, req)
	if err != nil {
		return nil, err
	}

	if err = c.service.CreateSecret(ctx, dto); err != nil {
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
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, err
	}

	return &gophkeeperv1.CreateSecretResponse{}, nil
}

// GetSecret handles gRPC GetSecretRequest. Extracts userID from JWT context.
// Returns NotFound if secret doesn't exist for user, delegates other errors.
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

// GetSecrets handles gRPC GetSecretsRequest. Extracts userID from JWT and lists user's secrets.
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

// DeleteSecret handles gRPC DeleteSecretRequest. Extracts userID from JWT context.
// Maps secret not found to gRPC NotFound status.
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
