package login

import domain "github.com/FrancoPesenda/eventra/internal/domain"

type loginResponse struct {
	ID       int    `json:"id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

func newLoginResponse(u domain.User) loginResponse {
	return loginResponse{
		ID:       u.ID,
		UserName: u.UserName,
		Email:    u.Email,
		Role:     string(u.Role),
	}
}

type errorResponse struct {
	Error string `json:"error"`
}
