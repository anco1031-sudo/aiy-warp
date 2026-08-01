package main

import (
	"flag"
	"io"
	"os"
	"path/filepath"

	"github.com/anco1031-sudo/aiy-warp/install/src/internal/errs"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/ui"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/workspace"
)

func cmdInit(u *ui.UI, rest []string) int {
	fs := flag.NewFlagSet("init-workspace", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	home := fs.String("home", "", "destination root (default $HOME)")
	dryRun := fs.Bool("dry-run", false, "print actions without executing")
	cfgPath := fs.String("config", "", "path to warp.config")
	verbose := fs.Bool("verbose", false, "verbose output")
	if handled, code := parseCmdFlags(u, fs, rest); handled {
		return code
	}

	repoRoot, cfg, _, code := setup(u, *cfgPath, *verbose)
	if code != errs.CodeOK {
		return code
	}
	dest := *home
	if dest == "" {
		dest = cfg.Paths.Home
	}
	if dest == "" {
		var err error
		dest, err = os.UserHomeDir()
		if err != nil {
			return reportErr(u, errs.Wrapf(err, "resolve $HOME"))
		}
	}
	tmpl := filepath.Join(repoRoot, "scaffold", "PARA-template")

	res, err := workspace.Run(workspace.Options{Home: dest, TemplateDir: tmpl, DryRun: *dryRun})
	if err != nil {
		return reportErr(u, errs.Wrapf(err, "init-workspace"))
	}
	if *dryRun {
		u.Info("→ Would create %d dir(s): %s", len(res.Dirs), workspace.Summary(res.Dirs))
		u.Info("→ Would copy %d README(s): %s", len(res.Copied), workspace.Summary(res.Copied))
		if res.Noop {
			u.Info("PARA structure already present — nothing to do.")
			return errs.CodeNoop
		}
		return errs.CodeOK
	}
	if res.Noop {
		u.Info("PARA structure already present at %s — nothing to do.", dest)
		return errs.CodeNoop
	}
	u.OK("PARA workspace initialized at %s (%d dirs, %d READMEs)", dest, len(res.Dirs), len(res.Copied))
	for _, d := range res.Dirs {
		u.Verbosef("  + %s/", d)
	}
	for _, c := range res.Copied {
		u.Verbosef("  + %s", c)
	}
	return errs.CodeOK
}
