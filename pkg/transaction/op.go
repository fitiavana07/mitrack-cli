package transaction

// Operation is the type comprising Debit(Dr) and Credit(Cr).
type Operation uint8

const (
	// OpDebit is the Debit(Dr) operation.
	OpDebit Operation = 1

	// OpCredit is the Credit(Cr) operation.
	OpCredit Operation = 2
)

func (o Operation) String() string {
	switch o {
	case OpDebit:
		return "debit"
	case OpCredit:
		return "credit"
	default:
		return "noop"
	}
}
