package account

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
)

// ID is the type for account IDs in the program.
type ID [IDSize]byte

const (
	// IDSize is the size of ID, the length of byte array.
	IDSize = sha256.Size

	// IDStringLength is the length of the hex string.
	IDStringLength = 64

	// IDShortLength is the length of the hex string returned by Short().
	IDShortLength = 8
)

// DecodeID decodes a string and returns an ID. It will be used
// when an external user want to search an account using its iD.
func DecodeID(s string) (ID, error) {
	if len(s) != IDStringLength {
		return ID{}, ErrInvalidStringLength
	}

	b, err := hex.DecodeString(s)
	if err != nil {
		return ID{}, ErrInvalidHexString
	}

	id := ID{}
	copy(id[:], b[:])
	return id, err
}

// Short returns a short hex representation of an ID.
func (id ID) Short() string {
	return fmt.Sprintf("%x", id[:IDShortLength/2])
}

// Hex returns the hex representation of an ID.
func (id ID) Hex() string {
	return fmt.Sprintf("%x", id)
}

var (
	// ErrInvalidStringLength indicates invalid length string provided to DecodeID.
	ErrInvalidStringLength = errors.New("invalid string length")

	// ErrInvalidHexString indicates invalid hex string provided
	ErrInvalidHexString = errors.New("invalid hex string")
)
