package account

import (
	"fmt"

	"github.com/mohsanabbas/authorizer/commons/eventsource"
)

//ErrBalanceOut when you don't have sufficient limit to perform the operation
var (
	ErrBalanceOut        = "insufficient-limit"
	ErrCardInactive      = "card-not-active"
	ErrHighFrequency     = "high-frequency-small-interval"
	ErrDoubleTransaction = "doubled-transaction"
)

//Account of account
type Account struct {
	eventsource.BaseAggregate
	ActiveCard     bool
	AvailableLimit uint64
	Violations     []string
}

//ApplyChange to account
func (a *Account) ApplyChange(event eventsource.Event) {
	switch e := event.Data.(type) {
	case *AccountCreated:
		a.ActiveCard = e.ActiveCard
		a.AvailableLimit = e.AvailableLimit
		a.ID = event.AggregateID
		a.Violations=e.Violations
	case *TransactionPerformed:
		if len(a.Violations) == 0 {
			a.AvailableLimit -= e.Amount
		}
	}
	fmt.Printf("\naccount_output:%v\n",a)
}

//HandleCommand create events and validate based on such command
func (a *Account) HandleCommand(command eventsource.Command) error {
	event := eventsource.Event{
		AggregateID:   a.ID,
		AggregateType: "Account",
	}
	switch c := command.(type) {
	case CreateAccount:
		event.AggregateID = c.AggregateID
		event.Data = &AccountCreated{
			ActiveCard:     c.ActiveCard,
			AvailableLimit: c.AvailableLimit,
			Violations:     c.Violations,
		}

	case PerformTransaction:
		checkViolations(a, c, event)
		event.Data = &TransactionPerformed{
			Amount:   c.Amount,
			Time:     c.Time,
			Merchant: c.Merchant,
		}
	}

	a.BaseAggregate.ApplyChangeHelper(a, event, true)
	return nil
}

// checkViolations perform Violations validation
func checkViolations(ac *Account, c PerformTransaction, event eventsource.Event) {
	switch true {
	case ac.AvailableLimit < c.Amount:
		ac.Violations = append(ac.Violations, ErrBalanceOut)
	case !ac.ActiveCard:
		ac.Violations = append(ac.Violations, ErrCardInactive)
	case c.AggregateID != event.AggregateID:
		ac.Violations = append(ac.Violations, ErrDoubleTransaction)
	}
}
