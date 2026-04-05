package login

import domain "github.com/FrancoPesenda/eventra/internal/domain"

type request struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r request) toDomain() domain.User {
	return domain.User{
		UserName: r.UserName,
		Email:    r.Email,
		Password: r.Password,
	}
}
