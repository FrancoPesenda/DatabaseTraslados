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
)

// Build wires the application manually (alternativa al NewFxApp basado en FX).
// main.go usa NewFxApp; esta función se mantiene para referencia/testing.
func Build(ctx context.Context) (*App, error) {
	_ = ctx

	logger := log.New(os.Stdout, "", log.LstdFlags|log.LUTC)

	cfg, err := LoadInfraConfig(LoadConfigOptions{})
	if err != nil {
		return nil, err
	}

	db, err := openDB(cfg)
	if err != nil {
		return nil, err
	}

	var repo admin.Repository = eventradatabase.New(db)

	healthUC := usecase.NewHealthUsecase()
	createAdminUC := usecase.NewCreateAdminUsecase(repo)

	healthH := handlers.NewHealthHandler(healthUC)
	usersH := handlers.NewUsersHandler(createAdminUC)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthH.Get)
	mux.HandleFunc("POST /users", usersH.Create)

	srv := &http.Server{
		Addr:              cfg.HTTP.Addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &App{
		Logger: logger,
		HTTP:   srv,
		Close: func(context.Context) error {
			return db.Close()
		},
	}, nil
}

func openDB(cfg InfraConfig) (*sql.DB, error) {
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
