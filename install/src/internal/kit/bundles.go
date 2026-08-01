// Package kit resolves the warp payload (agents, skills, templates, playbooks)
// and computes install/export bundles per CLI_SPEC.md §4.
package kit

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Bundles is the interim org-chart manifest (install/bundles.yaml).
type Bundles struct {
	WarpVersion int      `yaml:"warp_version"`
	Departments []string `yaml:"departments"`
	Agents      []Agent  `yaml:"agents"`
	Skills      []Skill  `yaml:"skills"`
}

// Agent is one org-chart entry.
type Agent struct {
	Name       string   `yaml:"name"`
	Rank       string   `yaml:"rank"`
	Department string   `yaml:"department"`
	ReportsTo  *string  `yaml:"reports_to"`
	Skills     []string `yaml:"skills"`
}

// Skill is one skill-ownership entry.
type Skill struct {
	Name        string   `yaml:"name"`
	OwnedBy     string   `yaml:"owned_by"`
	Departments []string `yaml:"departments"`
}

// LoadBundles reads and parses install/bundles.yaml.
func LoadBundles(path string) (*Bundles, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var out Bundles
	if err := yaml.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Find returns the org entry for an agent name, or nil.
func (b *Bundles) Find(name string) *Agent {
	for i := range b.Agents {
		if b.Agents[i].Name == name {
			return &b.Agents[i]
		}
	}
	return nil
}

// Team returns the agent names of a department (head + executors), sorted.
func (b *Bundles) Team(dept string) []string {
	var out []string
	for _, a := range b.Agents {
		if a.Department == dept {
			out = append(out, a.Name)
		}
	}
	sort.Strings(out)
	return out
}

// SkillOwnedBy returns skills whose owner is the given agent, sorted.
func (b *Bundles) SkillOwnedBy(agent string) []string {
	var out []string
	for _, s := range b.Skills {
		if s.OwnedBy == agent {
			out = append(out, s.Name)
		}
	}
	sort.Strings(out)
	return out
}

// ResolveTeam maps a team selector to a department: an existing department
// name passes through; a department-head agent name (e.g. "kwan") resolves to
// its department; anything else is returned unchanged (caller reports unknown).
func (b *Bundles) ResolveTeam(name string) string {
	if b.HasDept(name) {
		return name
	}
	if a := b.Find(name); a != nil && a.Rank == "head" {
		return a.Department
	}
	return name
}

// HasDept reports whether dept is a declared department.
func (b *Bundles) HasDept(dept string) bool {
	for _, d := range b.Departments {
		if d == dept {
			return true
		}
	}
	return false
}

// HasSkill reports whether a skill is declared.
func (b *Bundles) HasSkill(name string) bool {
	for _, s := range b.Skills {
		if s.Name == name {
			return true
		}
	}
	return false
}

// AgentNames returns all agent names, sorted.
func (b *Bundles) AgentNames() []string {
	out := make([]string, 0, len(b.Agents))
	for _, a := range b.Agents {
		out = append(out, a.Name)
	}
	sort.Strings(out)
	return out
}

// SkillNames returns all skill names, sorted.
func (b *Bundles) SkillNames() []string {
	out := make([]string, 0, len(b.Skills))
	for _, s := range b.Skills {
		out = append(out, s.Name)
	}
	sort.Strings(out)
	return out
}

// Describe renders a compact org-chart line for `aiy warp list`.
func (b *Bundles) Describe(name string) string {
	a := b.Find(name)
	if a == nil {
		return fmt.Sprintf("%-6s (unknown)", name)
	}
	reports := ""
	if a.ReportsTo != nil {
		reports = *a.ReportsTo
	}
	if reports == "" {
		reports = "—"
	}
	return strings.TrimSpace(fmt.Sprintf("%-6s %-9s %-9s %s", a.Name, a.Rank, a.Department, reports))
}
