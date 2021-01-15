package account

import (
	"testing"
)

func TestAccountNew(t *testing.T) {
	acc := NewAccount("Cash in Wallet", TypeAsset)

	if acc.Timestamp == 0 {
		t.Error("Timestamp was not initialized")
	}

	if acc.Alias == "" {
		t.Error("Alias was not initialized")
	}

	emptyID := ID{}
	if acc.ID == emptyID {
		t.Error("ID was not initialized")
	}
}
