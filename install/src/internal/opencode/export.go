package opencode

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/anco1031-sudo/aiy-warp/install/src/internal/errs"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/kit"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/redact"
)

// ExportOptions configures `aiy warp export opencode`.
type ExportOptions struct {
	RepoRoot string
	Out      string // output directory (ignored when Stdout)
	Stdout   bool
	Bundles  *kit.Bundles
	Selector kit.Selector
	AllowIDs map[string]bool
	Collapse bool // must stay false for opencode in P0 (flat one-file-per-agent)
}

// Run exports a flat file set (no collapse) to Out or stdout. Returns the
// exported repo-relative paths. Returns errs.Secret (5) when the redaction
// gate trips (credential values, or identifiers without --allow-identifiers).
func RunExport(o ExportOptions) ([]string, error) {
	if o.Collapse {
		return nil, errs.Usage("--collapse is a web-chat renderer feature (P1); opencode export is flat one-file-per-agent")
	}
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
			"credential value detected — refusing to export:\n  %s",
			strings.Join(rep.BlockFiles(), "\n  ")))
	}
	if len(rep.Exports) > 0 {
		ids := map[string]bool{}
		for _, f := range rep.Exports {
			ids[f.Match] = true
		}
		keys := make([]string, 0, len(ids))
		for id := range ids {
			keys = append(keys, id)
		}
		sort.Strings(keys)
		return nil, errs.Secret(fmt.Sprintf(
			"sensitive identifiers detected — refusing to export:\n  %s\nDeclare the public IDs with --allow-identifiers %s (CLI_SPEC.md §5.4).",
			strings.Join(rep.ExportBlockFiles(), "\n  "), strings.Join(keys, ",")))
	}

	if o.Stdout {
		var sb strings.Builder
		for _, rel := range files {
			sb.WriteString("#### " + rel + "\n\n")
			sb.WriteString(contents[rel])
			if !strings.HasSuffix(contents[rel], "\n") {
				sb.WriteString("\n")
			}
			sb.WriteString("\n")
		}
		fmt.Print(sb.String())
		return files, nil
	}

	for _, rel := range files {
		dest := filepath.Join(o.Out, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return nil, errs.Wrapf(err, "mkdir %s", filepath.Dir(dest))
		}
		if err := os.WriteFile(dest, []byte(contents[rel]), 0o644); err != nil {
			return nil, errs.Wrapf(err, "write %s", dest)
		}
	}
	return files, nil
}
