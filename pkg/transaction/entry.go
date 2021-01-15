package transaction

import "github.com/fitiavana07/mitrack/pkg/account"

// Entry is a transaction entry. It may be a debit or a credit entry.
type Entry interface {

	// Operation is the operation taken.
	Operation() Operation

	// AccountID is the ID of the concerned account.
	// The ID is referenced here for lazy loading.
	AccountID() account.ID

	// Amount is the amount of the operation.
	Amount() int64
}

// NewEntry creates a new transaction entry.
func NewEntry(op Operation, accID account.ID, amount int64) Entry {
	return &entry{op, accID, amount}
}

type entry struct {
	operation Operation
	accountID account.ID
	amount    int64
}

func (e *entry) Operation() Operation {
	return e.operation
}
func (e *entry) AccountID() account.ID {
	return e.accountID
}
func (e *entry) Amount() int64 {
	return e.amount
}
