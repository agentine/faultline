package faultline

import (
	"fmt"
	"strings"
	"testing"
)

func TestStackTrace(t *testing.T) {
	err := New("test")
	type stackTracer interface {
		StackTrace() StackTrace
	}
	st, ok := err.(stackTracer)
	if !ok {
		t.Fatal("error does not implement stackTracer")
	}

	trace := st.StackTrace()
	if len(trace) == 0 {
		t.Fatal("empty stack trace")
	}

	// First frame should be in this test file
	f := trace[0]
	file := f.file()
	if !strings.HasSuffix(file, "stack_test.go") {
		t.Errorf("first frame file: got %q, want suffix stack_test.go", file)
	}

	line := f.line()
	if line == 0 {
		t.Error("first frame line is 0")
	}

	name := f.name()
	if !strings.Contains(name, "TestStackTrace") {
		t.Errorf("first frame name: got %q, want to contain TestStackTrace", name)
	}
}

func TestFrameFormat(t *testing.T) {
	err := New("test")
	type stackTracer interface {
		StackTrace() StackTrace
	}
	st := err.(stackTracer)
	trace := st.StackTrace()
	f := trace[0]

	// %s should be just the filename
	s := fmt.Sprintf("%s", f)
	if !strings.HasSuffix(s, "stack_test.go") {
		t.Errorf("%%s: got %q", s)
	}

	// %d should be a line number
	d := fmt.Sprintf("%d", f)
	if d == "0" || d == "" {
		t.Errorf("%%d: got %q", d)
	}

	// %n should be function name without package path
	n := fmt.Sprintf("%n", f)
	if n != "TestFrameFormat" {
		t.Errorf("%%n: got %q, want TestFrameFormat", n)
	}

	// %v should be file:line
	v := fmt.Sprintf("%v", f)
	if !strings.Contains(v, "stack_test.go:") {
		t.Errorf("%%v: got %q", v)
	}

	// %+v should include function name and full path
	pv := fmt.Sprintf("%+v", f)
	if !strings.Contains(pv, "TestFrameFormat") {
		t.Errorf("%%+v: got %q, want to contain TestFrameFormat", pv)
	}
	if !strings.Contains(pv, "\n\t") {
		t.Errorf("%%+v: got %q, want to contain newline+tab", pv)
	}
}

func TestStackTraceFormat(t *testing.T) {
	err := New("test")
	type stackTracer interface {
		StackTrace() StackTrace
	}
	st := err.(stackTracer)
	trace := st.StackTrace()

	// %+v should produce multi-line output
	pv := fmt.Sprintf("%+v", trace)
	if !strings.Contains(pv, "TestStackTraceFormat") {
		t.Errorf("%%+v: got %q, want to contain TestStackTraceFormat", pv)
	}

	// %v should produce a bracketed list
	v := fmt.Sprintf("%v", trace)
	if !strings.HasPrefix(v, "[") || !strings.HasSuffix(v, "]") {
		t.Errorf("%%v: got %q, want bracketed list", v)
	}
}
