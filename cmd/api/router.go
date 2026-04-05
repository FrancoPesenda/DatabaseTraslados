package main

import (
	"net/http"

	userhandler "github.com/FrancoPesenda/eventra/cmd/handler/user/createcompany"
)

func newHTTPMux(createHandler *userhandler.CreateHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /user/company", createHandler.Handle)

	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("pong"))
	})

	return mux
}
