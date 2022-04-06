package billing

import "github.com/google/uuid"

type TransactionType string

const (
	Credit TransactionType = "credit"
	Debit  TransactionType = "debit"
)

type TransactionEntry struct {
	Model

	TransactionID uuid.UUID       `json:"transaction_id"`
	WalletID      uuid.UUID       `json:"wallet_id"`
	Type          TransactionType `json:"type"`
	Amount        int64           `json:"amount"`
}
