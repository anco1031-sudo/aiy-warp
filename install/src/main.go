// Command `aiy warp` — warp the 18-agent kit to a host (P0: opencode).
//
// Exit codes (CLI_SPEC.md §1.3): 0 success · 1 runtime · 2 usage · 3 drift ·
// 4 no-op · 5 secret detected.
package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/anco1031-sudo/aiy-warp/install/src/internal/errs"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/ui"
)

const version = "0.1.0-p0"

func main() { os.Exit(run(os.Args[1:])) }

func run(argv []string) int {
	u := &ui.UI{Out: os.Stdout, Err: os.Stderr}

	// Accept `aiy warp <cmd>`, `aiy-warp <cmd>`, or `warp <cmd>`.
	if len(argv) > 0 && argv[0] == "warp" {
		argv = argv[1:]
	}
	if len(argv) == 0 {
		usage(u.Out)
		return errs.CodeUsage
	}
	if argv[0] == "--version" {
		fmt.Fprintf(u.Out, "aiy warp %s\n", version)
		return errs.CodeOK
	}
	cmd, rest := argv[0], argv[1:]
	switch cmd {
	case "install":
		return cmdInstall(u, rest)
	case "export":
		return cmdExport(u, rest)
	case "doctor":
		return cmdDoctor(u, rest)
	case "list":
		return cmdList(u, rest)
	case "help", "--help", "-h":
		usage(u.Out)
		return errs.CodeOK
	default:
		fmt.Fprintf(u.Err, "aiy warp: unknown command %q\n", cmd)
		usage(u.Err)
		return errs.CodeUsage
	}
}

func reportErr(u *ui.UI, err error) int {
	code := errs.CodeOf(err)
	fmt.Fprintf(u.Err, "aiy warp: %v\n", err)
	switch code {
	case errs.CodeDrift:
		fmt.Fprintln(u.Err, "Hint: run 'aiy warp doctor opencode' for detail.")
	case errs.CodeSecret:
		fmt.Fprintln(u.Err, "Hint: the kit carries identities + workflows, never credentials (WARP.md).")
	}
	return code
}

func mergeIDs(base map[string]bool, flagVal string) map[string]bool {
	out := make(map[string]bool, len(base))
	for k := range base {
		out[k] = true
	}
	for _, p := range strings.Split(flagVal, ",") {
		if t := strings.TrimSpace(p); t != "" {
			out[t] = true
		}
	}
	return out
}

func pick(s, def string) string {
	if s != "" {
		return s
	}
	return def
}

func usage(w io.Writer) {
	fmt.Fprintln(w, `aiy warp — warp the 18-agent kit to a host (P0: opencode)

Usage:
  aiy warp install <platform> [--agent <name>|--team <dept>] [--force] [--dry-run] [-y]
  aiy warp export  <platform> [--agent <name>|--team <dept>] [--out <path>|--stdout]
                   [--allow-identifiers <id,id>] [--collapse|--no-collapse]
  aiy warp doctor  [--platform <p>] [--json]
  aiy warp list    [--platform <p>]
  aiy warp --version

Global flags: --config <path> · --verbose
Exit codes: 0 success · 1 runtime · 2 usage · 3 drift · 4 no-op · 5 secret detected
Spec: install/CLI_SPEC.md · Design: install/DESIGN-NOTE.md`)
}
