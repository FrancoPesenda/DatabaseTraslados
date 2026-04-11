package create

import (
	"context"
	"log"

	domain "github.com/FrancoPesenda/eventra/internal/domain"
)

type UserRepository interface {
	GetUserByUserName(ctx context.Context, userName string) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	CreateEvent(ctx context.Context, event domain.Event) (domain.Event, error)
}

type UseCase struct {
	repository UserRepository
}

func NewUseCase(repository UserRepository) *UseCase {
	return &UseCase{repository: repository}
}

func (u *UseCase) CreateEvent(ctx context.Context, admin domain.User, event domain.Event) (domain.Event, error) {
	// El usuario ya viene autenticado del JWT, solo verificar que sea admin
	if admin.Role != domain.AdminRole {
		log.Printf("[Layer:UseCase][error_message: user is not admin][admin:%+v]", admin)
		return domain.Event{}, domain.ErrAdminRequired
	}

	if err := validateEvent(event); err != nil {
		log.Printf("[Layer:UseCase][error_message:%s][request_body:%+v]", err.Error(), event)
		return domain.Event{}, err
	}

	result, err := u.repository.CreateEvent(ctx, event)
	if err != nil {
		log.Printf("[Layer:UseCase][error_message:%s][request_body:%+v]", err.Error(), event)
		return domain.Event{}, err
	}

	return result, nil
}

func validateAdmin(admin domain.User) error {
	if !admin.IsUserNameValid() && !admin.IsEmailValid() {
		if !admin.IsUserNameValid() {
			return domain.ErrUserNameRequired
		}
		return domain.ErrEmailRequired
	}

	if !admin.IsPasswordValid() {
		return domain.ErrPasswordRequired
	}

	return nil
}

func (u *UseCase) getAdmin(ctx context.Context, admin domain.User) (domain.User, error) {
	if admin.IsUserNameValid() {
		return u.repository.GetUserByUserName(ctx, admin.UserName)
	}

	return u.repository.GetUserByEmail(ctx, admin.Email)
}

func validateEvent(event domain.Event) error {
	if !event.IsNameValid() {
		return domain.ErrEventNameRequired
	}

	if !event.IsLocationIDValid() {
		return domain.ErrLocationIDRequired
	}

	if !event.IsStartDateValid() {
		return domain.ErrStartDateRequired
	}

	if !event.IsEndDateValid() {
		return domain.ErrEndDateRequired
	}

	if !event.HasValidDateRange() {
		return domain.ErrInvalidEventDates
	}

	return nil
}
