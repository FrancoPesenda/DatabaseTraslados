package login

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	domain "github.com/FrancoPesenda/eventra/internal/domain"
	"github.com/stretchr/testify/assert"
)

type mockUseCase struct {
	loginFn func(ctx context.Context, user domain.User) (domain.User, error)
}

func (m *mockUseCase) Login(ctx context.Context, user domain.User) (domain.User, error) {
	return m.loginFn(ctx, user)
}

func newTestHandler(uc UseCase) *Handler {
	return NewHandler(uc, log.Default())
}

func TestHandler_Login_WhenValidCredentials_ShouldReturn200WithBody(t *testing.T) {
	uc := &mockUseCase{
		loginFn: func(_ context.Context, _ domain.User) (domain.User, error) {
			return domain.User{
				ID:       1,
				UserName: "johndoe",
				Email:    "john@example.com",
				Role:     domain.CompanyRole,
			}, nil
		},
	}

	body, _ := json.Marshal(map[string]string{"user_name": "johndoe", "password": "secret"})
	r := httptest.NewRequest(http.MethodPost, "/user/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	newTestHandler(uc).Login(w, r)

	var response loginResponse
	_ = json.NewDecoder(w.Body).Decode(&response)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 1, response.ID)
	assert.Equal(t, "company", response.Role)
	assert.NotEmpty(t, response.Token)
}

func TestHandler_Login_WhenInvalidJSON_ShouldReturn400(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/user/login", bytes.NewReader([]byte("not json")))
	w := httptest.NewRecorder()

	newTestHandler(&mockUseCase{}).Login(w, r)

	var response errorResponse
	_ = json.NewDecoder(w.Body).Decode(&response)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "invalid request body", response.Error)
}

func TestHandler_Login_WhenUserNameIsEmpty_ShouldReturn400(t *testing.T) {
	uc := &mockUseCase{
		loginFn: func(_ context.Context, _ domain.User) (domain.User, error) {
			return domain.User{}, domain.ErrUserNameRequired
		},
	}

	body, _ := json.Marshal(map[string]string{"user_name": "", "password": "secret"})
	r := httptest.NewRequest(http.MethodPost, "/user/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	newTestHandler(uc).Login(w, r)

	var response errorResponse
	_ = json.NewDecoder(w.Body).Decode(&response)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, domain.ErrUserNameRequired.Error(), response.Error)
}

func TestHandler_Login_WhenPasswordIsEmpty_ShouldReturn400(t *testing.T) {
	uc := &mockUseCase{
		loginFn: func(_ context.Context, _ domain.User) (domain.User, error) {
			return domain.User{}, domain.ErrPasswordRequired
		},
	}

	body, _ := json.Marshal(map[string]string{"user_name": "johndoe", "password": ""})
	r := httptest.NewRequest(http.MethodPost, "/user/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	newTestHandler(uc).Login(w, r)

	var response errorResponse
	_ = json.NewDecoder(w.Body).Decode(&response)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, domain.ErrPasswordRequired.Error(), response.Error)
}

func TestHandler_Login_WhenInvalidCredentials_ShouldReturn401(t *testing.T) {
	uc := &mockUseCase{
		loginFn: func(_ context.Context, _ domain.User) (domain.User, error) {
			return domain.User{}, domain.ErrInvalidCredentials
		},
	}

	body, _ := json.Marshal(map[string]string{"user_name": "johndoe", "password": "wrong"})
	r := httptest.NewRequest(http.MethodPost, "/user/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	newTestHandler(uc).Login(w, r)

	var response errorResponse
	_ = json.NewDecoder(w.Body).Decode(&response)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, domain.ErrInvalidCredentials.Error(), response.Error)
}

func TestHandler_Login_WhenLoginByEmail_WhenValidCredentials_ShouldReturn200WithBody(t *testing.T) {
	uc := &mockUseCase{
		loginFn: func(_ context.Context, _ domain.User) (domain.User, error) {
			return domain.User{
				ID:       1,
				UserName: "johndoe",
				Email:    "john@example.com",
				Role:     domain.CompanyRole,
			}, nil
		},
	}

	body, _ := json.Marshal(map[string]string{"email": "john@example.com", "password": "secret"})
	r := httptest.NewRequest(http.MethodPost, "/user/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	newTestHandler(uc).Login(w, r)

	var response loginResponse
	_ = json.NewDecoder(w.Body).Decode(&response)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 1, response.ID)
	assert.Equal(t, "company", response.Role)
	assert.NotEmpty(t, response.Token)
}

func TestHandler_Login_WhenUnexpectedError_ShouldReturn500(t *testing.T) {
	uc := &mockUseCase{
		loginFn: func(_ context.Context, _ domain.User) (domain.User, error) {
			return domain.User{}, errors.New("connection refused")
		},
	}

	body, _ := json.Marshal(map[string]string{"user_name": "johndoe", "password": "secret"})
	r := httptest.NewRequest(http.MethodPost, "/user/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	newTestHandler(uc).Login(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
