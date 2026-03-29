package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	Logger *log.Logger
	HTTP   *http.Server
	Close  func(context.Context) error
}

func (a *App) Run(ctx context.Context) error {
	if a.Logger == nil {
		a.Logger = log.New(os.Stdout, "", log.LstdFlags|log.LUTC)
	}
	if a.Close == nil {
		a.Close = func(context.Context) error { return nil }
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	errCh := make(chan error, 1)
	go func() {
		a.Logger.Printf("http listening on %s", a.HTTP.Addr)
		if err := a.HTTP.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
	case <-stop:
	case err := <-errCh:
		return err
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = a.HTTP.Shutdown(shutdownCtx)
	_ = a.Close(shutdownCtx)
	a.Logger.Println("shutdown complete")
	return nil
}

