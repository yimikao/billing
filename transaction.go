package billing

import "github.com/google/uuid"

// type TransactionDirection string
// type TransactionStatus string
// type TransactionPurpose string

type TransactionDirection uint8
type TransactionStatus uint8
type TransactionPurpose uint8

const (
	PendingTransactionStatus TransactionStatus = iota
	SuccessfulTransactionStatus
	FailedTransactionStatus
)

const (
	IncomingTransactionDirection TransactionDirection = iota
	OutgoingTransactionDirection
)

const (
	TransferTransactionPurpose TransactionPurpose = iota
	DepositTransactionPurpose
	WithdrawalTransactionPurpose
	ReversalTransactionPurpose
)

// const (
// 	Incoming TransactionDirection = "incoming"
// 	Outgoing TransactionDirection = "outgoing"

// 	Failed  TransactionStatus = "failed"
// 	Pending TransactionStatus = "pending"
// 	Success TransactionStatus = "success"

// 	Transfer   TransactionPurpose = "transfer"
// 	Deposit    TransactionPurpose = "deposit"
// 	Withdrawal TransactionPurpose = "withdrawal"
// 	Reversal   TransactionPurpose = "reversal"
// )

type WalletTransaction struct {
	Model

	From     uuid.UUID `json:"from"`
	To       uuid.UUID `json:"to"`
	WalletID uuid.UUID `json:"wallet_id"`
	Amount   int64     `json:"amount"`
	Currency string    `json:"currency"`

	Status        TransactionStatus    `json:"status"`
	FailureReason string               `json:"failure_reason"`
	Direction     TransactionDirection `json:"direction"`
	Purpose       TransactionPurpose   `json:"purpose"`
}
