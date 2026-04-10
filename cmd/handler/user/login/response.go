package login

import (
	"github.com/FrancoPesenda/eventra/internal/utils/token"

	domain "github.com/FrancoPesenda/eventra/internal/domain"
)

type loginResponse struct {
	ID    int    `json:"id"`
	Token string `json:"token"`
	Role  string `json:"role"`
}

func newLoginResponse(u domain.User) (loginResponse, error) {
	jwtToken, err := token.GenerateToken(u)
	if err != nil {
		return loginResponse{}, err
	}

	return loginResponse{
		ID:    u.ID,
		Token: jwtToken,
		Role:  string(u.Role),
	}, nil
}

type errorResponse struct {
	Error string `json:"error"`
}
