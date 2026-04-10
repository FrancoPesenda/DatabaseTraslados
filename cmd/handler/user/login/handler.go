package login

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	domain "github.com/FrancoPesenda/eventra/internal/domain"
)

type UseCase interface {
	Login(ctx context.Context, user domain.User) (domain.User, error)
}

type Handler struct {
	usecase UseCase
	logger  *log.Logger
}

func NewHandler(usecase UseCase, logger *log.Logger) *Handler {
	return &Handler{
		usecase: usecase,
		logger:  logger,
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("[Layer:Handler][error_message:%s][request_body:%+v]", err.Error(), req)
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	h.logger.Printf("Incoming Request: %+v", req)

	out, err := h.usecase.Login(r.Context(), req.toDomain())
	if err != nil {
		h.logger.Printf("[Layer:Handler][error_message:%s][request_body:%+v]", err.Error(), req)
		processError(w, err)
		return
	}

	response, err := newLoginResponse(out)
	if err != nil {
		h.logger.Printf("[Layer:Handler][error_message:%s][user:%+v]", err.Error(), out)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func processError(w http.ResponseWriter, err error) {
	switch err {
	case domain.ErrUserNameRequired, domain.ErrPasswordRequired:
		writeError(w, http.StatusBadRequest, err.Error())
	case domain.ErrInvalidCredentials:
		writeError(w, http.StatusUnauthorized, err.Error())
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errorResponse{Error: msg})
}
