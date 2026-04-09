package login

import (
	"context"
	"log"

	domain "github.com/FrancoPesenda/eventra/internal/domain"
	"github.com/FrancoPesenda/eventra/internal/utils/security"
)

type UserRepository interface {
	GetUserByUserName(ctx context.Context, userName string) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
}

type UseCase struct {
	repository UserRepository
}

func NewUseCase(repository UserRepository) *UseCase {
	return &UseCase{repository: repository}
}

func (u *UseCase) Login(ctx context.Context, user domain.User) (domain.User, error) {
	var (
		storedUser domain.User
		err        error
	)

	if err := validateCredentials(user); err != nil {
		log.Printf("[Layer:UseCase][error_message:%s][request_body:%+v]", err.Error(), user)
		return domain.User{}, err
	}

	if user.IsUserNameValid() {
		storedUser, err = u.repository.GetUserByUserName(ctx, user.UserName)
		if err != nil {
			log.Printf("[Layer:UseCase][error_message:%s][request_body:%+v]", err.Error(), user)
			return domain.User{}, domain.ErrInvalidCredentials
		}
	} else if user.IsEmailValid() {
		storedUser, err = u.repository.GetUserByEmail(ctx, user.Email)
		if err != nil {
			log.Printf("[Layer:UseCase][error_message:%s][request_body:%+v]", err.Error(), user)
			return domain.User{}, domain.ErrInvalidCredentials
		}
	}

	if err := security.CheckPassword(storedUser.Password, user.Password); err != nil {
		log.Printf("[Layer:UseCase][error_message:%s][request_body:%+v]", err.Error(), user)
		return domain.User{}, domain.ErrInvalidCredentials
	}

	storedUser.Password = ""
	return storedUser, nil
}

func validateCredentials(user domain.User) error {
	if !user.IsUserNameValid() && !user.IsEmailValid() {
		if !user.IsUserNameValid() {
			return domain.ErrUserNameRequired
		}
		if !user.IsEmailValid() {
			return domain.ErrEmailRequired
		}
		// capaz cambiar por  return errors.New("username or email is required")
		// ya que nunca verifica el email
	}

	if !user.IsPasswordValid() {
		return domain.ErrPasswordRequired
	}

	return nil
}
