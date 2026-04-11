package create

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	domain "github.com/FrancoPesenda/eventra/internal/domain"
	"github.com/FrancoPesenda/eventra/internal/utils/token"
	"github.com/stretchr/testify/assert"
)

type mockUseCase struct {
	createEventFn func(ctx context.Context, admin domain.User, event domain.Event) (domain.Event, error)
}

func (m *mockUseCase) CreateEvent(ctx context.Context, admin domain.User, event domain.Event) (domain.Event, error) {
	return m.createEventFn(ctx, admin, event)
}

func newTestHandler(uc UseCase) *CreateHandler {
	return NewCreateHandler(uc, log.New(io.Discard, "", 0))
}

func generateTestToken(t *testing.T, role domain.Role) string {
	t.Helper()
	user := domain.User{
		ID:       1,
		UserName: "testadmin",
		Email:    "admin@test.com",
		Role:     role,
	}
	tok, err := token.GenerateToken(user)
	assert.NoError(t, err)
	return tok
}

func TestHandler_CreateEvent_WhenValidRequest_ShouldReturn201WithBody(t *testing.T) {
	reqBody := request{
		Name:       "Festival",
		LocationID: 1,
		StartDate:  "2026-08-01",
		EndDate:    "2026-08-03",
	}
	body, _ := json.Marshal(reqBody)

	uc := &mockUseCase{
		createEventFn: func(_ context.Context, _ domain.User, event domain.Event) (domain.Event, error) {
			event.ID = 99
			return event, nil
		},
	}

	r := httptest.NewRequest(http.MethodPost, "/event", bytes.NewReader(body))
	r.Header.Set("Authorization", "Bearer "+generateTestToken(t, domain.AdminRole))
	w := httptest.NewRecorder()

	newTestHandler(uc).Handle(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp createResponse
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, 99, resp.ID)
	assert.Equal(t, "Festival", resp.Name)
}

func TestHandler_CreateEvent_WhenInvalidJSON_ShouldReturn400(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/event", bytes.NewReader([]byte("not json")))
	r.Header.Set("Authorization", "Bearer "+generateTestToken(t, domain.AdminRole))
	w := httptest.NewRecorder()

	newTestHandler(&mockUseCase{}).Handle(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response errorResponse
	assert.NoError(t, json.NewDecoder(w.Body).Decode(&response))
	assert.Equal(t, "invalid request body", response.Error)
}

func TestHandler_CreateEvent_WhenMissingToken_ShouldReturn401(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/event", bytes.NewReader(mustJSON(t, request{
		Name:       "Festival",
		LocationID: 1,
		StartDate:  "2026-08-01",
		EndDate:    "2026-08-03",
	})))
	// No Authorization header
	w := httptest.NewRecorder()

	newTestHandler(&mockUseCase{}).Handle(w, r)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestHandler_CreateEvent_WhenAdminNotAuthorized_ShouldReturn403(t *testing.T) {
	uc := &mockUseCase{
		createEventFn: func(_ context.Context, _ domain.User, _ domain.Event) (domain.Event, error) {
			return domain.Event{}, domain.ErrAdminRequired
		},
	}

	r := httptest.NewRequest(http.MethodPost, "/event", bytes.NewReader(mustJSON(t, request{
		Name:       "Festival",
		LocationID: 1,
		StartDate:  "2026-08-01",
		EndDate:    "2026-08-03",
	})))
	r.Header.Set("Authorization", "Bearer "+generateTestToken(t, domain.CompanyRole))
	w := httptest.NewRecorder()

	newTestHandler(uc).Handle(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)

	var response errorResponse
	assert.NoError(t, json.NewDecoder(w.Body).Decode(&response))
	assert.Equal(t, domain.ErrAdminRequired.Error(), response.Error)
}

func TestHandler_CreateEvent_WhenUnexpectedError_ShouldReturn500(t *testing.T) {
	uc := &mockUseCase{
		createEventFn: func(_ context.Context, _ domain.User, _ domain.Event) (domain.Event, error) {
			return domain.Event{}, errors.New("repository failure")
		},
	}

	r := httptest.NewRequest(http.MethodPost, "/event", bytes.NewReader(mustJSON(t, request{
		Name:       "Festival",
		LocationID: 1,
		StartDate:  "2026-08-01",
		EndDate:    "2026-08-03",
	})))
	r.Header.Set("Authorization", "Bearer "+generateTestToken(t, domain.AdminRole))
	w := httptest.NewRecorder()

	newTestHandler(uc).Handle(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func mustJSON(t *testing.T, v interface{}) []byte {
	t.Helper()

	body, err := json.Marshal(v)
	assert.NoError(t, err)
	return body
}
