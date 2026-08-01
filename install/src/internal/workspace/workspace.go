// Package workspace implements `aiy warp init-workspace` — creating the PARA
// home structure (00-Inbox … 05-Other) at $HOME from scaffold/PARA-template/.
// Idempotent: a fully-present structure is a no-op (exit 4).
package workspace

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Options configures the init-workspace run.
type Options struct {
	Home        string // destination root (default $HOME)
	TemplateDir string // scaffold/PARA-template/
	DryRun      bool
}

// Result summarizes what was created.
type Result struct {
	Dirs     []string // directories created
	Copied   []string // README.md files copied
	Existing []string // already present (no-op components)
	Noop     bool     // everything already present
}

// Run creates the PARA structure. Returns (result, nil); a missing template
// directory is a runtime error.
func Run(o Options) (*Result, error) {
	entries, err := os.ReadDir(o.TemplateDir)
	if err != nil {
		return nil, err
	}
	var dirs []string
	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, e.Name())
		}
	}
	sort.Strings(dirs)
	if len(dirs) == 0 {
		return nil, os.ErrNotExist // no template dirs
	}
	res := &Result{}
	for _, d := range dirs {
		dest := filepath.Join(o.Home, d)
		if _, err := os.Stat(dest); os.IsNotExist(err) {
			res.Dirs = append(res.Dirs, d)
			if !o.DryRun {
				if err := os.MkdirAll(dest, 0o755); err != nil {
					return nil, err
				}
			}
		} else {
			res.Existing = append(res.Existing, d)
		}
		readme := filepath.Join(dest, "README.md")
		if _, err := os.Stat(readme); os.IsNotExist(err) {
			res.Copied = append(res.Copied, filepath.Join(d, "README.md"))
			if !o.DryRun {
				src, err := os.ReadFile(filepath.Join(o.TemplateDir, d, "README.md"))
				if err != nil {
					return nil, err
				}
				if err := os.WriteFile(readme, src, 0o644); err != nil {
					return nil, err
				}
			}
		} else {
			res.Existing = append(res.Existing, filepath.Join(d, "README.md"))
		}
	}
	res.Noop = len(res.Dirs) == 0 && len(res.Copied) == 0
	return res, nil
}

// Summary renders a human line like "00-Inbox, 01-Projects, …".
func Summary(names []string) string {
	if len(names) == 0 {
		return "—"
	}
	return strings.Join(names, ", ")
}
