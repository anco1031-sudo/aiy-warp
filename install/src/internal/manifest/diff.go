package manifest

import (
	"sort"
	"strings"

	"github.com/anco1031-sudo/aiy-warp/install/src/internal/frontmatter"
	"gopkg.in/yaml.v3"
)

// State of a file in the 3-way diff (CLI_SPEC.md §6.2).
type State string

const (
	StateNew     State = "new"     // in repo, not installed
	StateSame    State = "same"    // both present, drift hashes equal
	StateDrifted State = "drifted" // both present, drift hashes differ
	StateOrphan  State = "orphan"  // installed, not in repo
)

// Entry is one file's diff state.
type Entry struct {
	Rel      string
	State    State
	RepoHash string
	HostHash string
}

// HostLocalFields are the opencode-specific frontmatter keys exempted from
// drift detection (CLI_SPEC.md §6.2 merge policy): host tuning of these is
// runtime configuration, not kit content.
var HostLocalFields = []string{"model", "permission"}

// KitHash returns the drift-comparison hash for a payload file. Agent files
// are normalized by stripping host-local frontmatter fields so that host
// tuning of `model`/`permission` never registers as drift.
func KitHash(content, rel string) string {
	if strings.HasPrefix(rel, "agents/") {
		fm, body, err := frontmatter.SplitFrontmatter(content)
		if err == nil && fm != "" {
			var m map[string]any
			if err := yaml.Unmarshal([]byte(fm), &m); err == nil {
				for _, k := range HostLocalFields {
					delete(m, k)
				}
				if b, err := yaml.Marshal(m); err == nil {
					return Sha256Hex("---\n" + string(b) + "---\n" + body)
				}
			}
		}
	}
	return Sha256Hex(content)
}

// Diff compares repo drift hashes against host drift hashes. Only repo files
// are returned as entries; orphans are reported separately.
func Diff(repo, host map[string]string) []Entry {
	rels := make([]string, 0, len(repo))
	for p := range repo {
		rels = append(rels, p)
	}
	sort.Strings(rels)
	out := make([]Entry, 0, len(rels))
	for _, p := range rels {
		e := Entry{Rel: p, RepoHash: repo[p]}
		hh, ok := host[p]
		switch {
		case !ok:
			e.State = StateNew
		case hh == repo[p]:
			e.State = StateSame
		default:
			e.State = StateDrifted
			e.HostHash = hh
		}
		out = append(out, e)
	}
	return out
}

// Orphans returns installed files not present in the repo, sorted.
func Orphans(repo, host map[string]string) []string {
	var out []string
	for p := range host {
		if _, ok := repo[p]; !ok {
			out = append(out, p)
		}
	}
	sort.Strings(out)
	return out
}
