package eventradatabase

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"

	domain "eventra/internal/domain"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateCompanyUser(ctx context.Context, user domain.User) (domain.User, error) {
	res, err := r.db.ExecContext(
		ctx,
		`INSERT INTO user (name, email, password, type)
		 VALUES (?, ?, ?, 'company')`,
		user.Name,
		user.Email,
		user.Password,
	)
	if err != nil {
		log.Printf("[Layer:Repository][error_message:%s][request_body:%+v]", err.Error(), user)
		return domain.User{}, fmt.Errorf("insert user: %w", err)
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		log.Printf("[Layer:Repository][error_message:%s][request_body:%+v]", err.Error(), user)
		return domain.User{}, fmt.Errorf("last insert id: %w", err)
	}

	var result domain.User
	err = r.db.QueryRowContext(
		ctx,
		`SELECT id, name, email FROM user WHERE id = ?`,
		lastID,
	).Scan(&result.ID, &result.Name, &result.Email)
	if err != nil {
		log.Printf("[Layer:Repository][error_message:%s][request_body:%+v]", err.Error(), user)
		return domain.User{}, fmt.Errorf("select user by id: %w", err)
	}

	return result, nil
}
