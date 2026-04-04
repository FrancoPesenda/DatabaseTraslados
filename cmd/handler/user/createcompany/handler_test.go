package user

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
	createFn func(ctx context.Context, user domain.User) (domain.User, error)
}

func (m *mockUseCase) CreateCompanyUser(ctx context.Context, user domain.User) (domain.User, error) {
	return m.createFn(ctx, user)
}

func newTestHandler(uc UseCase) *CreateHandler {
	return NewCreateHandler(uc, log.Default())
}

func TestHandler_CreateCompanyUser_WhenValidRequest_ShouldReturn201WithBody(t *testing.T) {
	expectedResponse := createResponse{
		ID:    1,
		Name:  "johndoe",
		Email: "john@example.com",
		Role:  string(domain.CompanyRole),
	}

	uc := &mockUseCase{
		createFn: func(_ context.Context, user domain.User) (domain.User, error) {
			return domain.User{
				ID:    1,
				Name:  "johndoe",
				Email: "john@example.com",
				Role:  domain.CompanyRole,
			}, nil
		},
	}

	body, _ := json.Marshal(map[string]string{
		"user_name": "johndoe",
		"email":     "john@example.com",
		"password":  "secret",
	})

	r := httptest.NewRequest(http.MethodPost, "/user/company", bytes.NewReader(body))
	w := httptest.NewRecorder()

	newTestHandler(uc).Handle(w, r)

	var response createResponse
	_ = json.NewDecoder(w.Body).Decode(&response)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, expectedResponse, response)
}

func TestHandler_CreateCompanyUser_WhenInvalidJSON_ShouldReturn400(t *testing.T) {
	expectedResponse := errorResponse{Error: "invalid request body"}

	r := httptest.NewRequest(http.MethodPost, "/user/company", bytes.NewReader([]byte("not json")))
	w := httptest.NewRecorder()

	newTestHandler(&mockUseCase{}).Handle(w, r)

	var response errorResponse
	_ = json.NewDecoder(w.Body).Decode(&response)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, expectedResponse, response)
}

func TestHandler_CreateCompanyUser_WhenUseCaseReturnsErrNameRequired_ShouldReturn400(t *testing.T) {
	expectedResponse := errorResponse{Error: "Bad Request error: " + domain.ErrNameRequired.Error()}

	uc := &mockUseCase{
		createFn: func(_ context.Context, user domain.User) (domain.User, error) {
			return domain.User{}, domain.ErrNameRequired
		},
	}

	body, _ := json.Marshal(map[string]string{"user_name": "johndoe", "email": "john@example.com", "password": "secret"})
	r := httptest.NewRequest(http.MethodPost, "/user/company", bytes.NewReader(body))
	w := httptest.NewRecorder()

	newTestHandler(uc).Handle(w, r)

	var response errorResponse
	_ = json.NewDecoder(w.Body).Decode(&response)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, expectedResponse, response)
}

func TestHandler_CreateCompanyUser_WhenUseCaseReturnsErrEmailRequired_ShouldReturn400(t *testing.T) {
	expectedResponse := errorResponse{Error: "Bad Request error: " + domain.ErrEmailRequired.Error()}

	uc := &mockUseCase{
		createFn: func(_ context.Context, user domain.User) (domain.User, error) {
			return domain.User{}, domain.ErrEmailRequired
		},
	}

	body, _ := json.Marshal(map[string]string{"user_name": "johndoe", "email": "john@example.com", "password": "secret"})
	r := httptest.NewRequest(http.MethodPost, "/user/company", bytes.NewReader(body))
	w := httptest.NewRecorder()

	newTestHandler(uc).Handle(w, r)

	var response errorResponse
	_ = json.NewDecoder(w.Body).Decode(&response)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, expectedResponse, response)
}

func TestHandler_CreateCompanyUser_WhenUseCaseReturnsErrPasswordRequired_ShouldReturn400(t *testing.T) {
	expectedResponse := errorResponse{Error: "Bad Request error: " + domain.ErrPasswordRequired.Error()}

	uc := &mockUseCase{
		createFn: func(_ context.Context, user domain.User) (domain.User, error) {
			return domain.User{}, domain.ErrPasswordRequired
		},
	}

	body, _ := json.Marshal(map[string]string{"user_name": "johndoe", "email": "john@example.com", "password": "secret"})
	r := httptest.NewRequest(http.MethodPost, "/user/company", bytes.NewReader(body))
	w := httptest.NewRecorder()

	newTestHandler(uc).Handle(w, r)

	var response errorResponse
	_ = json.NewDecoder(w.Body).Decode(&response)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, expectedResponse, response)
}

func TestHandler_CreateCompanyUser_WhenUseCaseReturnsUnexpectedError_ShouldReturn500(t *testing.T) {
	uc := &mockUseCase{
		createFn: func(_ context.Context, user domain.User) (domain.User, error) {
			return domain.User{}, errors.New("connection refused")
		},
	}

	body, _ := json.Marshal(map[string]string{"user_name": "johndoe", "email": "john@example.com", "password": "secret"})
	r := httptest.NewRequest(http.MethodPost, "/user/company", bytes.NewReader(body))
	w := httptest.NewRecorder()

	newTestHandler(uc).Handle(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
