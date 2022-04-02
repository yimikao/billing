package main

import "github.com/google/uuid"

type TransactionDirection string
type TransactionStatus string
type TransactionPurpose string

const (
	Incoming TransactionDirection = "incoming"
	Outgoing TransactionDirection = "outgoing"

	Failed  TransactionStatus = "failed"
	Pending TransactionStatus = "pending"
	Success TransactionStatus = "success"

	Transfer   TransactionPurpose = "transfer"
	Deposit    TransactionPurpose = "deposit"
	Withdrawal TransactionPurpose = "withdrawal"
	Reversal   TransactionPurpose = "reversal"
)

type Transaction struct {
	Model

	From     uuid.UUID `json:"from"`
	To       uuid.UUID `json:"to"`
	WalletID uuid.UUID `json:"wallet_id"`
	Amount   int64     `json:"amount"`

	Status        TransactionStatus    `json:"status"`
	FailureReason string               `json:"failure_reason"`
	Direction     TransactionDirection `json:"direction"`
	Purpose       TransactionPurpose   `json:"purpose"`
}
