package user

import domain "database/internal/domain"

type createRequest struct {
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r createRequest) toDomain() domain.User {
	return domain.User{
		Name:     r.Name,
		LastName: r.LastName,
		Email:    r.Email,
		Password: r.Password,
	}
}
