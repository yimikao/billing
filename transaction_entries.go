package billing

import "github.com/google/uuid"

// type TransactionType string
type TransactionType uint8

const (
	// Credit TransactionType = "credit"
	// Debit  TransactionType = "debit"

	Credit TransactionType = iota
	Debit
)

type TransactionEntry struct {
	Model

	TransactionID uuid.UUID       `json:"transaction_id"`
	WalletID      uuid.UUID       `json:"wallet_id"`
	Type          TransactionType `json:"type"`
	Amount        int64           `json:"amount"`
}
