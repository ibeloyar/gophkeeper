package app

import (
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcController "github.com/ibeloyar/gophkeeper/internal/controller/grpc"
	"github.com/ibeloyar/gophkeeper/internal/model"
	"github.com/ibeloyar/gophkeeper/internal/service"
	"github.com/ibeloyar/gophkeeper/pgk/auth"
	gophkeeperv1 "github.com/ibeloyar/gophkeeper/proto/gophkeeper/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func buildGRPCServer(lg *zap.SugaredLogger, appService *service.Service, tokenSecret string) (*grpc.Server, error) {
	grpcControllerInstance := grpcController.New(lg, appService)
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_zap.UnaryServerInterceptor(lg.Desugar()),
			auth.AuthGRPCUnaryInterceptor[model.TokenInfo](tokenSecret),
			grpc_recovery.UnaryServerInterceptor(
				grpc_recovery.WithRecoveryHandler(grpcControllerInstance.HandlePanic),
			),
		),
	)
	gophkeeperv1.RegisterGophkeeperServer(server, grpcControllerInstance)
	reflection.Register(server)

	return server, nil
}
