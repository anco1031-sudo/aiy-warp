// Package ui renders human summaries and doctor JSON output.
package ui

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/anco1031-sudo/aiy-warp/install/src/internal/doctor"
)

// UI wraps output streams.
type UI struct {
	Out, Err io.Writer
	Verbose  bool
}

// OK prints a success line.
func (u *UI) OK(format string, a ...any) { fmt.Fprintf(u.Out, "✓ "+format+"\n", a...) }

// Warn prints a warning line.
func (u *UI) Warn(format string, a ...any) { fmt.Fprintf(u.Out, "⚠ "+format+"\n", a...) }

// Info prints a plain line.
func (u *UI) Info(format string, a ...any) { fmt.Fprintf(u.Out, format+"\n", a...) }

// Verbosef prints a line only in verbose mode.
func (u *UI) Verbosef(format string, a ...any) {
	if u.Verbose {
		fmt.Fprintf(u.Out, "· "+format+"\n", a...)
	}
}

// Doctor renders the doctor report as a human table.
func (u *UI) Doctor(r *doctor.Result) {
	for _, c := range r.Checks {
		mark := "PASS"
		if c.Status == doctor.Fail {
			mark = "FAIL"
		}
		fmt.Fprintf(u.Out, "[%s] %d. %s\n", mark, c.ID, c.Name)
		for _, h := range c.Hints {
			fmt.Fprintf(u.Out, "     %s\n", h)
		}
		for _, w := range c.Warns {
			fmt.Fprintf(u.Out, "     · %s\n", w)
		}
	}
}

// JSON writes v as indented JSON.
func (u *UI) JSON(v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(u.Out, string(b))
	return err
}
