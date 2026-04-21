package app

import (
	"net/http"
	"time"

	"github.com/ibeloyar/gophkeeper/internal/model"
	"github.com/ibeloyar/gophkeeper/internal/service"
	"github.com/ibeloyar/gophkeeper/pgk/auth"
	"go.uber.org/zap"

	httpController "github.com/ibeloyar/gophkeeper/internal/controller/http"
)

func buildHTTPServer(lg *zap.SugaredLogger, appService *service.Service, addr, tokenSecret string) (*http.Server, error) {
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

	controller := httpController.New(lg, appService)

	// Общий mux для всего HTTP
	mux := http.NewServeMux()

	// Публичные эндпоинты — без auth
	mux.HandleFunc("POST /api/v1/register", controller.Register)
	mux.HandleFunc("POST /api/v1/login", controller.Login)

	// Закрытые эндпоинты (будут за middleware)
	protected := http.NewServeMux()

	// protected.HandleFunc("/api/v1/secrets-create", controller.CreateSecret)
	// protected.HandleFunc("/api/v1/secrets", controller.GetSecrets)
	// protected.HandleFunc("/api/v1/secret", controller.GetSecret)
	// protected.HandleFunc("/api/v1/secret", controller.DeleteSecret)

	// Middleware: auth + protected mux
	authHandler := auth.AuthBearerMiddlewareInit[model.TokenInfo](tokenSecret)(protected)

	mux.HandleFunc("/api/v1/", func(w http.ResponseWriter, r *http.Request) {
		authHandler.ServeHTTP(w, r)
	})

	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	return server, nil
}
