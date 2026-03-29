package payment_method

import "time"

type ID string
type UserID string

type PaymentMethod struct {
	ID        ID
	UserID    UserID
	Type      string
	Token     string
	Last4     string
	Brand     string
	ExpMonth  int
	ExpYear   int
	CreatedAt time.Time
	UpdatedAt time.Time
}

