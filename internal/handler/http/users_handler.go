package http

import (
	"encoding/json"
	"errors"
	nethttp "net/http"

	"database/internal/usecase"
)

// UsersHandler handles HTTP requests for the admins/users resource.
type UsersHandler struct {
	createUC *usecase.CreateAdminUsecase
}

func NewUsersHandler(createUC *usecase.CreateAdminUsecase) *UsersHandler {
	return &UsersHandler{createUC: createUC}
}

type createAdminRequest struct {
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
}

type createAdminResponse struct {
	ID       string `json:"id"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
}

// Create handles POST /users — receives nombre+apellido, persists a new admin
// and responds 201 with the created resource including the generated ID.
func (h *UsersHandler) Create(w nethttp.ResponseWriter, r *nethttp.Request) {
	var req createAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(nethttp.StatusBadRequest)
		return
	}

	out, err := h.createUC.Execute(r.Context(), usecase.CreateAdminInput{
		Name:     req.Nombre,
		LastName: req.Apellido,
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
	_ = json.NewEncoder(w).Encode(createAdminResponse{
		ID:       out.ID,
		Nombre:   out.Name,
		Apellido: out.LastName,
	})
}
