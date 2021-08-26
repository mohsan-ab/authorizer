package account

import (
	"time"

	"github.com/mohsanabbas/authorizer/commons/eventsource"
)

//CreateAccount
type CreateAccount struct {
	eventsource.BaseCommand
	ActiveCard     bool
	AvailableLimit uint64
	Violations     []string
}

//PerformTransaction to a given account
type PerformTransaction struct {
	eventsource.BaseCommand
	Amount   uint64
	Time     time.Time
	Merchant string
}
