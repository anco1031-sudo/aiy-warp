// Package migrate stamps canonical warp-v1 frontmatter (CLI_SPEC.md §2.1) into
// legacy agent files. ADDITIVE ONLY: legacy host fields (mode, model, color,
// permission, description) are preserved verbatim and the persona body is never
// rewritten. Idempotent: files already carrying warp_version are skipped.
package migrate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/anco1031-sudo/aiy-warp/install/src/internal/condense"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/frontmatter"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/kit"
	"gopkg.in/yaml.v3"
)

// Options configures the migrate run.
type Options struct {
	AgentsDir string // repo agents/ dir
	Bundles   *kit.Bundles
	DryRun    bool
	Backup    bool // write {name}.md.bak before stamping (default true)
}

// Result summarizes the run.
type Result struct {
	Stamped []string // files stamped with canonical frontmatter
	Skipped []string // already canonical (idempotent no-op)
	Errors  []string
}

// Noop reports whether nothing was (or would be) stamped.
func (r *Result) Noop() bool { return len(r.Stamped) == 0 && len(r.Errors) == 0 }

// canonicalFM is the additive frontmatter shape (canonical keys first, then
// preserved legacy host keys).
type canonicalFM struct {
	WarpVersion     int      `yaml:"warp_version"`
	Name            string   `yaml:"name"`
	DisplayName     string   `yaml:"display_name"`
	Role            string   `yaml:"role"`
	Color           string   `yaml:"color,omitempty"`
	Department      string   `yaml:"department"`
	Rank            string   `yaml:"rank"`
	ReportsTo       string   `yaml:"reports_to,omitempty"`
	ModelHint       string   `yaml:"model_hint,omitempty"`
	Description     string   `yaml:"description"`
	Mode            string   `yaml:"mode,omitempty"`
	Model           string   `yaml:"model,omitempty"`
	Permission      any      `yaml:"permission,omitempty"`
	Personality     string   `yaml:"personality,omitempty"`
	Directives      []string `yaml:"directives,omitempty"`
	Boundaries      []string `yaml:"boundaries,omitempty"`
	PlatformTargets []string `yaml:"platform_targets"`
	Skills          []string `yaml:"skills,omitempty"`
}

var descRE = regexp.MustCompile(`^([A-Z][A-Za-z]+)\s+\(([^)]+)\)\s*[—\-–]\s*(.+)$`)

// Run stamps canonical frontmatter into every legacy agent file.
func Run(o Options) (*Result, error) {
	paths, err := filepath.Glob(filepath.Join(o.AgentsDir, "*.md"))
	if err != nil {
		return nil, err
	}
	sort.Strings(paths)
	res := &Result{}
	for _, p := range paths {
		name := strings.TrimSuffix(filepath.Base(p), ".md")
		if o.Bundles.Find(name) == nil {
			res.Errors = append(res.Errors, name+": not in bundles.yaml — skipped")
			continue
		}
		content, err := os.ReadFile(p)
		if err != nil {
			res.Errors = append(res.Errors, name+": "+err.Error())
			continue
		}
		a, err := frontmatter.ParseAgent(filepath.Base(p), string(content))
		if err != nil {
			res.Errors = append(res.Errors, name+": parse: "+err.Error())
			continue
		}
		if a.Canonical {
			res.Skipped = append(res.Skipped, name)
			continue // idempotent — already stamped
		}
		org := o.Bundles.Find(name)
		next, err := stamp(a, org, string(content))
		if err != nil {
			res.Errors = append(res.Errors, name+": "+err.Error())
			continue
		}
		res.Stamped = append(res.Stamped, name)
		if o.DryRun {
			continue
		}
		if o.Backup {
			if err := os.WriteFile(p+".bak", content, 0o644); err != nil {
				res.Errors = append(res.Errors, name+": backup: "+err.Error())
				continue
			}
		}
		if err := os.WriteFile(p, []byte(next), 0o644); err != nil {
			res.Errors = append(res.Errors, name+": write: "+err.Error())
			continue
		}
	}
	return res, nil
}

// stamp builds the canonical content for one legacy agent file. The body is
// preserved byte-for-byte; only the frontmatter block is replaced.
func stamp(a *frontmatter.Agent, org *kit.Agent, content string) (string, error) {
	fm, body, err := frontmatter.SplitFrontmatter(content)
	if err != nil {
		return "", err
	}
	if fm == "" {
		return "", fmt.Errorf("no frontmatter to extend")
	}
	// Preserve the legacy permission block as-is (host tuning).
	var legacy map[string]any
	if err := yaml.Unmarshal([]byte(fm), &legacy); err != nil {
		return "", fmt.Errorf("legacy frontmatter: %w", err)
	}
	display, role := deriveIdentity(a.Name, a.Description)
	reports := ""
	if org.ReportsTo != nil {
		reports = *org.ReportsTo
	}
	per := condense.Extract(a, false)
	next := canonicalFM{
		WarpVersion:     1,
		Name:            a.Name,
		DisplayName:     display,
		Role:            role,
		Color:           a.Color,
		Department:      org.Department,
		Rank:            org.Rank,
		ReportsTo:       reports,
		ModelHint:       a.ModelHint,
		Description:     a.Description,
		Mode:            a.Mode,
		Model:           a.ModelHint,
		Permission:      legacy["permission"],
		Personality:     per.Personality,
		Directives:      per.Directives,
		Boundaries:      per.Boundaries,
		PlatformTargets: []string{"opencode", "chatgpt", "gemini", "web", "claude"},
		Skills:          org.Skills,
	}
	out, err := yaml.Marshal(&next)
	if err != nil {
		return "", fmt.Errorf("marshal canonical frontmatter: %w", err)
	}
	return "---\n" + string(out) + "---\n" + body, nil
}

// deriveIdentity extracts display_name and role from a legacy description of
// the form "NAME (ไทย) — Role." → display "ไทย (Name)", role "Role".
func deriveIdentity(name, description string) (display, role string) {
	display = name
	if m := descRE.FindStringSubmatch(description); m != nil {
		display = m[2] + " (" + titleCase(m[1]) + ")"
		role = strings.TrimSpace(m[3])
		if i := strings.Index(role, "."); i > 0 {
			role = role[:i]
		}
	}
	if role == "" {
		role = name
	}
	return display, role
}

func titleCase(s string) string {
	if s == "" {
		return s
	}
	r := []rune(strings.ToLower(s))
	r[0] = []rune(strings.ToUpper(string(r[0])))[0]
	return string(r)
}
