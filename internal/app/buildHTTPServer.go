package app

import (
	"net/http"
	"time"

	"github.com/ibeloyar/gophkeeper/internal/config"
	"github.com/ibeloyar/gophkeeper/internal/model"
	"github.com/ibeloyar/gophkeeper/internal/service"
	"github.com/ibeloyar/gophkeeper/pgk/auth"
	"go.uber.org/zap"

	httpController "github.com/ibeloyar/gophkeeper/internal/controller/http"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func buildHTTPServer(lg *zap.SugaredLogger, appService *service.Service, cfg *config.Config) (*http.Server, error) {
	controller := httpController.New(lg, appService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/register", controller.Register)
	mux.HandleFunc("POST /api/v1/login", controller.Login)

	if cfg.Swagger.Enabled {
		mux.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, cfg.Swagger.JsonPath)
		})

		mux.HandleFunc("/swagger", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/swagger/", http.StatusSeeOther)
		})

		mux.Handle("/swagger/", http.StripPrefix(
			"/swagger",
			httpSwagger.Handler(
				httpSwagger.URL("/swagger/doc.json"),
				httpSwagger.PersistAuthorization(true),
			),
		))
	}

	protected := http.NewServeMux()

	protected.HandleFunc("POST /api/v1/get-secret", controller.GetSecret)
	protected.HandleFunc("POST /api/v1/get-secrets", controller.GetSecrets)
	protected.HandleFunc("POST /api/v1/secrets-create", controller.CreateSecret)
	protected.HandleFunc("POST /api/v1/delete-secret", controller.DeleteSecret)

	// Middleware: auth + protected mux
	authHandler := auth.AuthBearerMiddlewareInit[model.TokenInfo](cfg.Security.TokenSecret)(protected)

	mux.HandleFunc("/api/v1/", func(w http.ResponseWriter, r *http.Request) {
		authHandler.ServeHTTP(w, r)
	})

	server := &http.Server{
		Addr:         cfg.HttpServer.Addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	return server, nil
}
