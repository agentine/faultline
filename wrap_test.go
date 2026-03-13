package faultline

import (
	"errors"
	"fmt"
	"testing"
)

func TestWrap(t *testing.T) {
	base := errors.New("base error")
	wrapped := Wrap(base, "context")

	if wrapped == nil {
		t.Fatal("Wrap returned nil for non-nil error")
	}
	if got := wrapped.Error(); got != "context: base error" {
		t.Errorf("Error(): got %q, want %q", got, "context: base error")
	}
}

func TestWrapNil(t *testing.T) {
	if got := Wrap(nil, "context"); got != nil {
		t.Errorf("Wrap(nil, ...): got %v, want nil", got)
	}
}

func TestWrapCause(t *testing.T) {
	base := errors.New("base")
	wrapped := Wrap(base, "layer1")
	cause := Cause(wrapped)
	if cause != base {
		t.Errorf("Cause: got %v, want %v", cause, base)
	}
}

func TestWrapf(t *testing.T) {
	base := errors.New("base")
	wrapped := Wrapf(base, "context %d", 42)

	if wrapped == nil {
		t.Fatal("Wrapf returned nil for non-nil error")
	}
	if got := wrapped.Error(); got != "context 42: base" {
		t.Errorf("Error(): got %q, want %q", got, "context 42: base")
	}
}

func TestWrapfNil(t *testing.T) {
	if got := Wrapf(nil, "fmt %d", 1); got != nil {
		t.Errorf("Wrapf(nil, ...): got %v, want nil", got)
	}
}

func TestWithStack(t *testing.T) {
	base := errors.New("base")
	wrapped := WithStack(base)

	if wrapped == nil {
		t.Fatal("WithStack returned nil for non-nil error")
	}
	if got := wrapped.Error(); got != "base" {
		t.Errorf("Error(): got %q, want %q", got, "base")
	}

	type stackTracer interface {
		StackTrace() StackTrace
	}
	st, ok := wrapped.(stackTracer)
	if !ok {
		t.Fatal("WithStack result does not implement stackTracer")
	}
	if len(st.StackTrace()) == 0 {
		t.Error("StackTrace() returned empty slice")
	}
}

func TestWithStackNil(t *testing.T) {
	if got := WithStack(nil); got != nil {
		t.Errorf("WithStack(nil): got %v, want nil", got)
	}
}

func TestWithMessage(t *testing.T) {
	base := errors.New("base")
	wrapped := WithMessage(base, "context")

	if wrapped == nil {
		t.Fatal("WithMessage returned nil for non-nil error")
	}
	if got := wrapped.Error(); got != "context: base" {
		t.Errorf("Error(): got %q, want %q", got, "context: base")
	}
}

func TestWithMessageNil(t *testing.T) {
	if got := WithMessage(nil, "context"); got != nil {
		t.Errorf("WithMessage(nil, ...): got %v, want nil", got)
	}
}

func TestWithMessagef(t *testing.T) {
	base := errors.New("base")
	wrapped := WithMessagef(base, "ctx %d", 99)

	if wrapped == nil {
		t.Fatal("WithMessagef returned nil for non-nil error")
	}
	if got := wrapped.Error(); got != "ctx 99: base" {
		t.Errorf("Error(): got %q, want %q", got, "ctx 99: base")
	}
}

func TestWithMessagefNil(t *testing.T) {
	if got := WithMessagef(nil, "fmt %d", 1); got != nil {
		t.Errorf("WithMessagef(nil, ...): got %v, want nil", got)
	}
}

func TestCauseDeepChain(t *testing.T) {
	root := errors.New("root")
	l1 := Wrap(root, "layer1")
	l2 := Wrap(l1, "layer2")
	l3 := WithMessage(l2, "layer3")

	cause := Cause(l3)
	if cause != root {
		t.Errorf("Cause: got %v, want %v", cause, root)
	}
}

func TestCauseNil(t *testing.T) {
	if got := Cause(nil); got != nil {
		t.Errorf("Cause(nil): got %v, want nil", got)
	}
}

func TestCauseStdlibError(t *testing.T) {
	base := errors.New("stdlib")
	// stdlib error has no Cause() or Unwrap(), so Cause should return it directly.
	if got := Cause(base); got != base {
		t.Errorf("Cause(stdlib error): got %v, want %v", got, base)
	}
}

func TestWrapFormat(t *testing.T) {
	base := errors.New("base")
	wrapped := Wrap(base, "context")

	// %s
	if got := fmt.Sprintf("%s", wrapped); got != "context: base" {
		t.Errorf("%%s: got %q, want %q", got, "context: base")
	}

	// %v
	if got := fmt.Sprintf("%v", wrapped); got != "context: base" {
		t.Errorf("%%v: got %q, want %q", got, "context: base")
	}

	// %q
	if got := fmt.Sprintf("%q", wrapped); got != `"context: base"` {
		t.Errorf("%%q: got %q, want %q", got, `"context: base"`)
	}

	// %+v should be longer than just the message
	plusV := fmt.Sprintf("%+v", wrapped)
	if len(plusV) <= len("context: base") {
		t.Errorf("%%+v output should include stack trace info, got %q", plusV)
	}
}
