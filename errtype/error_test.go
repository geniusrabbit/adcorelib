package errtype_test

import (
	"errors"
	"fmt"
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

func TestWrap_errorsIs(t *testing.T) {
	const sentinel = errtype.Error("response is skipped")
	cause := errors.New("rps limit exceeded")
	other := errors.New("other")

	wrapped := sentinel.Wrap(cause)
	if !errors.Is(wrapped, sentinel) {
		t.Fatal("errors.Is(wrapped, sentinel) = false, want true")
	}
	if !errors.Is(wrapped, cause) {
		t.Fatal("errors.Is(wrapped, cause) = false, want true")
	}
	if errors.Is(wrapped, other) {
		t.Fatal("errors.Is(wrapped, other) = true, want false")
	}
	if got, want := wrapped.Error(), "response is skipped: rps limit exceeded"; got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}

func TestWrap_nilReturnsSentinel(t *testing.T) {
	const sentinel = errtype.Error("response is skipped")

	got := sentinel.Wrap(nil)
	if got != sentinel {
		t.Fatalf("Wrap(nil) = %v (%T), want sentinel itself", got, got)
	}
}

func TestWrap_nested(t *testing.T) {
	const sentinel = errtype.Error("response is skipped")
	inner := errors.New("circuit open")
	cause := fmt.Errorf("rps limit exceeded: %w", inner)

	wrapped := sentinel.Wrap(cause)
	if !errors.Is(wrapped, sentinel) {
		t.Fatal("errors.Is(wrapped, sentinel) = false, want true")
	}
	if !errors.Is(wrapped, cause) {
		t.Fatal("errors.Is(wrapped, cause) = false, want true")
	}
	if !errors.Is(wrapped, inner) {
		t.Fatal("errors.Is(wrapped, inner) = false, want true")
	}
}
