package faultline

import (
	"errors"
	"fmt"
	"testing"
)

// sentinel is a custom error type for errors.Is / errors.As testing.
type sentinel struct {
	msg string
}

func (s *sentinel) Error() string { return s.msg }

func TestErrorsIsWrap(t *testing.T) {
	base := &sentinel{"base"}
	wrapped := Wrap(base, "layer1")
	wrapped2 := Wrap(wrapped, "layer2")

	if !errors.Is(wrapped, base) {
		t.Error("errors.Is(Wrap(base), base) should be true")
	}
	if !errors.Is(wrapped2, base) {
		t.Error("errors.Is(Wrap(Wrap(base)), base) should be true")
	}
}

func TestErrorsAsWrap(t *testing.T) {
	base := &sentinel{"base"}
	wrapped := Wrap(base, "layer1")
	wrapped2 := Wrap(wrapped, "layer2")

	var target *sentinel
	if !errors.As(wrapped, &target) {
		t.Error("errors.As(Wrap(base), *sentinel) should be true")
	}
	if target != base {
		t.Errorf("errors.As target: got %v, want %v", target, base)
	}

	target = nil
	if !errors.As(wrapped2, &target) {
		t.Error("errors.As(Wrap(Wrap(base)), *sentinel) should be true")
	}
	if target != base {
		t.Errorf("errors.As target: got %v, want %v", target, base)
	}
}

func TestErrorsIsWithStack(t *testing.T) {
	base := &sentinel{"base"}
	wrapped := WithStack(base)

	if !errors.Is(wrapped, base) {
		t.Error("errors.Is(WithStack(base), base) should be true")
	}
}

func TestErrorsIsWithMessage(t *testing.T) {
	base := &sentinel{"base"}
	wrapped := WithMessage(base, "context")

	if !errors.Is(wrapped, base) {
		t.Error("errors.Is(WithMessage(base), base) should be true")
	}
}

func TestErrorsIsNew(t *testing.T) {
	err := New("test")
	// New errors are unique, so errors.Is with a different error should be false.
	other := errors.New("other")
	if errors.Is(err, other) {
		t.Error("errors.Is(New, other) should be false")
	}
}

func TestErrorsAsMixedChain(t *testing.T) {
	base := &sentinel{"root"}
	chain := WithMessage(Wrap(WithStack(base), "wrap"), "msg")

	var target *sentinel
	if !errors.As(chain, &target) {
		t.Error("errors.As through mixed chain should find sentinel")
	}
	if target != base {
		t.Errorf("errors.As target: got %v, want %v", target, base)
	}
}

func TestErrorsIsThroughFmtErrorf(t *testing.T) {
	base := &sentinel{"base"}
	// Test that Cause() also walks stdlib Unwrap chains.
	stdWrapped := fmt.Errorf("stdlib: %w", base)
	cause := Cause(stdWrapped)
	if cause != base {
		t.Errorf("Cause through fmt.Errorf: got %v, want %v", cause, base)
	}
}
