# faultline

[![CI](https://github.com/agentine/faultline/actions/workflows/ci.yml/badge.svg)](https://github.com/agentine/faultline/actions/workflows/ci.yml)

A modern, lightweight drop-in replacement for [github.com/pkg/errors](https://github.com/pkg/errors).

## Why faultline?

`pkg/errors` is one of the most widely imported Go packages (339K+ importers, 8.2K stars) but has been **archived since 2021**. It has two critical issues:

1. **No maintenance** — 42 open issues, no security fixes, no bug fixes
2. **`Cause()` is invisible to `errors.Is`/`errors.As`** — wrapped errors don't implement `Unwrap()`, breaking Go 1.13+ error inspection

**faultline** fixes both by implementing `Cause()` *and* `Unwrap()` on all wrapper types. Zero dependencies, identical API.

## Installation

```
go get github.com/agentine/faultline
```

Requires Go 1.21+.

## Quick Start

```go
import "github.com/agentine/faultline"

// Create errors with stack traces
err := faultline.New("something failed")

// Wrap with context
wrapped := faultline.Wrap(err, "operation failed")

// Stack traces with %+v
fmt.Printf("%+v\n", wrapped)

// errors.Is / errors.As work through the entire chain
if errors.Is(wrapped, err) {
    // works! (broken in pkg/errors)
}

// Cause() still works for backwards compatibility
root := faultline.Cause(wrapped)
```

## API Reference

### Error Creation

| Function | Description |
|----------|-------------|
| `New(message string) error` | Create a new error with a stack trace |
| `Errorf(format string, args ...interface{}) error` | Create a formatted error with a stack trace |

### Wrapping

| Function | Description |
|----------|-------------|
| `Wrap(err error, message string) error` | Wrap with message + new stack trace (nil-safe) |
| `Wrapf(err error, format string, args ...interface{}) error` | Wrap with formatted message + new stack trace (nil-safe) |
| `WithStack(err error) error` | Add a stack trace without a message (nil-safe) |
| `WithMessage(err error, message string) error` | Add context message without a new stack trace (nil-safe) |
| `WithMessagef(err error, format string, args ...interface{}) error` | Add formatted context message without a new stack trace (nil-safe) |

### Inspection

| Function | Description |
|----------|-------------|
| `Cause(err error) error` | Walk the error chain to the root cause (supports both `Cause()` and `Unwrap()` chains) |

### Types

| Type | Description |
|------|-------------|
| `Frame` | A program counter representing a stack frame |
| `StackTrace` | A slice of `Frame` values from innermost to outermost |

## Accessing Stack Traces

Errors created by faultline expose a `StackTrace()` method via the `stackTracer` interface. This is useful for logging integrations (logrus, zerolog, slog, etc.) that want to extract structured frame data.

```go
type stackTracer interface {
    StackTrace() faultline.StackTrace
}

err := faultline.New("something failed")
if st, ok := err.(stackTracer); ok {
    for _, f := range st.StackTrace() {
        fmt.Printf("%+v\n", f)
    }
}
```

Each `faultline.Frame` supports `%s` (file), `%d` (line), `%n` (function name), and `%+s` / `%+v` for the full function+file format.

## Migration from pkg/errors

One-liner import replacement:

```
gofmt -r '"github.com/pkg/errors" -> "github.com/agentine/faultline"' -w .
```

No code changes needed — the API is identical.

## `%+v` Output Format

Stack trace output matches `pkg/errors` format exactly, so existing log parsers and tools continue to work:

```
operation failed
github.com/agentine/faultline_test.TestFormatWrap
	/path/to/faultline/format_test.go:42
testing.tRunner
	/usr/local/go/src/testing/testing.go:1595
```

## Benchmarks

```
BenchmarkNew-8             4022078     274.5 ns/op    280 B/op    2 allocs/op
BenchmarkErrorf-8          3131156     378.2 ns/op    328 B/op    5 allocs/op
BenchmarkWrap-8            4010323     282.0 ns/op    280 B/op    2 allocs/op
BenchmarkWrapf-8           3198045     379.4 ns/op    352 B/op    5 allocs/op
BenchmarkWithStack-8       4270057     280.5 ns/op    280 B/op    2 allocs/op
BenchmarkWithMessage-8     1000000000  0.32 ns/op     0 B/op      0 allocs/op
BenchmarkCause-8           100000000   10.74 ns/op    0 B/op      0 allocs/op
BenchmarkFormatPlusV-8     429022      2547 ns/op     912 B/op    17 allocs/op
```

## License

BSD-2-Clause (same as pkg/errors)
