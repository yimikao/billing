package billing

import (
	"database/sql/driver"
	"strings"
)

type User struct {
	Model

	Tag             string `json:"tag"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	IsEmailVerified bool   `json:"is_email_verified"`
	Password        string `json:"password"`
}

type Email string

func (e Email) String() string {
	return strings.TrimSpace(strings.ToLower(string(e)))
}

func (e Email) Value() (driver.Value, error) {
	return driver.Value(e.String()), nil
}

type UserRepository interface {
	Create(*User) error
	GetByReference(string) (*User, error)
}
