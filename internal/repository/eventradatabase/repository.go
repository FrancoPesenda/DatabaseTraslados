package eventradatabase

import (
	"context"
	"database/sql"
	"fmt"

	"database/internal/usecase"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// CreateUser inserts a row into `users` and returns the created user.
func (r *Repository) CreateUser(ctx context.Context, email, name string) (usecase.User, error) {
	res, err := r.db.ExecContext(
		ctx,
		`INSERT INTO users (email, name) VALUES (?, ?)`,
		email,
		name,
	)
	if err != nil {
		return usecase.User{}, fmt.Errorf("insert user: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return usecase.User{}, fmt.Errorf("last insert id: %w", err)
	}

	var u usecase.User
	err = r.db.QueryRowContext(
		ctx,
		`SELECT id, email, name, created_at, updated_at
		 FROM users
		 WHERE id = ?`,
		uint64(id),
	).Scan(&u.ID, &u.Email, &u.Name, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return usecase.User{}, fmt.Errorf("select user by id: %w", err)
	}

	return u, nil
}

