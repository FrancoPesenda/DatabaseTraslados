package http

import (
	"encoding/json"
	"errors"
	nethttp "net/http"

	"database/internal/usecase"
)

type UsersHandler struct {
	createUC *usecase.CreateUserUsecase
}

func NewUsersHandler(createUC *usecase.CreateUserUsecase) *UsersHandler {
	return &UsersHandler{createUC: createUC}
}

type createUserRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (h *UsersHandler) Create(w nethttp.ResponseWriter, r *nethttp.Request) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(nethttp.StatusBadRequest)
		return
	}

	out, err := h.createUC.Execute(r.Context(), usecase.CreateUserInput{
		Email: req.Email,
		Name:  req.Name,
	})
	if err != nil {
		if errors.Is(err, usecase.ErrValidation) {
			w.WriteHeader(nethttp.StatusBadRequest)
			return
		}
		w.WriteHeader(nethttp.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(nethttp.StatusCreated)
	_ = json.NewEncoder(w).Encode(out)
}

