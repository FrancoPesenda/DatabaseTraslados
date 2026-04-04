package create

import (
	"context"
	"errors"
	"testing"

	domain "database/internal/domain"

	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	createFn func(ctx context.Context, user domain.User) (domain.User, error)
}

func (m *mockRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	return m.createFn(ctx, user)
}

func TestUseCase_User_Create_WhenValidInput_ShouldReturnCreatedUser(t *testing.T) {
	expectedResponse := domain.User{ID: 1, Name: "John", LastName: "Doe", Email: "john@example.com", Password: "secret"}

	repo := &mockRepository{
		createFn: func(_ context.Context, user domain.User) (domain.User, error) {
			user.ID = 1
			return user, nil
		},
	}
	uc := NewUseCase(repo)

	response, err := uc.Execute(context.Background(), domain.User{
		Name:     "John",
		LastName: "Doe",
		Email:    "john@example.com",
		Password: "secret",
	})

	assert.Nil(t, err)
	assert.Equal(t, expectedResponse, response)
}

func TestUseCase_User_Create_WhenNameHasWhitespace_ShouldTrimAndSucceed(t *testing.T) {
	expectedResponse := domain.User{Name: "John", LastName: "Doe", Email: "john@example.com", Password: "secret"}

	repo := &mockRepository{
		createFn: func(_ context.Context, user domain.User) (domain.User, error) {
			return user, nil
		},
	}
	uc := NewUseCase(repo)

	response, err := uc.Execute(context.Background(), domain.User{
		Name:     "  John  ",
		LastName: "  Doe  ",
		Email:    "john@example.com",
		Password: "secret",
	})

	assert.Nil(t, err)
	assert.Equal(t, expectedResponse, response)
}

func TestUseCase_User_Create_WhenNameIsEmpty_ShouldReturnErrNameRequired(t *testing.T) {
	expectedError := domain.ErrNameRequired

	uc := NewUseCase(&mockRepository{})

	response, err := uc.Execute(context.Background(), domain.User{
		LastName: "Doe",
		Email:    "john@example.com",
		Password: "secret",
	})

	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, response)
}

func TestUseCase_User_Create_WhenNameIsWhitespace_ShouldReturnErrNameRequired(t *testing.T) {
	expectedError := domain.ErrNameRequired

	uc := NewUseCase(&mockRepository{})

	response, err := uc.Execute(context.Background(), domain.User{
		Name:     "   ",
		LastName: "Doe",
		Email:    "john@example.com",
		Password: "secret",
	})

	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, response)
}

func TestUseCase_User_Create_WhenLastNameIsEmpty_ShouldReturnErrLastNameRequired(t *testing.T) {
	expectedError := domain.ErrLastNameRequired

	uc := NewUseCase(&mockRepository{})

	response, err := uc.Execute(context.Background(), domain.User{
		Name:     "John",
		Email:    "john@example.com",
		Password: "secret",
	})

	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, response)
}

func TestUseCase_User_Create_WhenLastNameIsWhitespace_ShouldReturnErrLastNameRequired(t *testing.T) {
	expectedError := domain.ErrLastNameRequired

	uc := NewUseCase(&mockRepository{})

	response, err := uc.Execute(context.Background(), domain.User{
		Name:     "John",
		LastName: "   ",
		Email:    "john@example.com",
		Password: "secret",
	})

	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, response)
}

func TestUseCase_User_Create_WhenEmailIsEmpty_ShouldReturnErrEmailRequired(t *testing.T) {
	expectedError := domain.ErrEmailRequired

	uc := NewUseCase(&mockRepository{})

	response, err := uc.Execute(context.Background(), domain.User{
		Name:     "John",
		LastName: "Doe",
		Password: "secret",
	})

	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, response)
}

func TestUseCase_User_Create_WhenPasswordIsEmpty_ShouldReturnErrPasswordRequired(t *testing.T) {
	expectedError := domain.ErrPasswordRequired

	uc := NewUseCase(&mockRepository{})

	response, err := uc.Execute(context.Background(), domain.User{
		Name:     "John",
		LastName: "Doe",
		Email:    "john@example.com",
	})

	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, response)
}

func TestUseCase_User_Create_WhenRepositoryFails_ShouldReturnRepositoryError(t *testing.T) {
	expectedError := errors.New("connection refused")

	repo := &mockRepository{
		createFn: func(_ context.Context, user domain.User) (domain.User, error) {
			return domain.User{}, expectedError
		},
	}
	uc := NewUseCase(repo)

	response, err := uc.Execute(context.Background(), domain.User{
		Name:     "John",
		LastName: "Doe",
		Email:    "john@example.com",
		Password: "secret",
	})

	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, response)
}
