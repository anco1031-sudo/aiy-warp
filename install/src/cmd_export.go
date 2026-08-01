package main

import (
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/anco1031-sudo/aiy-warp/install/src/internal/config"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/errs"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/kit"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/opencode"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/render"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/ui"
)

// webPlatforms are the single-agent hosts that use the condensed persona
// engine (CLI_SPEC.md §3 renderer matrix).
var webPlatforms = map[string]bool{"chatgpt": true, "gemini": true, "web": true}

func cmdExport(u *ui.UI, rest []string) int {
	platform, rest := extractPlatform("export", rest)
	fs := flag.NewFlagSet("export", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	agent := fs.String("agent", "", "export exactly one agent + its owned skills")
	team := fs.String("team", "", "export a department head + executors")
	out := fs.String("out", "", "output directory (default export/<platform> under CWD)")
	stdout := fs.Bool("stdout", false, "print the bundle to stdout")
	collapse := fs.Bool("collapse", false, "merge a team into one conductor prompt (web platforms)")
	noCollapse := fs.Bool("no-collapse", false, "keep one file per agent (opencode default)")
	cfgPath := fs.String("config", "", "path to warp.config")
	verbose := fs.Bool("verbose", false, "verbose output")
	allowIDs := fs.String("allow-identifiers", "", "comma-separated public IDs exempt from the export gate")
	if handled, code := parseCmdFlags(u, fs, rest); handled {
		return code
	}
	if *collapse && *noCollapse {
		fmt.Fprintln(u.Err, "aiy warp: --collapse and --no-collapse are mutually exclusive")
		return errs.CodeUsage
	}
	if platform == "" {
		fmt.Fprintln(u.Err, "usage: aiy warp export <platform> [flags]")
		return errs.CodeUsage
	}
	if !webPlatforms[platform] && platform != "opencode" {
		fmt.Fprintf(u.Err, "aiy warp: unknown platform %q (have: opencode, %s)\n", platform, strings.Join(render.Names(), ", "))
		return errs.CodeUsage
	}

	repoRoot, cfg, bundles, code := setup(u, *cfgPath, *verbose)
	if code != errs.CodeOK {
		return code
	}

	if platform == "opencode" {
		return runOpenCodeExport(u, repoRoot, cfg, bundles, *out, *stdout, *agent, *team, *allowIDs, *collapse, *noCollapse)
	}

	if *collapse && *agent != "" {
		fmt.Fprintln(u.Err, "aiy warp: --collapse applies to --team, not --agent (single agents are already condensed)")
		return errs.CodeUsage
	}
	outDir := *out
	if outDir == "" {
		outDir = filepath.Join(".", "export", platform)
	}
	res, err := render.RunExport(render.ExportOptions{
		RepoRoot:   repoRoot,
		Out:        outDir,
		Stdout:     *stdout,
		Platform:   platform,
		Bundles:    bundles,
		Selector:   kit.Selector{Agent: *agent, Team: *team},
		AllowIDs:   mergeIDs(cfg.AllowIDSet(), *allowIDs),
		NoCollapse: *noCollapse,
	})
	if err != nil {
		return reportErr(u, err)
	}
	if *stdout {
		u.Verbosef("exported %s persona (%d words) to stdout", platform, res.Words)
	} else {
		u.OK("Exported %s persona (%d words) to %s", platform, res.Words, outDir)
	}
	return errs.CodeOK
}

func runOpenCodeExport(u *ui.UI, repoRoot string, cfg *config.Config, bundles *kit.Bundles, out string, stdout bool, agent, team, allowIDs string, collapse, noCollapse bool) int {
	outDir := out
	if outDir == "" {
		outDir = filepath.Join(".", "export", "opencode")
	}
	files, err := opencode.RunExport(opencode.ExportOptions{
		RepoRoot: repoRoot,
		Out:      outDir,
		Stdout:   stdout,
		Bundles:  bundles,
		Selector: kit.Selector{Agent: agent, Team: team},
		AllowIDs: mergeIDs(cfg.AllowIDSet(), allowIDs),
		Collapse: collapse,
	})
	if err != nil {
		return reportErr(u, err)
	}
	if stdout {
		u.Verbosef("exported %d file(s) to stdout", len(files))
	} else {
		u.OK("Exported %d file(s) to %s", len(files), outDir)
	}
	_ = noCollapse
	return errs.CodeOK
}
