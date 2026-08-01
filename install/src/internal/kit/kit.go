package kit

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Kit is the discovered payload of the warp repo.
type Kit struct {
	Root      string   // repo root (source of truth)
	Agents    []string // repo-relative paths, e.g. agents/aiy.md
	Skills    []string // repo-relative paths, e.g. skills/aiy-messaging/SKILL.md
	Templates []string
	Playbooks []string
	Rev       string // source git revision at discovery time
}

// Selector narrows a bundle (CLI_SPEC.md §1.1). Mutual exclusion of Agent and
// Team is validated inside Resolve.
type Selector struct {
	Agent string
	Team  string
}

// Discover validates repoRoot and inventories the payload directories.
func Discover(repoRoot string) (*Kit, error) {
	k := &Kit{Root: repoRoot, Rev: GitRev(repoRoot)}
	var err error
	if k.Agents, err = globRel(repoRoot, "agents", "*.md"); err != nil {
		return nil, err
	}
	if k.Skills, err = globRel(repoRoot, "skills", "*", "SKILL.md"); err != nil {
		return nil, err
	}
	if k.Templates, err = globRel(repoRoot, "templates", "*.md"); err != nil {
		return nil, err
	}
	if k.Playbooks, err = globRel(repoRoot, "playbooks", "*.md"); err != nil {
		return nil, err
	}
	if len(k.Agents) == 0 {
		return nil, fmt.Errorf("%s: no agents found — not a warp kit?", repoRoot)
	}
	if _, err := os.Stat(filepath.Join(repoRoot, "install", "bundles.yaml")); err != nil {
		return nil, fmt.Errorf("%s: install/bundles.yaml missing — not a warp kit?", repoRoot)
	}
	return k, nil
}

// All returns every payload file, sorted.
func (k *Kit) All() []string {
	return sortedUnique(append(append(append([]string{}, k.Agents...), k.Skills...),
		append(k.Templates, k.Playbooks...)...))
}

// Resolve maps a selector to a sorted bundle file set. The empty selector is
// the full kit. Unknown agent/team and combined selectors are usage errors.
func (k *Kit) Resolve(b *Bundles, sel Selector) ([]string, error) {
	switch {
	case sel.Agent != "" && sel.Team != "":
		return nil, fmt.Errorf("--agent and --team are mutually exclusive")
	case sel.Agent != "":
		return k.resolveAgent(b, sel.Agent)
	case sel.Team != "":
		return k.resolveTeam(b, sel.Team)
	default:
		return k.All(), nil
	}
}

func (k *Kit) resolveAgent(b *Bundles, name string) ([]string, error) {
	a := b.Find(name)
	if a == nil {
		return nil, fmt.Errorf("unknown agent %q — see 'aiy warp list'", name)
	}
	files := []string{filepath.ToSlash(filepath.Join("agents", name+".md"))}
	for _, s := range append(append([]string{}, a.Skills...), b.SkillOwnedBy(name)...) {
		files = append(files, k.skillFiles(s)...)
	}
	return sortedUnique(files), nil
}

func (k *Kit) resolveTeam(b *Bundles, dept string) ([]string, error) {
	if !b.HasDept(dept) {
		return nil, fmt.Errorf("unknown team %q — departments: %s", dept, strings.Join(b.Departments, ", "))
	}
	var files []string
	for _, n := range b.Team(dept) {
		files = append(files, filepath.ToSlash(filepath.Join("agents", n+".md")))
	}
	for _, s := range b.Skills {
		for _, d := range s.Departments {
			if d == dept {
				files = append(files, k.skillFiles(s.Name)...)
			}
		}
	}
	return sortedUnique(files), nil
}

// skillFiles returns the kit's files for a skill directory.
func (k *Kit) skillFiles(name string) []string {
	var out []string
	for _, p := range k.Skills {
		if filepath.Base(filepath.Dir(filepath.FromSlash(p))) == name {
			out = append(out, p)
		}
	}
	return out
}

// GitRev reads the repo's HEAD commit without shelling out.
func GitRev(root string) string {
	head, err := os.ReadFile(filepath.Join(root, ".git", "HEAD"))
	if err != nil {
		return "unknown"
	}
	s := strings.TrimSpace(string(head))
	if strings.HasPrefix(s, "ref: ") {
		if b, err := os.ReadFile(filepath.Join(root, ".git", strings.TrimPrefix(s, "ref: "))); err == nil {
			return strings.TrimSpace(string(b))
		}
	}
	return s
}

func globRel(root string, parts ...string) ([]string, error) {
	pattern := filepath.Join(append([]string{root}, parts...)...)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, m := range matches {
		rel, err := filepath.Rel(root, m)
		if err != nil {
			return nil, err
		}
		out = append(out, filepath.ToSlash(rel))
	}
	sort.Strings(out)
	return out, nil
}

func sortedUnique(in []string) []string {
	sort.Strings(in)
	out := in[:0]
	for _, s := range in {
		if len(out) == 0 || out[len(out)-1] != s {
			out = append(out, s)
		}
	}
	return out
}
