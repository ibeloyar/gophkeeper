package auth

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type tokenDataContextKeyType string

const tokenDataContextKey = tokenDataContextKeyType("token")

type Claims[T any] struct {
	jwt.RegisteredClaims
	TokenInfo T
}

func GenerateBearerToken[T any](input T, exp time.Duration, secret string) (token string, err error) {
	tokenData := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims[T]{
		TokenInfo: input,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	token, err = tokenData.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Bearer %s", token), nil
}

func VerifyJWTBearerToken[T any](tokenString, secret string) (*T, error) {
	claims := &Claims[T]{}

	tokenParts := strings.Split(tokenString, " ")
	if len(tokenParts) != 2 {
		return nil, jwt.ErrSignatureInvalid
	}
	if tokenParts[0] != "Bearer" {
		return nil, jwt.ErrInvalidType
	}

	token, err := jwt.ParseWithClaims(tokenParts[1], claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidKeyType
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return &claims.TokenInfo, nil
}

func AuthBearerMiddlewareInit[T any](secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenInfo, err := VerifyJWTBearerToken[T](r.Header.Get("Authorization"), secret)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), tokenDataContextKey, tokenInfo)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetTokenInfo[T any](r *http.Request) *T {
	ctx := r.Context()

	tokenInfo, ok := ctx.Value(tokenDataContextKey).(*T)
	if !ok {
		return nil
	}

	return tokenInfo
}

func GetTokenInfoFromContext[T any](ctx context.Context) *T {
	tokenInfo, ok := ctx.Value(tokenDataContextKey).(*T)
	if !ok {
		return nil
	}
	return tokenInfo
}

func AuthGRPCUnaryInterceptor[T any](secret string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		withoutAuthMethods := []string{
			"/gophkeeper.v1.Gophkeeper/Register",
			"/gophkeeper.v1.Gophkeeper/Login",
		}

		if slices.Contains(withoutAuthMethods, info.FullMethod) {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		authHeaders := md.Get("Authorization")
		if len(authHeaders) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		}

		tokenInfo, err := VerifyJWTBearerToken[T](authHeaders[0], secret)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		ctx = context.WithValue(ctx, tokenDataContextKey, tokenInfo)
		return handler(ctx, req)
	}
}
