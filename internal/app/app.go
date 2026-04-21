package app

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/coreos/go-systemd/daemon"
	"github.com/ibeloyar/gophkeeper/internal/config"
	"github.com/ibeloyar/gophkeeper/internal/model"
	"github.com/ibeloyar/gophkeeper/internal/repository/pgstorage"
	"github.com/ibeloyar/gophkeeper/internal/service"
	"github.com/ibeloyar/gophkeeper/pgk/logger"
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

	appService := service.New(
		lg, pgStorage,
		cfg.Security.UserPasswordCost,
		cfg.Security.SecretPasswordKey,
		cfg.Security.TokenLifetime,
		cfg.Security.TokenSecret,
	)

	grpcServer, err := buildGRPCServer(lg, appService, cfg.Security.TokenSecret)

	grpcListener, err := net.Listen("tcp", cfg.GrpcServer.Addr)
	if err != nil {
		lg.Errorf("%s: %s", model.ErrListeningToLocalAddress, err)

		return err
	}

	go func() {
		if serveErr := grpcServer.Serve(grpcListener); serveErr != nil {
			lg.Errorf("%s: %s", model.ErrStartingGrpcServer, serveErr)

			return
		}
	}()

	lg.Infof("%s starting GRPC server on %s", appName, cfg.GrpcServer.Addr)

	var httpServer *http.Server
	if cfg.HttpServer.Addr != "" {
		httpServer, err = buildHTTPServer(lg, appService, cfg.HttpServer.Addr, cfg.Security.TokenSecret)
		if err != nil {
			return err
		}

		httpListener, err := net.Listen("tcp", cfg.HttpServer.Addr)
		if err != nil {
			return err
		}

		go func() {
			lg.Infof("%s starting HTTP server on %s", appName, cfg.HttpServer.Addr)

			if serveErr := httpServer.Serve(httpListener); serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
				lg.Errorf("http server error: %s", serveErr)
			}
		}()
	}

	_, _ = daemon.SdNotify(false, daemon.SdNotifyReady)

	signalCtx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-signalCtx.Done()

	lg.Info("shutting down server...")

	stop := make(chan struct{})

	go func() {
		if shutdownStorageErr := pgStorage.Shutdown(); shutdownStorageErr != nil {
			lg.Errorf("shutdown error: %s", shutdownStorageErr)
		}

		grpcServer.GracefulStop()

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
