package service

import "time"

type ID string
type TypeID string

type Service struct {
	ID          ID
	TypeID      TypeID
	Name        string
	Description string
	PriceCents  int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

