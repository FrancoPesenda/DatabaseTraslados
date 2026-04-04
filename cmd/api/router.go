package main

import (
	"net/http"

	userhandler "github.com/FrancoPesenda/eventra/cmd/handler/user/create_company"
)

func newHTTPMux(createHandler *userhandler.CreateHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /user/company", createHandler.Handle)

	return mux
}
