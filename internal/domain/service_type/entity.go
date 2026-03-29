package service_type

import "time"

type ID string

type ServiceType struct {
	ID          ID
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

