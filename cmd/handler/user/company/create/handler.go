package user

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	domain "eventra/internal/domain"
	userCreate "eventra/internal/usecase/user/company/create"
)

type CreateHandler struct {
	usecase *userCreate.UseCase
	logger  *log.Logger
}

func NewCreateHandler(usecase *userCreate.UseCase, logger *log.Logger) *CreateHandler {
	return &CreateHandler{
		usecase: usecase,
		logger:  logger,
	}
}

func (h *CreateHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("[Layer:Handler][error_message:%s][request_body:%+v]", err.Error(), req)
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	h.logger.Printf("Incoming Request: %+v", req)

	out, err := h.usecase.CreateCompanyUser(r.Context(), req.toDomain())
	if err != nil {
		h.logger.Printf("[Layer:Handler][error_message:%s][request_body:%+v]", err.Error(), req)
		processError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(newCreateResponse(out))
}

func processError(w http.ResponseWriter, err error) {
	switch err {
	case domain.ErrNameRequired,
		domain.ErrLastNameRequired,
		domain.ErrEmailRequired,
		domain.ErrPasswordRequired:
		writeError(w, http.StatusBadRequest, fmt.Sprintf("Bad Request error: %s", err.Error()))
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errorResponse{Error: msg})
}
