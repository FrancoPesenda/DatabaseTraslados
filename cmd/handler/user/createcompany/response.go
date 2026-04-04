package user

import domain "github.com/FrancoPesenda/eventra/internal/domain"

type createResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func newCreateResponse(u domain.User) createResponse {
	return createResponse{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
		Role:  string(u.Role),
	}
}

type errorResponse struct {
	Error string `json:"error"`
}
