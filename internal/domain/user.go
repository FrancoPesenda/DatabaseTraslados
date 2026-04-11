package user

import "errors"

var (
	ErrNameRequired            = errors.New("field 'name' is required")
	ErrLastNameRequired        = errors.New("field 'last_name' is required")
	ErrUserNameRequired        = errors.New("field 'user_name' is required")
	ErrEmailRequired           = errors.New("field 'email' is required")
	ErrPasswordRequired        = errors.New("field 'password' is required")
	ErrInvalidCredentials      = errors.New("invalid credentials")
	ErrAdminRequired           = errors.New("admin role required")
	ErrUserNameorEmailRequired = errors.New("username or email is required")
)

const (
	AdminRole   Role = "admin"
	CompanyRole Role = "company"
	DefaultRole Role = "default"
)

type (
	Role string

	User struct {
		ID       int
		Name     string
		LastName string
		UserName string
		Email    string
		Password string
		Role     Role
	}
)

func (u *User) IsNameValid() bool {
	return u.Name != ""
}

func (u *User) IsLastNameValid() bool {
	return u.LastName != ""
}

func (u *User) IsUserNameValid() bool {
	return u.UserName != ""
}

func (u *User) IsEmailValid() bool {
	return u.Email != ""
}

func (u *User) IsPasswordValid() bool {
	return u.Password != ""
}
