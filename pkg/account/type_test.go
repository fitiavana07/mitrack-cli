package account

import "testing"

func TestTypeString(t *testing.T) {
	tt := []struct {
		accType    Type
		wantString string
	}{
		{TypeAsset, "asset"},
		{TypeLiability, "liability"},
		{TypeEquity, "equity"},
		{TypeExpense, "expense"},
		{TypeRevenue, "revenue"},
	}

	for _, tc := range tt {
		got := tc.accType.String()
		if got != tc.wantString {
			t.Errorf("wrong string, got %q, want %q", got, tc.wantString)
		}
	}
}

func TestTypeIsValid(t *testing.T) {
	tt := []struct {
		accType     Type
		wantIsValid bool
	}{
		{TypeAsset, true},
		{TypeLiability, true},
		{TypeEquity, true},
		{TypeExpense, true},
		{TypeRevenue, true},
		{0, false},
	}

	for _, tc := range tt {
		got := tc.accType.IsValid()
		if got != tc.wantIsValid {
			t.Errorf("got %v, want %v for int value %d", got, tc.wantIsValid, tc.accType)
		}
	}
}
