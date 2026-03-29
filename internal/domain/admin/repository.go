package admin

import "context"

// Repository defines the persistence contract for the Admin aggregate.
type Repository interface {
	Create(ctx context.Context, a Admin) (Admin, error)
}
