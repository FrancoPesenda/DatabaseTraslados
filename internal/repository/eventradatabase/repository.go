package eventradatabase

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"

	domain "github.com/FrancoPesenda/eventra/internal/domain"
)

type DB interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	Close() error
}

type Repository struct {
	db DB
}

func NewRepository(db DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateCompanyUser(ctx context.Context, user domain.User) (domain.User, error) {
	lastID, err := r.insertNewUserCompany(ctx, user)
	if err != nil {
		return domain.User{}, err
	}

	return r.getUserByID(ctx, lastID, user)
}

func (r *Repository) insertNewUserCompany(ctx context.Context, user domain.User) (int64, error) {
	res, err := r.db.ExecContext(
		ctx,
		`INSERT INTO user (username, email, password, role_id)
		 VALUES (?, ?, ?, (SELECT id FROM role WHERE name = ?))`,
		user.UserName,
		user.Email,
		user.Password,
		domain.CompanyRole,
	)
	if err != nil {
		log.Printf("[Layer:Repository][error_message:%s][request_body:%+v]", err.Error(), user)
		return 0, fmt.Errorf("insert user: %w", err)
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		log.Printf("[Layer:Repository][error_message:%s][request_body:%+v]", err.Error(), user)
		return 0, fmt.Errorf("last insert id: %w", err)
	}

	return lastID, nil
}

func (r *Repository) getUserByID(ctx context.Context, id int64, user domain.User) (domain.User, error) {
	var result domain.User
	var roleName string

	err := r.db.QueryRowContext(
		ctx,
		`SELECT u.id, u.username, u.email, r.name
		 FROM user u
		 LEFT JOIN role r ON r.id = u.role_id
		 WHERE u.id = ?`,
		id,
	).Scan(&result.ID, &result.Name, &result.Email, &roleName)
	if err != nil {
		log.Printf("[Layer:Repository][error_message:%s][request_body:%+v]", err.Error(), user)
		return domain.User{}, fmt.Errorf("select user by id: %w", err)
	}

	result.Role = toRole(roleName)
	return result, nil
}

func toRole(name string) domain.Role {
	switch domain.Role(name) {
	case domain.CompanyRole:
		return domain.CompanyRole
	default:
		return domain.Role(name)
	}
}
