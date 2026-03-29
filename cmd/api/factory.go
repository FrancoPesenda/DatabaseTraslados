package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	handlers "database/internal/handler/http"
	mysqlrepo "database/internal/repository/mysql"
	"database/internal/usecase"
)

func Build(ctx context.Context) (*App, error) {
	_ = ctx

	logger := log.New(os.Stdout, "", log.LstdFlags|log.LUTC)

	cfg, err := LoadInfraConfig(LoadConfigOptions{})
	if err != nil {
		return nil, err
	}

	db, err := mysqlrepo.Open(cfg.MySQL.DSN)
	if err != nil {
		return nil, err
	}

	healthUC := usecase.NewHealthUsecase()
	healthH := handlers.NewHealthHandler(healthUC)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthH.Get)

	srv := &http.Server{
		Addr:              cfg.HTTP.Addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	app := &App{
		Logger: logger,
		HTTP:   srv,
		Close: func(context.Context) error {
			return db.Close()
		},
	}

	return app, nil
}

