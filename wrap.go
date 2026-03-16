package faultline

import (
	"fmt"
)

// Wrap returns an error annotating err with a stack trace at the point
// Wrap is called, and the supplied message. If err is nil, Wrap returns nil.
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return &withCause{
		cause: err,
		msg:   message,
		stack: callers(),
	}
}

// Wrapf returns an error annotating err with a stack trace at the point
// Wrapf is called, and the format specifier. If err is nil, Wrapf returns nil.
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &withCause{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
		stack: callers(),
	}
}

// WithStack annotates err with a stack trace at the point WithStack was called.
// If err is nil, WithStack returns nil.
func WithStack(err error) error {
	if err == nil {
		return nil
	}
	return &withStack{
		cause: err,
		stack: callers(),
	}
}

// WithMessage annotates an error with a new message.
// If err is nil, WithMessage returns nil.
func WithMessage(err error, message string) error {
	if err == nil {
		return nil
	}
	return &withMessage{
		cause: err,
		msg:   message,
	}
}

// WithMessagef annotates an error with the format specifier.
// If err is nil, WithMessagef returns nil.
func WithMessagef(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &withMessage{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}
}

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the following interface:
//
//	type causer interface {
//	       Cause() error
//	}
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func Cause(err error) error {
	type causer interface {
		Cause() error
	}
	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			// Also try Unwrap for compatibility with stdlib errors.
			unwrapper, ok := err.(interface{ Unwrap() error })
			if !ok {
				break
			}
			err = unwrapper.Unwrap()
			continue
		}
		c := cause.Cause()
		if c == nil {
			break
		}
		err = c
	}
	return err
}

// withStack is an error that has a cause and a stack trace but no additional message.
type withStack struct {
	cause error
	*stack
}

func (w *withStack) Error() string { return w.cause.Error() }
func (w *withStack) Cause() error  { return w.cause }
func (w *withStack) Unwrap() error { return w.cause }

func (w *withStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v", w.Cause())
			w.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		_, _ = fmt.Fprintf(s, "%s", w.cause)
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", w.cause)
	}
}

// withMessage is an error that has a cause and an additional message but no stack trace.
type withMessage struct {
	cause error
	msg   string
}

func (w *withMessage) Error() string { return w.msg + ": " + w.cause.Error() }
func (w *withMessage) Cause() error  { return w.cause }
func (w *withMessage) Unwrap() error { return w.cause }

func (w *withMessage) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v\n", w.Cause())
			_, _ = fmt.Fprintf(s, "%s", w.msg)
			return
		}
		fallthrough
	case 's':
		_, _ = fmt.Fprintf(s, "%s", w.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", w.Error())
	}
}

// withCause is an error that has a cause, an additional message, and a stack trace.
// This is the result of Wrap/Wrapf.
type withCause struct {
	cause error
	msg   string
	*stack
}

func (w *withCause) Error() string { return w.msg + ": " + w.cause.Error() }
func (w *withCause) Cause() error  { return w.cause }
func (w *withCause) Unwrap() error { return w.cause }

func (w *withCause) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v\n", w.Cause())
			_, _ = fmt.Fprintf(s, "%s", w.msg)
			w.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		_, _ = fmt.Fprintf(s, "%s", w.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", w.Error())
	}
}
