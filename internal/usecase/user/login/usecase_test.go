package login

import (
	"context"
	"errors"
	"testing"

	domain "github.com/FrancoPesenda/eventra/internal/domain"
	"github.com/FrancoPesenda/eventra/internal/utils/security"
	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	getUserByUserNameFn func(ctx context.Context, userName string) (domain.User, error)
	getUserByEmailFn    func(ctx context.Context, email string) (domain.User, error)
}

func (m *mockRepository) GetUserByUserName(ctx context.Context, userName string) (domain.User, error) {
	return m.getUserByUserNameFn(ctx, userName)
}

func (m *mockRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	return m.getUserByEmailFn(ctx, email)
}

func hashedPassword(t *testing.T, plain string) string {
	t.Helper()
	h, err := security.HashPassword(plain)
	if err != nil {
		t.Fatal(err)
	}
	return h
}

func TestUseCase_User_Login_WhenValidCredentials_ShouldReturnUser(t *testing.T) {
	expectedResponse := domain.User{
		ID:       1,
		UserName: "johndoe",
		Email:    "john@example.com",
		Role:     domain.CompanyRole,
	}

	repo := &mockRepository{
		getUserByUserNameFn: func(_ context.Context, _ string) (domain.User, error) {
			return domain.User{
				ID:       1,
				UserName: "johndoe",
				Email:    "john@example.com",
				Password: hashedPassword(t, "secret"),
				Role:     domain.CompanyRole,
			}, nil
		},
	}
	uc := NewUseCase(repo)

	response, err := uc.Login(context.Background(), domain.User{
		UserName: "johndoe",
		Password: "secret",
	})

	assert.Nil(t, err)
	assert.Equal(t, expectedResponse, response)
}

func TestUseCase_User_Login_WhenLoginByEmail_WhenValidCredentials_ShouldReturnUser(t *testing.T) {
	expectedResponse := domain.User{
		ID:       1,
		UserName: "johndoe",
		Email:    "john@example.com",
		Role:     domain.CompanyRole,
	}

	repo := &mockRepository{
		getUserByEmailFn: func(_ context.Context, _ string) (domain.User, error) {
			return domain.User{
				ID:       1,
				UserName: "johndoe",
				Email:    "john@example.com",
				Password: hashedPassword(t, "secret"),
				Role:     domain.CompanyRole,
			}, nil
		},
	}
	uc := NewUseCase(repo)

	response, err := uc.Login(context.Background(), domain.User{
		Email:    "john@example.com",
		Password: "secret",
	})

	assert.Nil(t, err)
	assert.Equal(t, expectedResponse, response)
}

func TestUseCase_User_Login_WhenBothUserNameAndEmailAreEmpty_ShouldReturnErrUserNameRequired(t *testing.T) {
	expectedError := domain.ErrUserNameRequired

	uc := NewUseCase(&mockRepository{})

	response, err := uc.Login(context.Background(), domain.User{
		Password: "secret",
	})

	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, response)
}

func TestUseCase_User_Login_WhenPasswordIsEmpty_ShouldReturnErrPasswordRequired(t *testing.T) {
	expectedError := domain.ErrPasswordRequired

	uc := NewUseCase(&mockRepository{})

	response, err := uc.Login(context.Background(), domain.User{
		UserName: "johndoe",
	})

	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, response)
}

func TestUseCase_User_Login_WhenUserNotFound_ShouldReturnErrInvalidCredentials(t *testing.T) {
	expectedError := domain.ErrInvalidCredentials

	repo := &mockRepository{
		getUserByUserNameFn: func(_ context.Context, _ string) (domain.User, error) {
			return domain.User{}, errors.New("user not found")
		},
	}
	uc := NewUseCase(repo)

	response, err := uc.Login(context.Background(), domain.User{
		UserName: "johndoe",
		Password: "secret",
	})

	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, response)
}

func TestUseCase_User_Login_WhenLoginByEmail_WhenUserNotFound_ShouldReturnErrInvalidCredentials(t *testing.T) {
	expectedError := domain.ErrInvalidCredentials

	repo := &mockRepository{
		getUserByEmailFn: func(_ context.Context, _ string) (domain.User, error) {
			return domain.User{}, errors.New("user not found")
		},
	}
	uc := NewUseCase(repo)

	response, err := uc.Login(context.Background(), domain.User{
		Email:    "john@example.com",
		Password: "secret",
	})

	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, response)
}

func TestUseCase_User_Login_WhenPasswordDoesNotMatch_ShouldReturnErrInvalidCredentials(t *testing.T) {
	expectedError := domain.ErrInvalidCredentials

	repo := &mockRepository{
		getUserByUserNameFn: func(_ context.Context, _ string) (domain.User, error) {
			return domain.User{
				UserName: "johndoe",
				Password: hashedPassword(t, "secret"),
			}, nil
		},
	}
	uc := NewUseCase(repo)

	response, err := uc.Login(context.Background(), domain.User{
		UserName: "johndoe",
		Password: "wrongpassword",
	})

	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, response)
}

func TestUseCase_User_Login_WhenLoginByEmail_WhenPasswordDoesNotMatch_ShouldReturnErrInvalidCredentials(t *testing.T) {
	expectedError := domain.ErrInvalidCredentials

	repo := &mockRepository{
		getUserByEmailFn: func(_ context.Context, _ string) (domain.User, error) {
			return domain.User{
				Email:    "john@example.com",
				Password: hashedPassword(t, "secret"),
			}, nil
		},
	}
	uc := NewUseCase(repo)

	response, err := uc.Login(context.Background(), domain.User{
		Email:    "john@example.com",
		Password: "wrongpassword",
	})

	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, response)
}
