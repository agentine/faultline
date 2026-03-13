package faultline

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestFormatNew(t *testing.T) {
	err := New("test message")

	tests := []struct {
		format string
		want   string
	}{
		{"%s", "test message"},
		{"%v", "test message"},
		{"%q", `"test message"`},
	}

	for _, tt := range tests {
		got := fmt.Sprintf(tt.format, err)
		if got != tt.want {
			t.Errorf("fmt.Sprintf(%q, New(...)): got %q, want %q", tt.format, got, tt.want)
		}
	}

	// %+v should contain both message and stack
	pv := fmt.Sprintf("%+v", err)
	if !strings.HasPrefix(pv, "test message") {
		t.Errorf("%%+v should start with message, got %q", pv)
	}
	if !strings.Contains(pv, "TestFormatNew") {
		t.Errorf("%%+v should contain function name, got %q", pv)
	}
}

func TestFormatWrap(t *testing.T) {
	base := errors.New("base error")
	err := Wrap(base, "wrapped")

	tests := []struct {
		format string
		want   string
	}{
		{"%s", "wrapped: base error"},
		{"%v", "wrapped: base error"},
		{"%q", `"wrapped: base error"`},
	}

	for _, tt := range tests {
		got := fmt.Sprintf(tt.format, err)
		if got != tt.want {
			t.Errorf("fmt.Sprintf(%q, Wrap(...)): got %q, want %q", tt.format, got, tt.want)
		}
	}

	// %+v should contain both messages and stack
	pv := fmt.Sprintf("%+v", err)
	if !strings.Contains(pv, "wrapped") {
		t.Errorf("%%+v should contain wrap message, got %q", pv)
	}
	if !strings.Contains(pv, "base error") {
		t.Errorf("%%+v should contain base message, got %q", pv)
	}
}

func TestFormatWithStack(t *testing.T) {
	base := errors.New("base")
	err := WithStack(base)

	if got := fmt.Sprintf("%s", err); got != "base" {
		t.Errorf("%%s: got %q, want %q", got, "base")
	}

	if got := fmt.Sprintf("%v", err); got != "base" {
		t.Errorf("%%v: got %q, want %q", got, "base")
	}

	pv := fmt.Sprintf("%+v", err)
	if !strings.Contains(pv, "base") {
		t.Errorf("%%+v should contain message, got %q", pv)
	}
	if !strings.Contains(pv, "TestFormatWithStack") {
		t.Errorf("%%+v should contain function name, got %q", pv)
	}
}

func TestFormatWithMessage(t *testing.T) {
	base := New("base")
	err := WithMessage(base, "context")

	if got := fmt.Sprintf("%s", err); got != "context: base" {
		t.Errorf("%%s: got %q, want %q", got, "context: base")
	}

	if got := fmt.Sprintf("%v", err); got != "context: base" {
		t.Errorf("%%v: got %q, want %q", got, "context: base")
	}

	pv := fmt.Sprintf("%+v", err)
	if !strings.Contains(pv, "context") {
		t.Errorf("%%+v should contain context message, got %q", pv)
	}
	if !strings.Contains(pv, "base") {
		t.Errorf("%%+v should contain base message, got %q", pv)
	}
}
