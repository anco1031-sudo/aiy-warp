// Package errs defines typed errors that map to the CLI exit codes of
// CLI_SPEC.md §1.3. Plain errors are treated as runtime (exit 1).
package errs

import "fmt"

// Exit codes (CLI_SPEC.md §1.3).
const (
	CodeOK      = 0 // success / no drift
	CodeRuntime = 1 // io, parse, network
	CodeUsage   = 2 // bad flag, unknown platform/agent/team
	CodeDrift   = 3 // drift detected
	CodeNoop    = 4 // nothing to do (no-op)
	CodeSecret  = 5 // credential detected — refusing
)

// Error is a typed failure carrying the process exit code.
type Error struct {
	Code int
	Msg  string
}

func (e *Error) Error() string { return e.Msg }

// Runtime wraps a runtime failure (exit 1).
func Runtime(msg string) *Error { return &Error{Code: CodeRuntime, Msg: msg} }

// Usage wraps a usage error (exit 2).
func Usage(msg string) *Error { return &Error{Code: CodeUsage, Msg: msg} }

// Drift wraps a drift condition (exit 3).
func Drift(msg string) *Error { return &Error{Code: CodeDrift, Msg: msg} }

// Noop wraps a nothing-to-do condition (exit 4).
func Noop(msg string) *Error { return &Error{Code: CodeNoop, Msg: msg} }

// Secret wraps a redaction-gate rejection (exit 5).
func Secret(msg string) *Error { return &Error{Code: CodeSecret, Msg: msg} }

// Wrapf annotates an underlying error. A wrapped *Error keeps its exit code;
// any other error becomes a runtime error.
func Wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	a := append([]any{}, args...)
	if e, ok := err.(*Error); ok {
		return &Error{Code: e.Code, Msg: fmt.Sprintf(format+": %v", append(a, e.Msg)...)}
	}
	return &Error{Code: CodeRuntime, Msg: fmt.Sprintf(format+": %v", append(a, err)...)}
}

// CodeOf extracts the exit code for any error (nil → 0).
func CodeOf(err error) int {
	if err == nil {
		return CodeOK
	}
	if e, ok := err.(*Error); ok {
		return e.Code
	}
	return CodeRuntime
}
