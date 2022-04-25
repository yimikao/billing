package billing

import (
	"database/sql/driver"
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Tag             string `json:"tag"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	IsEmailVerified bool   `json:"is_email_verified"`
	Password        string `json:"password"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Model
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

// findUser
type FindUserOptions struct {
	ID    uuid.UUID
	Email Email
}
