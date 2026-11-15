package admodels

import (
	"errors"
	"testing"
)

func TestNilAccount(t *testing.T) {
	var acc *Account

	if acc.ID() != 0 {
		t.Fatalf("ID = %d, want 0", acc.ID())
	}
	if acc.TestBudgetValue() {
		t.Fatal("TestBudgetValue = true, want false")
	}
	if err := acc.TestRequest(1); !errors.Is(err, ErrAccessPointNotAllowed) {
		t.Fatalf("TestRequest = %v, want ErrAccessPointNotAllowed", err)
	}
}
