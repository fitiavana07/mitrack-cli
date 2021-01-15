package account

import (
	"errors"
	"fmt"
)

// Enum code is inspired by https://blog.learngoprogramming.com/golang-const-type-enums-iota-bc4befd096d3.

// Type represents account types. The account type defines the effect
// of Debit(Dr) and Credit(Cr) to the account.
type Type uint8

// Account Types
const (
	// TypeAsset is the type of assets account.
	TypeAsset Type = 1

	// TypeLiability is the type of liabilities account.
	TypeLiability Type = 2

	// TypeEquity is the type of equity account.
	TypeEquity Type = 3

	// TypeExpense is the type of expenses account.
	TypeExpense Type = 4

	// TypeRevenue is the type of revenues account.
	TypeRevenue Type = 5
)

// IsValid returns wether the given value is a valid account Type.
func (t Type) IsValid() bool {
	switch t {
	case TypeAsset, TypeLiability, TypeEquity, TypeExpense, TypeRevenue:
		return true
	}
	return false
}

// String returns the string representation of an accout type.
// It will panic when called with an invalid type.
func (t Type) String() string {
	if !t.IsValid() {
		panic(fmt.Sprintf("trying to get a String of an invalid Account Type: %d", t))
	}

	stringMap := map[Type]string{
		TypeAsset:     "asset",
		TypeLiability: "liability",
		TypeEquity:    "equity",
		TypeExpense:   "expense",
		TypeRevenue:   "revenue",
	}

	return stringMap[t]
}

// Initial returns the Initial used for an account type:
// - A: Asset
// - L: Liability
// - O: (Owner's) Equity
// - E: Expense
// - R: Revenue
func (t Type) Initial() string {
	if !t.IsValid() {
		panic(fmt.Sprintf("trying to get a String of an invalid Account Type: %d", t))
	}

	initialMap := map[Type]string{
		TypeAsset:     "A",
		TypeLiability: "L",
		TypeEquity:    "O",
		TypeExpense:   "E",
		TypeRevenue:   "R",
	}
	return initialMap[t]
}

// TypeFromString returns the account type corresponding to a given string.
func TypeFromString(s string) (t Type, err error) {
	switch s {
	case "Asset", "A", "assets", "asset":
		t = TypeAsset
	case "Liability", "L", "liabilities", "liability":
		t = TypeLiability
	case "Equity", "O", "equity":
		t = TypeEquity
	case "Expense", "E", "expenses", "expense":
		t = TypeExpense
	case "Revenue", "R", "revenues", "revenue":
		t = TypeRevenue
	default:
		err = ErrInvalidTypeString
	}
	return
}

// ErrInvalidTypeString is returned in case of invalid account type string.
var ErrInvalidTypeString = errors.New("invalid account type string")
