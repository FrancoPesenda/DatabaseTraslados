package event

import "context"

type Repository interface {
	GetByID(ctx context.Context, id ID) (Event, error)
	ListByUser(ctx context.Context, userID UserID) ([]Event, error)
	Create(ctx context.Context, e Event) (Event, error)
	Update(ctx context.Context, e Event) (Event, error)
}

