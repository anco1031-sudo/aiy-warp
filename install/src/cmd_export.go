package main

import (
	"flag"
	"fmt"
	"io"
	"path/filepath"

	"github.com/anco1031-sudo/aiy-warp/install/src/internal/errs"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/kit"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/opencode"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/ui"
)

func cmdExport(u *ui.UI, rest []string) int {
	platform, rest := extractPlatform("export", rest)
	fs := flag.NewFlagSet("export", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	agent := fs.String("agent", "", "export exactly one agent + its owned skills")
	team := fs.String("team", "", "export a department head + executors")
	out := fs.String("out", "", "output directory (default export/opencode under CWD)")
	stdout := fs.Bool("stdout", false, "print the bundle to stdout")
	collapse := fs.Bool("collapse", false, "merge a team into one conductor prompt (P1, web chats)")
	noCollapse := fs.Bool("no-collapse", false, "keep one file per agent (default; opencode P0)")
	cfgPath := fs.String("config", "", "path to warp.config")
	verbose := fs.Bool("verbose", false, "verbose output")
	allowIDs := fs.String("allow-identifiers", "", "comma-separated public IDs exempt from the export gate")
	if err := fs.Parse(rest); err != nil {
		fmt.Fprintln(u.Err, "aiy warp export:", err)
		return errs.CodeUsage
	}
	if *collapse && *noCollapse {
		fmt.Fprintln(u.Err, "aiy warp: --collapse and --no-collapse are mutually exclusive")
		return errs.CodeUsage
	}
	if platform == "" {
		fmt.Fprintln(u.Err, "usage: aiy warp export <platform> [flags]")
		return errs.CodeUsage
	}
	if platform != "opencode" {
		fmt.Fprintf(u.Err, "aiy warp: unknown platform %q (P0 supports: opencode)\n", platform)
		return errs.CodeUsage
	}

	repoRoot, cfg, bundles, code := setup(u, *cfgPath, *verbose)
	if code != errs.CodeOK {
		return code
	}

	outDir := *out
	if outDir == "" {
		outDir = filepath.Join(".", "export", "opencode")
	}
	files, err := opencode.RunExport(opencode.ExportOptions{
		RepoRoot: repoRoot,
		Out:      outDir,
		Stdout:   *stdout,
		Bundles:  bundles,
		Selector: kit.Selector{Agent: *agent, Team: *team},
		AllowIDs: mergeIDs(cfg.AllowIDSet(), *allowIDs),
		Collapse: *collapse,
	})
	if err != nil {
		return reportErr(u, err)
	}
	if *stdout {
		u.Verbosef("exported %d file(s) to stdout", len(files))
	} else {
		u.OK("Exported %d file(s) to %s", len(files), outDir)
	}
	return errs.CodeOK
}
