package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/fx"

	userhandler "github.com/FrancoPesenda/eventra/cmd/handler/user/create_company"
	"github.com/FrancoPesenda/eventra/internal/config"
	"github.com/FrancoPesenda/eventra/internal/repository/eventradatabase"
	userCreate "github.com/FrancoPesenda/eventra/internal/usecase/user/company/create"
)

func NewFxApp() *fx.App {
	return fx.New(
		fx.Provide(
			newLogger,
			newInfraConfig,
			newDB,
			fx.Annotate(
				eventradatabase.NewRepository,
				fx.As(new(userCreate.UserRepository)),
			),
			userCreate.NewUseCase,
			userhandler.NewCreateHandler,
			newHTTPMux,
			newHTTPServer,
		),
		fx.Invoke(registerLifecycle),
	)
}

func newLogger() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags|log.LUTC)
}

func newInfraConfig() (config.InfraConfig, error) {
	return config.LoadInfraConfig(config.LoadConfigOptions{})
}

func newDB(cfg config.InfraConfig) (eventradatabase.DB, error) {
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

func newHTTPServer(cfg config.InfraConfig, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              cfg.HTTP.Addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func registerLifecycle(lc fx.Lifecycle, logger *log.Logger, srv *http.Server, db eventradatabase.DB) {
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
