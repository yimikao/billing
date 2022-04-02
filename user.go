package cashapp

import (
	"time"

	"github.com/google/uuid"
)

type Model struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	Reference string    `json:"reference"`
	Tag       string    `json:"tag"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`

	Model
}
