package faultline

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		msg  string
		want string
	}{
		{"", ""},
		{"hello", "hello"},
		{"foo bar", "foo bar"},
	}

	for _, tt := range tests {
		got := New(tt.msg)
		if got.Error() != tt.want {
			t.Errorf("New(%q): got %q, want %q", tt.msg, got.Error(), tt.want)
		}
	}
}

func TestNewHasStack(t *testing.T) {
	err := New("test")
	type stackTracer interface {
		StackTrace() StackTrace
	}
	st, ok := err.(stackTracer)
	if !ok {
		t.Fatal("New error does not implement stackTracer")
	}
	trace := st.StackTrace()
	if len(trace) == 0 {
		t.Fatal("StackTrace() returned empty slice")
	}
}

func TestErrorf(t *testing.T) {
	tests := []struct {
		format string
		args   []interface{}
		want   string
	}{
		{"hello %s", []interface{}{"world"}, "hello world"},
		{"count: %d", []interface{}{42}, "count: 42"},
		{"no args", nil, "no args"},
	}

	for _, tt := range tests {
		got := Errorf(tt.format, tt.args...)
		if got.Error() != tt.want {
			t.Errorf("Errorf(%q, %v): got %q, want %q", tt.format, tt.args, got.Error(), tt.want)
		}
	}
}

func TestErrorfHasStack(t *testing.T) {
	err := Errorf("test %d", 1)
	type stackTracer interface {
		StackTrace() StackTrace
	}
	st, ok := err.(stackTracer)
	if !ok {
		t.Fatal("Errorf error does not implement stackTracer")
	}
	trace := st.StackTrace()
	if len(trace) == 0 {
		t.Fatal("StackTrace() returned empty slice")
	}
}

func TestFundamentalCauseIsNil(t *testing.T) {
	err := New("test")
	type causer interface {
		Cause() error
	}
	c, ok := err.(causer)
	if !ok {
		t.Fatal("fundamental does not implement causer")
	}
	if c.Cause() != nil {
		t.Error("fundamental.Cause() should be nil")
	}
}

func TestFundamentalUnwrapIsNil(t *testing.T) {
	err := New("test")
	type unwrapper interface {
		Unwrap() error
	}
	u, ok := err.(unwrapper)
	if !ok {
		t.Fatal("fundamental does not implement unwrapper")
	}
	if u.Unwrap() != nil {
		t.Error("fundamental.Unwrap() should be nil")
	}
}

func TestNewFormat(t *testing.T) {
	err := New("test message")

	// %s
	if got := fmt.Sprintf("%s", err); got != "test message" {
		t.Errorf("%%s: got %q, want %q", got, "test message")
	}

	// %v
	if got := fmt.Sprintf("%v", err); got != "test message" {
		t.Errorf("%%v: got %q, want %q", got, "test message")
	}

	// %q
	if got := fmt.Sprintf("%q", err); got != `"test message"` {
		t.Errorf("%%q: got %q, want %q", got, `"test message"`)
	}

	// %+v should contain the message and stack trace
	plusV := fmt.Sprintf("%+v", err)
	if len(plusV) <= len("test message") {
		t.Errorf("%%+v output should include stack trace, got %q", plusV)
	}
}
