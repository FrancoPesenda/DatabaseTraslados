package user

import "errors"

var (
	ErrNameRequired     = errors.New("field 'name' is required")
	ErrLastNameRequired = errors.New("field 'last_name' is required")
	ErrEmailRequired    = errors.New("field 'email' is required")
	ErrPasswordRequired = errors.New("field 'password' is required")
)

type User struct {
	ID       int
	Name     string
	LastName string
	Email    string
	Password string
}

func (u *User) IsNameValid() bool {
	return u.Name != ""
}

func (u *User) IsLastNameValid() bool {
	return u.LastName != ""
}

func (u *User) IsEmailValid() bool {
	return u.Email != ""
}

func (u *User) IsPasswordValid() bool {
	return u.Password != ""
}
