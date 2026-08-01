package main

import (
	"flag"
	"fmt"
	"io"

	"github.com/anco1031-sudo/aiy-warp/install/src/internal/errs"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/kit"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/opencode"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/ui"
)

func cmdInstall(u *ui.UI, rest []string) int {
	platform, rest := extractPlatform("install", rest)
	fs := flag.NewFlagSet("install", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	agent := fs.String("agent", "", "warp exactly one agent (e.g. aiy) + its owned skills")
	team := fs.String("team", "", "warp a department head + executors (e.g. kwan)")
	force := fs.Bool("force", false, "overwrite drifted local files")
	dryRun := fs.Bool("dry-run", false, "print actions without executing")
	yes := fs.Bool("y", false, "skip prompts (accepted)")
	cfgPath := fs.String("config", "", "path to warp.config (default ~/.config/aiy-warp/warp.config)")
	verbose := fs.Bool("verbose", false, "verbose output")
	allowIDs := fs.String("allow-identifiers", "", "comma-separated public IDs exempt from the export gate")
	if err := fs.Parse(rest); err != nil {
		fmt.Fprintln(u.Err, "aiy warp install:", err)
		return errs.CodeUsage
	}
	_ = yes
	if platform == "" {
		fmt.Fprintln(u.Err, "usage: aiy warp install <platform> [flags]")
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

	res, err := opencode.RunInstall(opencode.InstallOptions{
		RepoRoot: repoRoot,
		DestDir:  cfg.OpencodeDest(),
		StateDir: cfg.StateDirResolved(),
		Bundles:  bundles,
		Selector: kit.Selector{Agent: *agent, Team: *team},
		Force:    *force,
		DryRun:   *dryRun,
		AllowIDs: mergeIDs(cfg.AllowIDSet(), *allowIDs),
	})
	if err != nil {
		return reportErr(u, err)
	}

	for _, w := range res.Warns {
		u.Verbosef("%s: %s (%s)", w.File, w.Match, w.Pattern)
	}
	if len(res.Warns) > 0 {
		u.Verbosef("note: %d warning(s) — credential references/host paths (P1 parameterization)", len(res.Warns))
	}

	if *dryRun {
		u.Info("→ Would install %d file(s) (new), update %d (--force), skip %d drifted",
			len(res.Installed), len(res.Updated), len(res.Skipped))
		for _, f := range res.Installed {
			u.Verbosef("  + %s", f)
		}
		for _, f := range res.Skipped {
			u.Warn("  ~ %s (drifted, needs --force)", f)
		}
		return errs.CodeOK
	}

	if res.Noop {
		u.Info("Already in sync — nothing to do.")
		return errs.CodeNoop
	}
	if len(res.Skipped) > 0 {
		u.Warn("%d file(s) skipped (drifted) — run with --force to override, or 'aiy warp doctor opencode'", len(res.Skipped))
	}
	u.OK("Installed %d file(s), updated %d, %d in sync → %s",
		len(res.Installed), len(res.Updated), res.Same, cfg.OpencodeDest())
	if len(res.Installed) == 0 && len(res.Updated) == 0 && len(res.Skipped) > 0 {
		u.Warn("Install made no progress — all files drifted. See: aiy warp doctor opencode")
		return errs.CodeDrift
	}
	return errs.CodeOK
}
