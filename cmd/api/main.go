package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"database/internal/config"
	"database/internal/infrastructure/persistence/mysql"
	"database/internal/interfaces/http/server"
)

func main() {
	cfg := config.FromEnv()

	logger := log.New(os.Stdout, "", log.LstdFlags|log.LUTC)

	db, err := mysql.Open(cfg.MySQL.DSN)
	if err != nil {
		logger.Fatalf("mysql open: %v", err)
	}
	defer db.Close()

	srv := server.New(server.Config{
		Addr:   cfg.HTTP.Addr,
		Logger: logger,
	})

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Printf("http listening on %s", cfg.HTTP.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("http serve: %v", err)
		}
	}()

	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	logger.Println("shutdown complete")
}

