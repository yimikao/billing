package billing

import (
	"strings"

	"github.com/google/uuid"
)

type Email string

func (e Email) String() string {
	return strings.ToLower(string(e))
}

type User struct {
	Tag      string `json:"tag"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`

	Model
}

type UserRepository interface {
	Create(*User) error
}

// findUser
type FindUserOptions struct {
	ID    uuid.UUID
	Email Email
}
