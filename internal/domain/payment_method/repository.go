package payment_method

import "context"

type Repository interface {
	GetByID(ctx context.Context, id ID) (PaymentMethod, error)
	ListByUser(ctx context.Context, userID UserID) ([]PaymentMethod, error)
	Create(ctx context.Context, pm PaymentMethod) (PaymentMethod, error)
	Update(ctx context.Context, pm PaymentMethod) (PaymentMethod, error)
}

