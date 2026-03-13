// Package faultline provides error handling primitives with stack traces.
//
// faultline is a drop-in replacement for github.com/pkg/errors that fixes
// the compatibility issue between Cause() and errors.Is/errors.As by
// implementing both Cause() and Unwrap() on all wrapper types.
//
// The formatting verbs are implemented directly on each error type
// (fundamental, withStack, withMessage, withCause) in their respective
// source files (errors.go, wrap.go).
package faultline
