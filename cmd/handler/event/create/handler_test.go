package create

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	domain "github.com/FrancoPesenda/eventra/internal/domain"
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

func TestHandler_CreateEvent_WhenValidRequest_ShouldReturn201WithBody(t *testing.T) {
	reqBody := request{
		AdminUserName: "admin",
		AdminPassword: "secret",
		Name:          "Festival",
		LocationID:    1,
		StartDate:     "2026-08-01",
		EndDate:       "2026-08-03",
	}
	body, _ := json.Marshal(reqBody)

	uc := &mockUseCase{
		createEventFn: func(_ context.Context, _ domain.User, event domain.Event) (domain.Event, error) {
			event.ID = 99
			return event, nil
		},
	}

	r := httptest.NewRequest(http.MethodPost, "/event", bytes.NewReader(body))
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
	w := httptest.NewRecorder()

	newTestHandler(&mockUseCase{}).Handle(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_CreateEvent_WhenAdminNotAuthorized_ShouldReturn403(t *testing.T) {
	uc := &mockUseCase{
		createEventFn: func(_ context.Context, _ domain.User, _ domain.Event) (domain.Event, error) {
			return domain.Event{}, domain.ErrAdminRequired
		},
	}

	reqBody := request{
		AdminUserName: "admin",
		AdminPassword: "secret",
		Name:          "Festival",
		LocationID:    1,
		StartDate:     "2026-08-01",
		EndDate:       "2026-08-03",
	}
	body, _ := json.Marshal(reqBody)

	r := httptest.NewRequest(http.MethodPost, "/event", bytes.NewReader(body))
	w := httptest.NewRecorder()

	newTestHandler(uc).Handle(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}
