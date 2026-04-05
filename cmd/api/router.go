package main

import (
	"net/http"

	createcompanyhandler "github.com/FrancoPesenda/eventra/cmd/handler/user/createcompany"
	loginhandler "github.com/FrancoPesenda/eventra/cmd/handler/user/login"
)

func newHTTPMux(createHandler *createcompanyhandler.CreateHandler, loginHandler *loginhandler.Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /user/company", createHandler.Handle)
	mux.HandleFunc("POST /user/login", loginHandler.Login)

	return mux
}
