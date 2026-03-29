package http

import (
	"encoding/json"
	nethttp "net/http"

	"database/internal/usecase"
)

// HealthHandler handles liveness/readiness checks.
type HealthHandler struct {
	uc *usecase.HealthUsecase
}

func NewHealthHandler(uc *usecase.HealthUsecase) *HealthHandler {
	return &HealthHandler{uc: uc}
}

// Get handles GET /health.
func (h *HealthHandler) Get(w nethttp.ResponseWriter, r *nethttp.Request) {
	status, err := h.uc.Execute(r.Context())
	if err != nil {
		w.WriteHeader(nethttp.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(status)
}
