package account

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"strings"
	"time"
)

// Account represents an a financial account in a double-entry
// accounting system
type Account struct {

	// ID is the identifier of the account, generated on creation.
	// It is obtained by calculating the hash of initial attributes
	// of the account using SHA256 algorithm.
	ID ID

	// Name is the name of the account.
	Name string

	// Alias is an alias to the account.
	Alias string

	// Description is the description of the account.
	Description string

	// Type is the type of the account.
	Type Type

	// ParentID is the ID of the parent of this account in the account tree.
	// We don't really need to declare parent here, because we will not need
	// to access it often.
	ParentID ID

	// Timestamp is the creation date of the account, in timestamp.
	// It is obtained using time.Now().UTC().Unix().
	Timestamp int64
}

// NewAccount returns a new initialized Account.
func NewAccount(name string, t Type) *Account {
	a := &Account{
		Name: name,
		Type: t,
	}

	a.Alias = strings.ReplaceAll(strings.ToLower(a.Name), " ", "-")
	a.Timestamp = time.Now().UTC().Unix()

	data := []interface{}{
		a.Name,
		a.Alias,
		a.Description,
		a.Type,
		a.ParentID,
		a.Timestamp,
	}

	all := &bytes.Buffer{}
	for _, value := range data {
		switch value.(type) {
		case string:
			all.Write([]byte(value.(string)))
		default:
			binary.Write(all, binary.LittleEndian, value)
		}
	}

	a.ID = sha256.Sum256(all.Bytes())

	return a
}

// TODO
// func NewWithParentID(name string, type Type, parentID ID) *Account
// func NewWithDescription() ...

func (a Account) String() string {
	// fmt.Sprintf("")
	return fmt.Sprintf(
		"Account(ID=%x,Name='%s',Alias='%s',Type='%x'",
		a.ID[:4],
		a.Name,
		a.Alias,
		a.Type,
	)
}
