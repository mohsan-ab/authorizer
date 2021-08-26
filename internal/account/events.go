package account

import "time"

//AccountCreated event
type AccountCreated struct {
	ActiveCard     bool     `json:"active-card"`
	AvailableLimit uint64   `json:"available-limit"`
	Violations     []string `json:"violations"`
}

//TransactionPerformed event
type TransactionPerformed struct {
	Amount   uint64    `json:"amount"`
	Time     time.Time `json:"time"`
	Merchant string    `json:"merchant"`
}
