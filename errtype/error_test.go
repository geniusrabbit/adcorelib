package errtype_test

import (
	"errors"
	"testing"

	"github.com/geniusrabbit/adcorelib/errtype"
)

func TestWithMessage_errorsIs(t *testing.T) {
	const sentinel = errtype.Error("response is skipped")

	wrapped := sentinel.WithMessage("rps limit exceeded")
	if !errors.Is(wrapped, sentinel) {
		t.Fatalf("errors.Is(wrapped, sentinel) = false, want true")
	}
	if got, want := wrapped.Error(), "response is skipped: rps limit exceeded"; got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}

func TestWithMessage_emptyReturnsSentinel(t *testing.T) {
	const sentinel = errtype.Error("response is skipped")

	got := sentinel.WithMessage("")
	if got != sentinel {
		t.Fatalf("WithMessage(\"\") = %v (%T), want sentinel itself", got, got)
	}
}

func TestWithMessage_Message(t *testing.T) {
	const sentinel = errtype.Error("response is skipped")

	wrapped := sentinel.WithMessage("error circuit open")
	var wm interface{ Message() string }
	if !errors.As(wrapped, &wm) {
		t.Fatal("errors.As(wrapped, Message()) failed")
	}
	if got, want := wm.Message(), "error circuit open"; got != want {
		t.Fatalf("Message() = %q, want %q", got, want)
	}
}
