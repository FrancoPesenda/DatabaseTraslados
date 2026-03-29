package service_type

import "context"

type Repository interface {
	GetByID(ctx context.Context, id ID) (ServiceType, error)
	List(ctx context.Context) ([]ServiceType, error)
	Create(ctx context.Context, st ServiceType) (ServiceType, error)
	Update(ctx context.Context, st ServiceType) (ServiceType, error)
}

