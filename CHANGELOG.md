# Changelog

## v0.1.0

Initial release — drop-in replacement for `github.com/pkg/errors`.

### Features

- `New` and `Errorf` for creating errors with stack traces
- `Wrap` and `Wrapf` for wrapping errors with context and a new stack trace
- `WithStack` for adding a stack trace to an existing error
- `WithMessage` and `WithMessagef` for adding context without a new stack trace
- `Cause` for walking the error chain to the root cause
- `Frame` and `StackTrace` types for stack trace inspection
- `%+v` formatting matching `pkg/errors` output format
- Dual `Cause()`/`Unwrap()` implementation on all wrapper types — `errors.Is` and `errors.As` work through the entire chain
- Zero dependencies (stdlib only)
- Go 1.21+ support
- GitHub Actions CI with Go 1.21, 1.22, 1.23, and stable
- Comprehensive test suite with unit, compatibility, benchmark, and example tests
