package api

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	mysqlrepo "database/internal/repository/mysql"
	handlers "database/internal/handler/http"
	"database/internal/usecase"

	"go.uber.org/fx"
)

// NewFxApp builds the application using Uber FX (DI + lifecycle).
func NewFxApp() *fx.App {
	return fx.New(
		fx.Provide(
			newLogger,
			newInfraConfig,
			newDB,
			newHTTPHandler,
			newHTTPServer,
		),
		fx.Invoke(registerLifecycle),
	)
}

func newLogger() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags|log.LUTC)
}

func newInfraConfig() (InfraConfig, error) {
	return LoadInfraConfig(LoadConfigOptions{})
}

func newDB(cfg InfraConfig) (*sql.DB, error) {
	return mysqlrepo.Open(cfg.MySQL.DSN)
}

func newHTTPHandler(healthUC *usecase.HealthUsecase) http.Handler {
	healthH := handlers.NewHealthHandler(healthUC)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthH.Get)
	return mux
}

func newHTTPServer(cfg InfraConfig, h http.Handler) *http.Server {
	return &http.Server{
		Addr:              cfg.HTTP.Addr,
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func registerLifecycle(lc fx.Lifecycle, logger *log.Logger, srv *http.Server, db *sql.DB) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Printf("http listening on %s", srv.Addr)
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Printf("http serve error: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			_ = srv.Shutdown(ctx)
			return db.Close()
		},
	})
}

