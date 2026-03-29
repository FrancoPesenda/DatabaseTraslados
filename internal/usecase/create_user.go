package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

type CreateUserInput struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type User struct {
	ID        uint64    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var ErrValidation = errors.New("validation error")

type UserRepository interface {
	CreateUser(ctx context.Context, email, name string) (User, error)
}

type CreateUserUsecase struct {
	repo UserRepository
}

func NewCreateUserUsecase(repo UserRepository) *CreateUserUsecase {
	return &CreateUserUsecase{repo: repo}
}

func (u *CreateUserUsecase) Execute(ctx context.Context, in CreateUserInput) (User, error) {
	in.Email = strings.TrimSpace(in.Email)
	in.Name = strings.TrimSpace(in.Name)

	if in.Email == "" || in.Name == "" {
		return User{}, fmt.Errorf("%w: email and name are required", ErrValidation)
	}

	return u.repo.CreateUser(ctx, in.Email, in.Name)
}

