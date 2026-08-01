package doctor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/anco1031-sudo/aiy-warp/install/src/internal/frontmatter"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/kit"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/redact"
	"gopkg.in/yaml.v3"
)

// check2 — agent frontmatter parses (legacy or canonical); canonical name must
// equal the filename.
func check2(k *kit.Kit, repoRoot string) Check {
	c := Check{ID: 2, Name: "agent frontmatter parses; name = filename"}
	failed := 0
	for _, rel := range k.Agents {
		content, err := os.ReadFile(filepath.Join(repoRoot, filepath.FromSlash(rel)))
		if err != nil {
			c.Hints = append(c.Hints, rel+": unreadable")
			failed++
			continue
		}
		a, err := frontmatter.ParseAgent(rel, string(content))
		if err != nil {
			c.Hints = append(c.Hints, err.Error())
			failed++
			continue
		}
		if a.Canonical {
			want := strings.TrimSuffix(filepath.Base(rel), filepath.Ext(rel))
			if a.Name != want {
				c.Hints = append(c.Hints, fmt.Sprintf("%s: canonical name %q != filename", rel, a.Name))
				failed++
			}
			if a.Department == "" || a.Rank == "" {
				c.Hints = append(c.Hints, rel+": canonical file missing department/rank")
				failed++
			}
		}
	}
	if failed == 0 {
		c.Status = Pass
		c.Hints = append(c.Hints, fmt.Sprintf("%d agent files parse (legacy or canonical)", len(k.Agents)))
	} else {
		c.Status = Fail
	}
	return c
}

// check3 — org-chart integrity: every reports_to/department/rank resolves.
func check3(b *kit.Bundles) Check {
	c := Check{ID: 3, Name: "org-chart integrity (bundles.yaml)"}
	var problems []string
	for _, a := range b.Agents {
		if !validRank(a.Rank) {
			problems = append(problems, a.Name+": invalid rank "+a.Rank)
		}
		if !b.HasDept(a.Department) {
			problems = append(problems, a.Name+": unknown department "+a.Department)
		}
		if a.ReportsTo != nil && *a.ReportsTo != "" && b.Find(*a.ReportsTo) == nil {
			problems = append(problems, a.Name+": dangling reports_to "+*a.ReportsTo)
		}
	}
	if len(problems) == 0 {
		c.Status = Pass
		c.Hints = append(c.Hints, fmt.Sprintf("%d agents across %d departments", len(b.Agents), len(b.Departments)))
	} else {
		c.Status = Fail
		c.Hints = problems
	}
	return c
}

// check4 — skill references resolve to skills/<name>/SKILL.md and known owners.
func check4(k *kit.Kit, b *kit.Bundles) Check {
	c := Check{ID: 4, Name: "skill references resolve"}
	have := map[string]bool{}
	for _, s := range k.Skills {
		have[filepath.Base(filepath.Dir(filepath.FromSlash(s)))] = true
	}
	var problems []string
	for _, s := range b.Skills {
		if !have[s.Name] {
			problems = append(problems, s.Name+": skills/<name>/SKILL.md missing")
		}
		if s.OwnedBy != "" && b.Find(s.OwnedBy) == nil {
			problems = append(problems, s.Name+": owned_by "+s.OwnedBy+" is not an agent")
		}
	}
	for _, a := range b.Agents {
		for _, s := range a.Skills {
			if !b.HasSkill(s) {
				problems = append(problems, a.Name+": unknown skill "+s)
			}
		}
	}
	if len(problems) == 0 {
		c.Status = Pass
		c.Hints = append(c.Hints, fmt.Sprintf("%d skills resolve", len(b.Skills)))
	} else {
		c.Status = Fail
		c.Hints = problems
	}
	return c
}

// check5 — no credential values in the payload (CLI_SPEC.md §5.4).
func check5(k *kit.Kit, repoRoot string) Check {
	c := Check{ID: 5, Name: "no credential values in payload"}
	scan := map[string]string{}
	for _, rel := range k.All() {
		if b, err := os.ReadFile(filepath.Join(repoRoot, filepath.FromSlash(rel))); err == nil {
			scan[rel] = string(b)
		}
	}
	rep := redact.Scan(scan, nil)
	if rep.HasBlocks() {
		c.Status = Fail
		c.Hints = append(c.Hints, "credential values found:")
		for _, f := range rep.Blocks {
			c.Hints = append(c.Hints, fmt.Sprintf("  %s:%d (%s)", f.File, f.Line, f.Pattern))
		}
	} else {
		c.Status = Pass
		c.Hints = append(c.Hints, "no credential values")
	}
	for _, f := range rep.Exports {
		c.Warns = append(c.Warns, fmt.Sprintf("%s:%d numeric id — P1 parameterization debt", f.File, f.Line))
	}
	for _, f := range capWarns(rep.Warns, 8) {
		c.Warns = append(c.Warns, fmt.Sprintf("%s:%d %s", f.File, f.Line, f.Pattern))
	}
	return c
}

// check6 — warp.config paths resolve on this machine.
func check6(configPath string) Check {
	c := Check{ID: 6, Name: "warp.config paths resolve"}
	b, err := os.ReadFile(configPath)
	if err != nil {
		c.Status = Pass
		c.Hints = append(c.Hints, "no warp.config — $HOME defaults in use")
		return c
	}
	var cfg struct {
		Paths struct {
			Home         string `yaml:"home"`
			Workspace    string `yaml:"workspace"`
			ObsidianRoot string `yaml:"obsidian_root"`
		} `yaml:"paths"`
	}
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		c.Status = Fail
		c.Hints = append(c.Hints, "unparseable: "+err.Error())
		return c
	}
	var missing []string
	for _, p := range []string{cfg.Paths.Home, cfg.Paths.Workspace, cfg.Paths.ObsidianRoot} {
		if p == "" {
			continue
		}
		if _, err := os.Stat(expandHome(p)); err != nil {
			missing = append(missing, p)
		}
	}
	if len(missing) == 0 {
		c.Status = Pass
		c.Hints = append(c.Hints, "configured paths resolve")
	} else {
		c.Status = Fail
		c.Hints = append(c.Hints, "missing paths: "+strings.Join(missing, ", "))
	}
	return c
}

func validRank(r string) bool {
	return r == "primary" || r == "head" || r == "executor"
}

func capWarns(fs []redact.Finding, n int) []redact.Finding {
	if len(fs) <= n {
		return fs
	}
	return fs[:n]
}

func expandHome(p string) string {
	if p == "~" || strings.HasPrefix(p, "~/") {
		if h, err := os.UserHomeDir(); err == nil {
			return filepath.Join(h, strings.TrimPrefix(p, "~/"))
		}
	}
	return p
}
