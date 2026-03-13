package faultline

import (
	"errors"
	"fmt"
	"testing"
)

var benchErr error

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchErr = New("error message")
	}
}

func BenchmarkErrorf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchErr = Errorf("error %d", i)
	}
}

func BenchmarkWrap(b *testing.B) {
	base := errors.New("base")
	for i := 0; i < b.N; i++ {
		benchErr = Wrap(base, "wrapped")
	}
}

func BenchmarkWrapf(b *testing.B) {
	base := errors.New("base")
	for i := 0; i < b.N; i++ {
		benchErr = Wrapf(base, "wrapped %d", i)
	}
}

func BenchmarkWithStack(b *testing.B) {
	base := errors.New("base")
	for i := 0; i < b.N; i++ {
		benchErr = WithStack(base)
	}
}

func BenchmarkWithMessage(b *testing.B) {
	base := errors.New("base")
	for i := 0; i < b.N; i++ {
		benchErr = WithMessage(base, "context")
	}
}

func BenchmarkCause(b *testing.B) {
	base := errors.New("base")
	chain := Wrap(Wrap(Wrap(base, "l1"), "l2"), "l3")
	for i := 0; i < b.N; i++ {
		Cause(chain)
	}
}

var benchSink string

func BenchmarkFormatPlusV(b *testing.B) {
	err := Wrap(New("base"), "context")
	for i := 0; i < b.N; i++ {
		benchSink = fmt.Sprintf("%+v", err)
	}
}
