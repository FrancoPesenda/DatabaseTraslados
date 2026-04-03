package http

import (
	"encoding/json"
	"errors"
	"log"
	nethttp "net/http"

	"database/internal/usecase"
)

// UsersHandler handles HTTP requests for the admins/users resource.
type UsersHandler struct {
	createUC *usecase.CreateAdminUsecase
	logger   *log.Logger
}

func NewUsersHandler(createUC *usecase.CreateAdminUsecase, logger *log.Logger) *UsersHandler {
	return &UsersHandler{createUC: createUC, logger: logger}
}

type createAdminRequest struct {
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type createAdminResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
}

type errorResponse struct {
	Error string `json:"error"`
}

// Create handles POST /users — receives name, last_name, email and password,
// persists a new admin and responds 201 with the created resource.
func (h *UsersHandler) Create(w nethttp.ResponseWriter, r *nethttp.Request) {
	h.logger.Printf("POST /users - request received from %s", r.RemoteAddr)

	var req createAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("POST /users - decode error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(nethttp.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{Error: "invalid request body"})
		return
	}

	if req.Name == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(nethttp.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{Error: "el campo 'name' es obligatorio"})
		return
	}
	if req.LastName == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(nethttp.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{Error: "el campo 'last_name' es obligatorio"})
		return
	}
	if req.Email == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(nethttp.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{Error: "el campo 'email' es obligatorio"})
		return
	}
	if req.Password == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(nethttp.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{Error: "el campo 'password' es obligatorio"})
		return
	}

	h.logger.Printf("POST /users - creating user: name=%q last_name=%q email=%q", req.Name, req.LastName, req.Email)

	out, err := h.createUC.Execute(r.Context(), usecase.CreateAdminInput{
		Name:     req.Name,
		LastName: req.LastName,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		h.logger.Printf("POST /users - error: %v", err)
		if errors.Is(err, usecase.ErrValidation) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(nethttp.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
			return
		}
		w.WriteHeader(nethttp.StatusInternalServerError)
		return
	}

	h.logger.Printf("POST /users - user created: id=%d", out.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(nethttp.StatusCreated)
	_ = json.NewEncoder(w).Encode(createAdminResponse{
		ID:       out.ID,
		Name:     out.Name,
		LastName: out.LastName,
		Email:    out.Email,
	})
}
