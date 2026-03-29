package event

import "time"

type ID string
type UserID string
type ServiceID string

type Event struct {
	ID        ID
	UserID    UserID
	ServiceID ServiceID
	StartsAt  time.Time
	EndsAt    time.Time
	Notes     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

