package create

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	domain "github.com/FrancoPesenda/eventra/internal/domain"
	"github.com/FrancoPesenda/eventra/internal/utils/token"
)

type UseCase interface {
	CreateEvent(ctx context.Context, admin domain.User, event domain.Event) (domain.Event, error)
}

type CreateHandler struct {
	usecase UseCase
	logger  *log.Logger
}

func NewCreateHandler(usecase UseCase, logger *log.Logger) *CreateHandler {
	return &CreateHandler{
		usecase: usecase,
		logger:  logger,
	}
}

func (h *CreateHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// Extraer y validar token JWT
	authHeader := r.Header.Get("Authorization")
	tokenString, err := token.ExtractTokenFromHeader(authHeader)
	if err != nil {
		h.logger.Printf("[Layer:Handler][error_message:%s]", err.Error())
		writeError(w, http.StatusUnauthorized, "missing or invalid authorization token")
		return
	}

	claims, err := token.ValidateToken(tokenString)
	if err != nil {
		h.logger.Printf("[Layer:Handler][error_message:%s]", err.Error())
		writeError(w, http.StatusUnauthorized, "invalid or expired token")
		return
	}

	admin := token.ClaimsToUser(claims)

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("[Layer:Handler][error_message:%s][request_body:%+v]", err.Error(), req)
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	h.logger.Printf("Incoming Request: %+v [User:%s]", req, admin.UserName)

	out, err := h.usecase.CreateEvent(r.Context(), admin, req.eventDomain())
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
	case domain.ErrUserNameRequired,
		domain.ErrEmailRequired,
		domain.ErrPasswordRequired,
		domain.ErrEventNameRequired,
		domain.ErrLocationIDRequired,
		domain.ErrStartDateRequired,
		domain.ErrEndDateRequired,
		domain.ErrInvalidEventDates,
		domain.ErrLocationNotFound:
		writeError(w, http.StatusBadRequest, fmt.Sprintf("Bad Request error: %s", err.Error()))
	case domain.ErrInvalidCredentials:
		writeError(w, http.StatusUnauthorized, err.Error())
	case domain.ErrAdminRequired:
		writeError(w, http.StatusForbidden, err.Error())
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errorResponse{Error: msg})
}
