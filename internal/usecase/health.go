package usecase

import (
	"context"
	"time"
)

type HealthStatus struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

type HealthUsecase struct{}

func NewHealthUsecase() *HealthUsecase { return &HealthUsecase{} }

func (u *HealthUsecase) Execute(ctx context.Context) (HealthStatus, error) {
	_ = ctx
	return HealthStatus{
		Status:    "ok",
		Timestamp: time.Now().UTC(),
	}, nil
}

