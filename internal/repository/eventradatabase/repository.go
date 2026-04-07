package eventradatabase

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"

	domain "github.com/FrancoPesenda/eventra/internal/domain"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
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

func (r *Repository) CreateEvent(ctx context.Context, event domain.Event) (domain.Event, error) {
	lastID, err := r.insertEvent(ctx, event)
	if err != nil {
		return domain.Event{}, err
	}

	return r.getEventByID(ctx, lastID)
}

func (r *Repository) insertEvent(ctx context.Context, event domain.Event) (int64, error) {
	res, err := r.db.ExecContext(
		ctx,
		`INSERT INTO event (name, location_id, start_date, end_date, image)
		 VALUES (?, ?, ?, ?, ?)`,
		event.Name,
		event.LocationID,
		event.StartDate,
		event.EndDate,
		event.Image,
	)
	if err != nil {
		log.Printf("[Layer:Repository][error_message:%s][request_body:%+v]", err.Error(), event)
		return 0, fmt.Errorf("insert event: %w", err)
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		log.Printf("[Layer:Repository][error_message:%s][request_body:%+v]", err.Error(), event)
		return 0, fmt.Errorf("last insert id: %w", err)
	}

	return lastID, nil
}

func (r *Repository) getEventByID(ctx context.Context, id int64) (domain.Event, error) {
	var result domain.Event

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, name, location_id, start_date, end_date, image
		 FROM event
		 WHERE id = ?`,
		id,
	).Scan(&result.ID, &result.Name, &result.LocationID, &result.StartDate, &result.EndDate, &result.Image)
	if err != nil {
		log.Printf("[Layer:Repository][error_message:%s][event_id:%d]", err.Error(), id)
		return domain.Event{}, fmt.Errorf("select event by id: %w", err)
	}

	return result, nil
}

func (r *Repository) GetUserByUserName(ctx context.Context, userName string) (domain.User, error) {
	var result domain.User
	var roleName string

	err := r.db.QueryRowContext(
		ctx,
		`SELECT u.id, u.username, u.email, u.password, r.name
		 FROM user u
		 LEFT JOIN role r ON r.id = u.role_id
		 WHERE u.username = ?`,
		userName,
	).Scan(&result.ID, &result.UserName, &result.Email, &result.Password, &roleName)
	if err != nil {
		log.Printf("[Layer:Repository][error_message:%s][request_body:%+v]", err.Error(), userName)
		return domain.User{}, fmt.Errorf("select user by username: %w", err)
	}

	result.Role = toRole(roleName)
	return result, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	var result domain.User
	var roleName string

	err := r.db.QueryRowContext(
		ctx,
		`SELECT u.id, u.username, u.email, u.password, r.name
		 FROM user u
		 LEFT JOIN role r ON r.id = u.role_id
		 WHERE u.email = ?`,
		email,
	).Scan(&result.ID, &result.UserName, &result.Email, &result.Password, &roleName)
	if err != nil {
		log.Printf("[Layer:Repository][error_message:%s][request_body:%+v]", err.Error(), email)
		return domain.User{}, fmt.Errorf("select user by email: %w", err)
	}

	result.Role = toRole(roleName)
	return result, nil
}

func toRole(name string) domain.Role {
	switch domain.Role(name) {
	case domain.CompanyRole:
		return domain.CompanyRole
	case domain.AdminRole:
		return domain.AdminRole
	default:
		return domain.Role(name)
	}
}
