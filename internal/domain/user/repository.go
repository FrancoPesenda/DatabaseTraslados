package user

import (
	"context"
)

type Repository interface {
	GetByID(ctx context.Context, id ID) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	Create(ctx context.Context, u User) (User, error)
	Update(ctx context.Context, u User) (User, error)
}

