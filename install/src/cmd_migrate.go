package main

import (
	"flag"
	"io"
	"path/filepath"
	"strings"

	"github.com/anco1031-sudo/aiy-warp/install/src/internal/errs"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/migrate"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/ui"
)

func cmdMigrate(u *ui.UI, rest []string) int {
	fs := flag.NewFlagSet("migrate", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	dryRun := fs.Bool("dry-run", false, "report what would be stamped without writing")
	noBackup := fs.Bool("no-backup", false, "skip writing {name}.md.bak originals")
	agentsDir := fs.String("agents-dir", "", "agents directory (default <repo>/agents)")
	cfgPath := fs.String("config", "", "path to warp.config")
	verbose := fs.Bool("verbose", false, "verbose output")
	if handled, code := parseCmdFlags(u, fs, rest); handled {
		return code
	}

	repoRoot, _, bundles, code := setup(u, *cfgPath, *verbose)
	if code != errs.CodeOK {
		return code
	}
	dir := *agentsDir
	if dir == "" {
		dir = filepath.Join(repoRoot, "agents")
	}

	res, err := migrate.Run(migrate.Options{
		AgentsDir: dir,
		Bundles:   bundles,
		DryRun:    *dryRun,
		Backup:    !*noBackup,
	})
	if err != nil {
		return reportErr(u, errs.Wrapf(err, "migrate"))
	}
	for _, e := range res.Errors {
		u.Warn("%s", e)
	}
	if *dryRun {
		u.Info("→ Would stamp %d file(s): %s", len(res.Stamped), strings.Join(res.Stamped, ", "))
		if len(res.Skipped) > 0 {
			u.Verbosef("already canonical (skipped): %s", strings.Join(res.Skipped, ", "))
		}
		if res.Noop() {
			u.Info("Nothing to stamp — all files already canonical.")
			return errs.CodeNoop
		}
		return errs.CodeOK
	}
	if res.Noop() {
		u.Info("Nothing to stamp — all agent files already canonical (idempotent).")
		return errs.CodeNoop
	}
	u.OK("Stamped canonical frontmatter into %d file(s)", len(res.Stamped))
	for _, n := range res.Stamped {
		u.Verbosef("  + %s.md", n)
	}
	if len(res.Skipped) > 0 {
		u.Verbosef("already canonical: %s", strings.Join(res.Skipped, ", "))
	}
	if !*noBackup {
		u.Info("Originals backed up as agents/*.md.bak (gitignored).")
	}
	return errs.CodeOK
}
