package billing

import "github.com/google/uuid"

type Wallet struct {
	UserID    uuid.UUID `json:"user_id"`
	IsPrimary bool      `json:"is_primary"`
	Currency  string    `json:"currency"`

	Model
}
