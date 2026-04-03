package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"database/internal/domain/admin"
	handlers "database/internal/handler/http"
	"database/internal/repository/eventradatabase"
	"database/internal/usecase"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/fx"
)

// NewFxApp builds the application using Uber FX (DI + lifecycle).
func NewFxApp() *fx.App {
	return fx.New(
		fx.Provide(
			newLogger,
			newInfraConfig,
			newDB,
			newEventraRepo,
			usecase.NewHealthUsecase,
			usecase.NewCreateAdminUsecase,
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
	db, err := sql.Open("mysql", cfg.MySQL.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MySQL.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MySQL.MaxIdleConns)
	d, err := time.ParseDuration(cfg.MySQL.ConnMaxLifetime)
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("invalid MYSQL_CONN_MAX_LIFETIME: %w", err)
	}
	db.SetConnMaxLifetime(d)

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

func newEventraRepo(db *sql.DB) admin.Repository {
	return eventradatabase.New(db)
}

func newHTTPHandler(healthUC *usecase.HealthUsecase, createAdminUC *usecase.CreateAdminUsecase, logger *log.Logger) http.Handler {
	healthH := handlers.NewHealthHandler(healthUC)
	usersH := handlers.NewUsersHandler(createAdminUC, logger)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthH.Get)
	mux.HandleFunc("POST /users", usersH.Create)
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
