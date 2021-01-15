package account

import (
	"crypto/rand"
	"fmt"
	"reflect"
	"testing"
)

func TestIDDecodeID(t *testing.T) {
	t.Run("decode string into an ID", func(t *testing.T) {
		want := makeRandomID(t)

		idStr := fmt.Sprintf("%x", want)
		got, err := DecodeID(idStr)

		if err != nil {
			t.Errorf("got error %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %q want %q", got, want)
		}
	})
	t.Run("invalid string length", func(t *testing.T) {
		idStr := "56838956cecc009acc91a426738c402e0638f01da11b1638a4cde4f"
		_, err := DecodeID(idStr)
		if err != ErrInvalidStringLength {
			t.Errorf("got error %q, want %q", err, ErrInvalidStringLength)
		}
	})
	t.Run("invalid hex string", func(t *testing.T) {
		idStr := "56838956cecc009acc91a426738c402e0638f01da11b1638a4cde4fd4eeca10t"
		_, err := DecodeID(idStr)

		if err != ErrInvalidHexString {
			t.Errorf("got error %q, want %q", err, ErrInvalidHexString)
		}
	})
}

func makeRandomID(t testing.TB) ID {
	t.Helper()

	b := make([]byte, IDSize)
	_, err := rand.Read(b)
	if err != nil {
		t.Fatalf("failed to create random bytes: %v", err)
	}

	id := ID{}
	copy(id[:], b[:])
	return id
}
