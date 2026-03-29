package user

import "time"

type ID string

type User struct {
	ID        ID
	Email     string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

