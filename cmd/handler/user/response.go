package user

import domain "eventra/internal/domain"

type createResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
}

func newCreateResponse(u domain.User) createResponse {
	return createResponse{
		ID:       u.ID,
		Name:     u.Name,
		LastName: u.LastName,
		Email:    u.Email,
	}
}

type errorResponse struct {
	Error string `json:"error"`
}
