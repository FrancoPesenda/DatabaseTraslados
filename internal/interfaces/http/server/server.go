package server

import (
	"log"
	"net/http"
	"time"

	"database/internal/interfaces/http/routes"
)

type Config struct {
	Addr   string
	Logger *log.Logger
}

func New(cfg Config) *http.Server {
	if cfg.Logger == nil {
		cfg.Logger = log.Default()
	}

	h := routes.New(routes.Config{
		Logger: cfg.Logger,
	})

	return &http.Server{
		Addr:              cfg.Addr,
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
	}
}

