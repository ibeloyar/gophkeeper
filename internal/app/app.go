package app

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/coreos/go-systemd/daemon"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/ibeloyar/gophkeeper/internal/config"
	"github.com/ibeloyar/gophkeeper/internal/model"
	"github.com/ibeloyar/gophkeeper/internal/repository/pgstorage"
	"github.com/ibeloyar/gophkeeper/internal/service"
	"github.com/ibeloyar/gophkeeper/pgk/auth"
	"github.com/ibeloyar/gophkeeper/pgk/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpcController "github.com/ibeloyar/gophkeeper/internal/controller/grpc"
	gophkeeperv1 "github.com/ibeloyar/gophkeeper/proto/gophkeeper/v1"
)

const (
	appName = "gophkeeper"

	shutdownTimeout = 5 * time.Second
)

func Run(cfg *config.Config) error {
	lg, err := logger.New()
	if err != nil {
		return err
	}
	defer lg.Sync()

	pgStorage, err := pgstorage.New(cfg.Database.DSN)
	if err != nil {
		return err
	}

	appService := service.New(lg, pgStorage, cfg.Security.UserPasswordCost, cfg.Security.SecretPasswordKey, cfg.Security.TokenLifetime, cfg.Security.TokenSecret)
	grpcControllerInstance := grpcController.New(lg, appService)
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_zap.UnaryServerInterceptor(lg.Desugar()),
			auth.AuthGRPCUnaryInterceptor[model.TokenInfo](cfg.Security.TokenSecret),
			grpc_recovery.UnaryServerInterceptor(
				grpc_recovery.WithRecoveryHandler(grpcControllerInstance.HandlePanic),
			),
		),
		grpc.ChainStreamInterceptor(
			grpc_zap.StreamServerInterceptor(lg.Desugar()),
			//auth.AuthGRPCStreamInterceptor[model.TokenInfo](cfg.Security.TokenSecret),
			grpc_recovery.StreamServerInterceptor(
				grpc_recovery.WithRecoveryHandler(grpcControllerInstance.HandlePanic),
			),
		),
	)
	gophkeeperv1.RegisterGophkeeperServer(server, grpcControllerInstance)
	reflection.Register(server)

	listener, err := net.Listen("tcp", cfg.GrpcServer.Addr)
	if err != nil {
		lg.Errorf("%s: %s", model.ErrListeningToLocalAddress, err)

		return err
	}

	go func() {
		if serveErr := server.Serve(listener); serveErr != nil {
			lg.Errorf("%s: %s", model.ErrStartingGrpcServer, serveErr)

			return
		}
	}()

	//if c.cfg.Swagger.Enabled {
	//	switch r.URL.Path {
	//	case "/swagger/doc.json":
	//		http.ServeFile(w, r, c.cfg.Swagger.JsonPath)
	//		return
	//	case "/swagger":
	//		http.Redirect(w, r, "/swagger/", http.StatusSeeOther)
	//		return
	//	}
	//
	//	if strings.Contains(r.URL.Path, "swagger") {
	//		httpSwagger.Handler(
	//			httpSwagger.URL("/swagger/doc.json"),
	//			httpSwagger.DocExpansion("none"),
	//			httpSwagger.Layout(httpSwagger.BaseLayout),
	//			httpSwagger.PersistAuthorization(true),
	//		).ServeHTTP(w, r)
	//		return
	//	}
	//}

	lg.Infof("%s starting server on %s", appName, cfg.GrpcServer.Addr)

	_, _ = daemon.SdNotify(false, daemon.SdNotifyReady)

	signalCtx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-signalCtx.Done()

	lg.Info("shutting down server...")

	stop := make(chan struct{})

	go func() {
		if shutdownStorageErr := pgStorage.Shutdown(); shutdownStorageErr != nil {
			lg.Errorf("shutdown error: %s", shutdownStorageErr)
		}

		server.GracefulStop()

		stop <- struct{}{}
	}()

	select {
	case <-stop:
		lg.Info("shutting down gracefully")
	case <-time.NewTicker(shutdownTimeout).C:
		lg.Error("shutting down timed out error")
	}

	return nil
}
