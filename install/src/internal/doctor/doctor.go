// Package doctor implements `aiy warp doctor` — the six integrity checks of
// CLI_SPEC.md §6.3. Exit: 0 clean, 3 drift, 1 runtime error.
package doctor

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/anco1031-sudo/aiy-warp/install/src/internal/kit"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/manifest"
)

// Status of a check.
type Status string

const (
	Pass Status = "PASS"
	Fail Status = "FAIL"
)

// Check is one doctor check result.
type Check struct {
	ID     int      `json:"id"`
	Name   string   `json:"name"`
	Status Status   `json:"status"`
	Hints  []string `json:"hints,omitempty"`
	Warns  []string `json:"warns,omitempty"`
}

// Result is the full doctor report.
type Result struct {
	Checks []Check `json:"checks"`
	Exit   int     `json:"exit"` // 0 clean / 3 drift
}

// Options configures the doctor run.
type Options struct {
	RepoRoot   string
	DestDir    string
	StateDir   string
	ConfigPath string
}

// Run executes all six checks. Hard runtime failures return an error (exit 1);
// drift is encoded in Result.Exit (3).
func Run(o Options) (*Result, error) {
	k, err := kit.Discover(o.RepoRoot)
	if err != nil {
		return nil, err
	}
	b, err := kit.LoadBundles(filepath.Join(o.RepoRoot, "install", "bundles.yaml"))
	if err != nil {
		return nil, err
	}

	res := &Result{}
	res.Checks = append(res.Checks,
		check1(k, o.DestDir, o.StateDir),
		check2(k, o.RepoRoot),
		check3(b),
		check4(k, b),
		check5(k, o.RepoRoot),
		check6(o.ConfigPath),
	)
	res.Exit = 0
	for _, c := range res.Checks {
		if c.Status == Fail {
			res.Exit = 3
			break
		}
	}
	return res, nil
}

// check1 — all files recorded in warp.lock are present at the destination with
// drift hashes matching the repo (source of truth, CLI_SPEC.md §6.1). When a
// lock exists, only its file set is verified: kit files intentionally not
// installed (selective install) are warnings, not drift.
func check1(k *kit.Kit, destDir, stateDir string) Check {
	c := Check{ID: 1, Name: "installed files present & hashes match"}
	lock, err := manifest.Load(filepath.Join(stateDir, "warp.lock"))
	if err != nil {
		c.Status = Fail
		c.Hints = append(c.Hints, "warp.lock unreadable: "+err.Error())
		return c
	}
	rels := k.All()
	if lock != nil && len(lock.Files) > 0 {
		rels = make([]string, 0, len(lock.Files))
		for r := range lock.Files {
			rels = append(rels, r)
		}
		sort.Strings(rels)
		if notInstalled := kitFilesNotIn(k.All(), rels); len(notInstalled) > 0 {
			shown, suffix := notInstalled, ""
			if len(shown) > 5 {
				shown, suffix = shown[:5], fmt.Sprintf(" … %d more", len(notInstalled)-5)
			}
			c.Warns = append(c.Warns, fmt.Sprintf(
				"%d kit file(s) not installed (selective install): %s%s — not drift",
				len(notInstalled), strings.Join(shown, ", "), suffix))
		}
	} else {
		c.Warns = append(c.Warns, "no warp.lock — installed state unknown (run 'aiy warp install opencode')")
	}
	missing, drifted, repoUnreadable := 0, 0, 0
	for _, rel := range rels {
		dest, err := os.ReadFile(filepath.Join(destDir, filepath.FromSlash(rel)))
		if err != nil {
			missing++
			continue
		}
		repo, err := os.ReadFile(filepath.Join(k.Root, filepath.FromSlash(rel)))
		if err != nil {
			repoUnreadable++
			c.Hints = append(c.Hints, rel+": repo file unreadable ("+err.Error()+")")
			continue
		}
		if manifest.KitHash(string(dest), rel) != manifest.KitHash(string(repo), rel) {
			drifted++
		}
	}
	if missing == 0 && drifted == 0 && repoUnreadable == 0 {
		c.Status = Pass
		c.Hints = append(c.Hints, fmt.Sprintf("%d files in sync", len(rels)))
	} else {
		c.Status = Fail
		c.Hints = append(c.Hints, fmt.Sprintf(
			"%d missing, %d drifted, %d repo-unreadable — run 'aiy warp install opencode' (--force if needed)",
			missing, drifted, repoUnreadable))
	}
	return c
}

// kitFilesNotIn returns the sorted kit rels absent from the given subset.
func kitFilesNotIn(all, subset []string) []string {
	in := make(map[string]bool, len(subset))
	for _, r := range subset {
		in[r] = true
	}
	var out []string
	for _, r := range all {
		if !in[r] {
			out = append(out, r)
		}
	}
	return out
}
