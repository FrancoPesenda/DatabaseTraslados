package eventradatabase

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	domain "database/internal/domain"
	_ "github.com/go-sql-driver/mysql"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, u domain.User) (domain.User, error) {
	res, err := r.db.ExecContext(
		ctx,
		`INSERT INTO user (name, last_name, email, password, type)
		 VALUES (?, ?, ?, ?, 'admin')`,
		u.Name,
		u.LastName,
		u.Email,
		u.Password,
	)
	if err != nil {
		log.Printf("[Layer:Repository][error_message:%s][request_body:%+v]", err.Error(), u)
		return domain.User{}, fmt.Errorf("insert user: %w", err)
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		log.Printf("[Layer:Repository][error_message:%s][request_body:%+v]", err.Error(), u)
		return domain.User{}, fmt.Errorf("last insert id: %w", err)
	}

	var result domain.User
	err = r.db.QueryRowContext(
		ctx,
		`SELECT id, name, last_name, email FROM user WHERE id = ?`,
		lastID,
	).Scan(&result.ID, &result.Name, &result.LastName, &result.Email)
	if err != nil {
		log.Printf("[Layer:Repository][error_message:%s][request_body:%+v]", err.Error(), u)
		return domain.User{}, fmt.Errorf("select user by id: %w", err)
	}

	return result, nil
}
