package migrate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/anco1031-sudo/aiy-warp/install/src/internal/frontmatter"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/kit"
)

func sptr(s string) *string { return &s }

func testBundles() *kit.Bundles {
	return &kit.Bundles{
		Departments: []string{"core", "tech", "trading"},
		Agents: []kit.Agent{
			{Name: "aiy", Rank: "head", Department: "core"},
			{Name: "lin", Rank: "head", Department: "tech", ReportsTo: sptr("aiy")},
			{Name: "kwan", Rank: "head", Department: "trading", ReportsTo: sptr("aiy")},
			{Name: "bee", Rank: "executor", Department: "trading", ReportsTo: sptr("kwan")},
		},
	}
}

func TestDeriveIdentity(t *testing.T) {
	display, role := deriveIdentity("lin", "LIN (หลิน) — The Elite Product Owner & Governor.")
	if display != "หลิน (Lin)" {
		t.Errorf("display = %q, want %q", display, "หลิน (Lin)")
	}
	if role != "The Elite Product Owner & Governor" {
		t.Errorf("role = %q", role)
	}
}

func TestStampPreservesBodyAndPermission(t *testing.T) {
	content := "---\ndescription: \"LIN (หลิน) — The Elite Product Owner & Governor.\"\nmode: subagent\nmodel: opencode/deepseek-v4-flash-free\ncolor: \"#E53935\"\npermission:\n  edit: true\n---\n# Body\n1. ROLE & MINDSET\nYou are Lin."
	a, err := frontmatter.ParseAgent("agents/lin.md", content)
	if err != nil {
		t.Fatal(err)
	}
	if a.Canonical {
		t.Fatal("legacy file must not be canonical")
	}
	next, err := stamp(a, testBundles().Find("lin"), content)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"warp_version: 1",
		"name: lin",
		"หลิน (Lin)",
		"department: tech",
		"rank: head",
		"reports_to: aiy",
		"platform_targets",
		"permission:",
		"edit: true",
	} {
		if !strings.Contains(next, want) {
			t.Errorf("stamped frontmatter missing %q:\n%s", want, next)
		}
	}
	// Body preserved byte-for-byte.
	if !strings.HasSuffix(next, "# Body\n1. ROLE & MINDSET\nYou are Lin.") {
		t.Errorf("body not preserved:\n%s", next)
	}
	// Re-parse → canonical.
	re, err := frontmatter.ParseAgent("agents/lin.md", next)
	if err != nil {
		t.Fatal(err)
	}
	if !re.Canonical || re.WarpVersion != 1 {
		t.Fatalf("stamped file must re-parse canonical: %+v", re)
	}
}

func TestRunStampsAndSkips(t *testing.T) {
	dir := t.TempDir()
	agentsDir := filepath.Join(dir, "agents")
	if err := os.MkdirAll(agentsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	legacy := "---\ndescription: \"KWAN (ขวัญ) — The Head Trader.\"\nmode: subagent\n---\nbody-kwan"
	if err := os.WriteFile(filepath.Join(agentsDir, "kwan.md"), []byte(legacy), 0o644); err != nil {
		t.Fatal(err)
	}
	canonical := "---\nwarp_version: 1\nname: aiy\nrole: Muse\n---\nbody-aiy"
	if err := os.WriteFile(filepath.Join(agentsDir, "aiy.md"), []byte(canonical), 0o644); err != nil {
		t.Fatal(err)
	}

	b := testBundles()
	res, err := Run(Options{AgentsDir: agentsDir, Bundles: b, Backup: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Stamped) != 1 || res.Stamped[0] != "kwan" {
		t.Fatalf("stamped = %v", res.Stamped)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "aiy" {
		t.Fatalf("skipped = %v", res.Skipped)
	}
	if res.Noop() {
		t.Fatal("run that stamped must not be a no-op")
	}
	// .bak written with the original content.
	bak, err := os.ReadFile(filepath.Join(agentsDir, "kwan.md.bak"))
	if err != nil {
		t.Fatal(err)
	}
	if string(bak) != legacy {
		t.Fatalf(".bak content mismatch")
	}
	got, err := os.ReadFile(filepath.Join(agentsDir, "kwan.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "warp_version: 1") {
		t.Fatalf("kwan.md not stamped:\n%s", got)
	}

	// Second run → fully idempotent.
	res2, err := Run(Options{AgentsDir: agentsDir, Bundles: b})
	if err != nil {
		t.Fatal(err)
	}
	if len(res2.Stamped) != 0 || len(res2.Skipped) != 2 {
		t.Fatalf("second run: stamped=%v skipped=%v", res2.Stamped, res2.Skipped)
	}
	if !res2.Noop() {
		t.Fatal("second run must be a no-op")
	}
}

func TestRunDryRun(t *testing.T) {
	dir := t.TempDir()
	agentsDir := filepath.Join(dir, "agents")
	if err := os.MkdirAll(agentsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	legacy := "---\ndescription: \"BEE (บี) — The Fundamental Analyst.\"\n---\nbody-bee"
	if err := os.WriteFile(filepath.Join(agentsDir, "bee.md"), []byte(legacy), 0o644); err != nil {
		t.Fatal(err)
	}
	res, err := Run(Options{AgentsDir: agentsDir, Bundles: testBundles(), DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Stamped) != 1 {
		t.Fatalf("dry-run must report stamped, got %v", res.Stamped)
	}
	got, err := os.ReadFile(filepath.Join(agentsDir, "bee.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != legacy {
		t.Fatalf("dry-run must not modify the file")
	}
	if _, err := os.Stat(filepath.Join(agentsDir, "bee.md.bak")); !os.IsNotExist(err) {
		t.Fatalf("dry-run must not write .bak")
	}
}

func TestResultNoop(t *testing.T) {
	if !(&Result{}).Noop() {
		t.Error("empty result must be no-op")
	}
	if (&Result{Stamped: []string{"x"}}).Noop() {
		t.Error("stamped result must not be no-op")
	}
	if (&Result{Errors: []string{"e"}}).Noop() {
		t.Error("errored result must not be no-op")
	}
}
