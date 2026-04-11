package main

import (
	"net/http"

	createeventhandler "github.com/FrancoPesenda/eventra/cmd/handler/event/create"
	createcompanyhandler "github.com/FrancoPesenda/eventra/cmd/handler/user/createcompany"
	loginhandler "github.com/FrancoPesenda/eventra/cmd/handler/user/login"
)

func newHTTPMux(createHandler *createcompanyhandler.CreateHandler, loginHandler *loginhandler.Handler, createEventHandler *createeventhandler.CreateHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /user/company", createHandler.Handle)
	mux.HandleFunc("POST /user/login", loginHandler.Login)
	mux.HandleFunc("POST /event", createEventHandler.Handle)
	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("pong"))
	})

	return mux
}
