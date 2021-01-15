package transaction

import (
	"crypto/sha256"
	"time"
)

// Transaction represents a financial transaction in a double-entry
// bookkeeping system.
type Transaction interface {

	// Hash is the hash of transaction's data.
	// Typically, the hash is not recorded in the (blockchain) transactions database.
	Hash() [sha256.Size]byte

	// Timestamp is date of recording of the transaction.
	// It is obtained using time.Now().UTC().Unix().
	Timestamp() int64

	// Entries are the transaction entries. It is composed of debits and credits.
	// For reference purposes, entries may have simple indices, like c1, d1.
	Entries() []Entry

	// Note returns the description or reason of a transaction.
	Note() string

	// Height is the height of the transaction in the whole set of transactions.
	// (If in a blockchain, it may be just the index/height of the transaction
	// in the block.)
	// Height() int64
}

// NewFromMaps returns a new transaction from debits and credits map (alias->amount).
func NewFromMaps(note string, debits, credits map[string]int64) (Transaction, error) {

	// TODO validation: sum debits = sum credits

	// TODO test
	// lines := []Line{}

	// for alias, amount := range debits {
	// 	lines = append(lines, newLineFromAlias(alias, amount, OpDebit))
	// }
	// for alias, amount := range credits {
	// 	lines = append(lines, newLineFromAlias(alias, amount, OpCredit))
	// }

	// for alias, amount := range debits {
	// 	acc, err := mitrackCli.AccountDB().FindByAlias(alias)
	// 	if err != nil {
	// 		// maybe account not found, read error ...
	// 		return err
	// 	}

	// 	line := transaction.Line{
	// 		Operation: transaction.OpDebit,
	// 		Account:   acc,
	// 		Amount:    amount,
	// 	}
	// 	lines = append(lines, line)
	// }

	// for alias, amount := range options.credits {
	// 	acc, err := mitrackCli.AccountDB().FindByAlias(alias)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	line := transaction.Line{
	// 		Operation: transaction.OpCredit,
	// 		Account:   acc,
	// 		Amount:    amount,
	// 	}
	// 	lines = append(lines, line)
	// }
	return &transaction{
		timestamp: time.Now().UTC().Unix(),
	}, nil
}

type transaction struct {
	hash      [sha256.Size]byte
	timestamp int64
	entries   []Entry
	note      string
}

func (t *transaction) Hash() [sha256.Size]byte {
	return t.hash
}

func (t *transaction) Timestamp() int64 {
	return t.timestamp
}

func (t *transaction) Entries() []Entry {
	return t.entries
}

func (t *transaction) Note() string {
	return t.note
}
