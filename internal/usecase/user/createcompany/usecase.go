package createcompany

import (
	"context"
	"log"
	"strings"

	domain "github.com/FrancoPesenda/eventra/internal/domain"
	"github.com/FrancoPesenda/eventra/internal/utils/security"
)

type (
	UserRepository interface {
		CreateCompanyUser(ctx context.Context, user domain.User) (domain.User, error)
	}

	UseCase struct {
		repository UserRepository
	}
)

func NewUseCase(repository UserRepository) *UseCase {
	return &UseCase{repository: repository}
}

func (u *UseCase) CreateCompanyUser(ctx context.Context, user domain.User) (domain.User, error) {
	user.Name = strings.TrimSpace(user.Name)

	if err := validateUser(user); err != nil {
		log.Printf("[Layer:UseCase][error_message:%s][request_body:%+v]", err.Error(), user)
		return domain.User{}, err
	}

	setCompanyRole(&user)

	hashed, err := security.HashPassword(user.Password)
	if err != nil {
		log.Printf("[Layer:UseCase][error_message:%s][request_body:%+v]", err.Error(), user)
		return domain.User{}, err
	}
	user.Password = hashed

	result, err := u.repository.CreateCompanyUser(ctx, user)
	if err != nil {
		log.Printf("[Layer:UseCase][error_message:%s][request_body:%+v]", err.Error(), user)
		return domain.User{}, err
	}

	return result, nil
}

func validateUser(user domain.User) error {
	if !user.IsUserNameValid() {
		return domain.ErrUserNameRequired
	}

	if !user.IsEmailValid() {
		return domain.ErrEmailRequired
	}

	if !user.IsPasswordValid() {
		return domain.ErrPasswordRequired
	}

	return nil
}

func setCompanyRole(user *domain.User) {
	user.Role = domain.CompanyRole
}
