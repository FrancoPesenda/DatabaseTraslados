package user

import domain "eventra/internal/domain"

type request struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *request) toDomain() domain.User {
	return domain.User{
		Name:     r.Name,
		Email:    r.Email,
		Password: r.Password,
	}
}
