package create

import (
	"context"
	"testing"

	domain "github.com/FrancoPesenda/eventra/internal/domain"
	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	getByUserNameFn func(ctx context.Context, userName string) (domain.User, error)
	getByEmailFn    func(ctx context.Context, email string) (domain.User, error)
	createEventFn   func(ctx context.Context, event domain.Event) (domain.Event, error)
}

func (m *mockRepository) GetUserByUserName(ctx context.Context, userName string) (domain.User, error) {
	return m.getByUserNameFn(ctx, userName)
}

func (m *mockRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	return m.getByEmailFn(ctx, email)
}

func (m *mockRepository) CreateEvent(ctx context.Context, event domain.Event) (domain.Event, error) {
	return m.createEventFn(ctx, event)
}

func TestUseCase_CreateEvent_WhenAdminCanCreateEvent(t *testing.T) {
	repo := &mockRepository{
		createEventFn: func(_ context.Context, event domain.Event) (domain.Event, error) {
			event.ID = 123
			return event, nil
		},
	}

	uc := NewUseCase(repo)

	// User already authenticated via JWT, just needs admin role
	result, err := uc.CreateEvent(context.Background(), domain.User{ID: 1, UserName: "admin", Role: domain.AdminRole}, domain.Event{
		Name:       "Event 1",
		LocationID: 1,
		StartDate:  "2026-05-01",
		EndDate:    "2026-05-02",
	})

	assert.NoError(t, err)
	assert.Equal(t, 123, result.ID)
	assert.Equal(t, "Event 1", result.Name)
}

func TestUseCase_CreateEvent_WhenAdminIsNotAdmin_ShouldReturnErrAdminRequired(t *testing.T) {
	repo := &mockRepository{}

	uc := NewUseCase(repo)

	// User with company role (not admin)
	_, err := uc.CreateEvent(context.Background(), domain.User{ID: 1, UserName: "company", Role: domain.CompanyRole}, domain.Event{
		Name:       "Event 1",
		LocationID: 1,
		StartDate:  "2026-05-01",
		EndDate:    "2026-05-02",
	})

	assert.ErrorIs(t, err, domain.ErrAdminRequired)
}

func TestUseCase_CreateEvent_WhenEventHasInvalidDates_ShouldReturnErrInvalidEventDates(t *testing.T) {
	repo := &mockRepository{}

	uc := NewUseCase(repo)

	_, err := uc.CreateEvent(context.Background(), domain.User{ID: 1, UserName: "admin", Role: domain.AdminRole}, domain.Event{
		Name:       "Event 1",
		LocationID: 1,
		StartDate:  "2026-05-05",
		EndDate:    "2026-05-01",
	})

	assert.ErrorIs(t, err, domain.ErrInvalidEventDates)
}

func TestUseCase_CreateEvent_WhenLocationNotFound_ShouldReturnErrLocationNotFound(t *testing.T) {
	repo := &mockRepository{
		createEventFn: func(_ context.Context, event domain.Event) (domain.Event, error) {
			return domain.Event{}, domain.ErrLocationNotFound
		},
	}

	uc := NewUseCase(repo)

	_, err := uc.CreateEvent(context.Background(), domain.User{ID: 1, UserName: "admin", Role: domain.AdminRole}, domain.Event{
		Name:       "Event 1",
		LocationID: 999, // Non-existent location
		StartDate:  "2026-05-01",
		EndDate:    "2026-05-02",
	})

	assert.ErrorIs(t, err, domain.ErrLocationNotFound)
}
