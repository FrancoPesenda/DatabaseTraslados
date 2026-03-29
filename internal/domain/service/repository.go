package service

import "context"

type Repository interface {
	GetByID(ctx context.Context, id ID) (Service, error)
	ListByType(ctx context.Context, typeID TypeID) ([]Service, error)
	Create(ctx context.Context, s Service) (Service, error)
	Update(ctx context.Context, s Service) (Service, error)
}

