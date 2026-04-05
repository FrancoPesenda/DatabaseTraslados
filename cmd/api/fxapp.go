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

	createcompanyhandler "github.com/FrancoPesenda/eventra/cmd/handler/user/createcompany"
	loginhandler "github.com/FrancoPesenda/eventra/cmd/handler/user/login"
	"github.com/FrancoPesenda/eventra/internal/repository/eventradatabase"
	"github.com/FrancoPesenda/eventra/internal/usecase/user/createcompany"
	"github.com/FrancoPesenda/eventra/internal/usecase/user/login"
	"github.com/FrancoPesenda/eventra/internal/utils/config"
)

func NewFxApp() *fx.App {
	return fx.New(
		fx.Provide(
			newLogger,
			newInfraConfig,
			newDB,
			eventradatabase.NewRepository,
			func(r *eventradatabase.Repository) createcompany.UserRepository { return r },
			func(r *eventradatabase.Repository) login.UserRepository { return r },
			fx.Annotate(
				createcompany.NewUseCase,
				fx.As(new(createcompanyhandler.UseCase)),
			),
			createcompanyhandler.NewCreateHandler,
			fx.Annotate(
				login.NewUseCase,
				fx.As(new(loginhandler.UseCase)),
			),
			loginhandler.NewHandler,
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

func newDB(cfg config.InfraConfig) (*sql.DB, error) {
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

	idleTime, err := time.ParseDuration(cfg.MySQL.ConnMaxIdleTime)
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("invalid MYSQL_CONN_MAX_IDLE_TIME: %w", err)
	}
	db.SetConnMaxIdleTime(idleTime)

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
