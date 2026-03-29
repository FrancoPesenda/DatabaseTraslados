package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"database/internal/domain/admin"
)

// ErrValidation is returned when input data fails validation.
var ErrValidation = errors.New("validation error")

// CreateAdminInput holds the data required to create a new admin.
type CreateAdminInput struct {
	Name     string
	LastName string
}

// CreateAdminUsecase orchestrates the creation of an admin.
type CreateAdminUsecase struct {
	repo admin.Repository
}

func NewCreateAdminUsecase(repo admin.Repository) *CreateAdminUsecase {
	return &CreateAdminUsecase{repo: repo}
}

// Execute validates the input, delegates to the repository and returns the
// persisted admin with the DB-assigned ID.
func (u *CreateAdminUsecase) Execute(ctx context.Context, in CreateAdminInput) (admin.Admin, error) {
	in.Name = strings.TrimSpace(in.Name)
	in.LastName = strings.TrimSpace(in.LastName)

	if in.Name == "" || in.LastName == "" {
		return admin.Admin{}, fmt.Errorf("%w: nombre y apellido son requeridos", ErrValidation)
	}

	return u.repo.Create(ctx, admin.Admin{
		Name:     in.Name,
		LastName: in.LastName,
	})
}
