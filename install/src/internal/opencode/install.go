// Package opencode implements the install and export renderers for the
// opencode platform (CLI_SPEC.md §3, §6).
package opencode

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/anco1031-sudo/aiy-warp/install/src/internal/errs"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/kit"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/manifest"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/redact"
)

// InstallOptions configures `aiy warp install opencode`.
type InstallOptions struct {
	RepoRoot string
	DestDir  string // e.g. ~/.config/opencode
	StateDir string // where warp.lock lives
	Bundles  *kit.Bundles
	Selector kit.Selector
	Force    bool
	DryRun   bool
	AllowIDs map[string]bool
}

// InstallResult summarizes what happened.
type InstallResult struct {
	Installed []string // newly written
	Updated   []string // force-rewritten (merged for agents)
	Skipped   []string // drifted, not forced
	Same      int
	Noop      bool // nothing to do (already in sync)
	Warns     []redact.Finding
	Notes     []string // informational (e.g. symlink replaced)
}

// Run executes the install. Returns errs.Secret (5) on credential values and
// errs.Noop (4) when everything is already in sync.
func RunInstall(o InstallOptions) (*InstallResult, error) {
	k, err := kit.Discover(o.RepoRoot)
	if err != nil {
		return nil, errs.Wrapf(err, "discover kit")
	}
	files, err := k.Resolve(o.Bundles, o.Selector)
	if err != nil {
		return nil, errs.Usage(err.Error())
	}

	contents := map[string]string{}
	for _, rel := range files {
		c, err := os.ReadFile(filepath.Join(o.RepoRoot, filepath.FromSlash(rel)))
		if err != nil {
			return nil, errs.Wrapf(err, "read %s", rel)
		}
		contents[rel] = string(c)
	}

	rep := redact.Scan(contents, o.AllowIDs)
	if rep.HasBlocks() {
		return nil, errs.Secret(fmt.Sprintf(
			"credential value detected — refusing to install:\n  %s\nRun 'aiy warp doctor opencode' for detail.",
			strings.Join(rep.BlockFiles(), "\n  ")))
	}
	// Identifier findings (export gate only) do not block install: they are
	// host-routing metadata; P1 parameterization removes them entirely.

	repoHashes := manifest.Hashes(contents)

	hostHashes := map[string]string{}
	for rel := range repoHashes {
		if c, err := os.ReadFile(destPath(o.DestDir, rel)); err == nil {
			hostHashes[rel] = manifest.KitHash(string(c), rel)
		}
	}

	entries := manifest.Diff(repoHashes, hostHashes)
	res := &InstallResult{Warns: rep.Warns}
	type writeOp struct{ rel, content string }
	var toWrite []writeOp

	for _, e := range entries {
		switch e.State {
		case manifest.StateSame:
			res.Same++
		case manifest.StateNew:
			res.Installed = append(res.Installed, e.Rel)
			toWrite = append(toWrite, writeOp{e.Rel, contents[e.Rel]})
		case manifest.StateDrifted:
			if !o.Force {
				res.Skipped = append(res.Skipped, e.Rel)
				continue
			}
			merged := contents[e.Rel]
			if strings.HasPrefix(e.Rel, "agents/") {
				if host, err := os.ReadFile(destPath(o.DestDir, e.Rel)); err == nil {
					if m, warns, err := manifest.MergeAgentContentV(contents[e.Rel], string(host)); err == nil {
						merged = m
						res.Notes = append(res.Notes, warns...)
					}
				}
			}
			res.Updated = append(res.Updated, e.Rel)
			toWrite = append(toWrite, writeOp{e.Rel, merged})
		}
	}

	res.Noop = len(res.Installed) == 0 && len(res.Updated) == 0 && len(res.Skipped) == 0
	if o.DryRun {
		return res, nil
	}

	for _, w := range toWrite {
		dest := destPath(o.DestDir, w.rel)
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return nil, errs.Wrapf(err, "mkdir %s", filepath.Dir(dest))
		}
		if fi, err := os.Lstat(dest); err == nil && fi.Mode()&os.ModeSymlink != 0 {
			res.Notes = append(res.Notes, "replaced symlink at "+w.rel)
			if err := os.Remove(dest); err != nil {
				return nil, errs.Wrapf(err, "remove symlink %s", dest)
			}
		}
		if err := os.WriteFile(dest, []byte(w.content), 0o644); err != nil {
			return nil, errs.Wrapf(err, "write %s", dest)
		}
	}

	lock := manifest.New("opencode", k.Rev, repoHashes)
	if err := manifest.Save(filepath.Join(o.StateDir, "warp.lock"), lock); err != nil {
		return nil, errs.Wrapf(err, "write warp.lock")
	}

	return res, nil
}

// destPath maps a repo-relative path to the opencode destination layout:
// agents/ → <dest>/agents/, skills/ → <dest>/skills/, templates/ and
// playbooks/ likewise.
func destPath(destDir, rel string) string {
	return filepath.Join(destDir, filepath.FromSlash(rel))
}
