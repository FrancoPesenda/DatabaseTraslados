package eventradatabase

import (
	"context"
	"database/sql/driver"
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	domain "github.com/FrancoPesenda/eventra/internal/domain"
	"github.com/stretchr/testify/assert"
)

type failingLastInsertID struct{}

func (failingLastInsertID) LastInsertId() (int64, error) {
	return 0, errors.New("last insert id failed")
}
func (failingLastInsertID) RowsAffected() (int64, error) { return 1, nil }

var _ driver.Result = failingLastInsertID{}

func TestRepository_CreateCompanyUser_WhenValidInput_ShouldReturnCreatedUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	expectedResponse := domain.User{
		ID:    1,
		Name:  "johndoe",
		Email: "john@example.com",
		Role:  domain.CompanyRole,
	}

	mock.ExpectExec("INSERT INTO user").
		WithArgs("johndoe", "john@example.com", "secret", string(domain.CompanyRole)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery("SELECT u.id").
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "name"}).
			AddRow(1, "johndoe", "john@example.com", "company"))

	repo := NewRepository(db)
	response, err := repo.CreateCompanyUser(context.Background(), domain.User{
		UserName: "johndoe",
		Email:    "john@example.com",
		Password: "secret",
	})

	assert.Nil(t, err)
	assert.Equal(t, expectedResponse, response)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestRepository_CreateCompanyUser_WhenInsertFails_ShouldReturnError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	expectedError := errors.New("insert user: connection refused")

	mock.ExpectExec("INSERT INTO user").
		WillReturnError(errors.New("connection refused"))

	repo := NewRepository(db)
	response, err := repo.CreateCompanyUser(context.Background(), domain.User{
		UserName: "johndoe",
		Email:    "john@example.com",
		Password: "secret",
	})

	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, response)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestRepository_CreateCompanyUser_WhenLastInsertIDFails_ShouldReturnError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	expectedError := errors.New("last insert id: last insert id failed")

	mock.ExpectExec("INSERT INTO user").
		WillReturnResult(failingLastInsertID{})

	repo := NewRepository(db)
	response, err := repo.CreateCompanyUser(context.Background(), domain.User{
		UserName: "johndoe",
		Email:    "john@example.com",
		Password: "secret",
	})

	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, response)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestRepository_CreateCompanyUser_WhenSelectFails_ShouldReturnError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	expectedError := errors.New("select user by id: sql: no rows in result set")

	mock.ExpectExec("INSERT INTO user").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery("SELECT u.id").
		WithArgs(int64(1)).
		WillReturnError(errors.New("sql: no rows in result set"))

	repo := NewRepository(db)
	response, err := repo.CreateCompanyUser(context.Background(), domain.User{
		UserName: "johndoe",
		Email:    "john@example.com",
		Password: "secret",
	})

	assert.EqualError(t, err, expectedError.Error())
	assert.Empty(t, response)
	assert.Nil(t, mock.ExpectationsWereMet())
}
