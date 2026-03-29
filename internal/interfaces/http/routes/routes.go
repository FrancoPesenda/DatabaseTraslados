package routes

import (
	"log"
	"net/http"

	"database/internal/interfaces/http/handlers"
)

type Config struct {
	Logger *log.Logger
}

func New(cfg Config) http.Handler {
	if cfg.Logger == nil {
		cfg.Logger = log.Default()
	}

	mux := http.NewServeMux()

	health := handlers.NewHealthHandler()
	mux.HandleFunc("GET /health", health.Get)

	return mux
}

