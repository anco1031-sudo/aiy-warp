package doctor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/anco1031-sudo/aiy-warp/install/src/internal/manifest"
)

func testKit(t *testing.T, root string) {
	t.Helper()
	for _, d := range []string{"agents", "install", "skills", "templates", "playbooks"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	agents := map[string]string{
		"aiy": "---\ndescription: \"AIY (อัย) — test\"\nmode: primary\n---\naiy body\n",
		"lin": "---\ndescription: \"LIN (หลิน) — test\"\nmode: subagent\n---\nlin body\n",
	}
	for name, content := range agents {
		if err := os.WriteFile(filepath.Join(root, "agents", name+".md"), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	bundles := "warp_version: 1\ndepartments: [core, tech]\nagents:\n  - name: aiy\n    rank: primary\n    department: core\n    reports_to: null\n    skills: []\n  - name: lin\n    rank: head\n    department: tech\n    reports_to: aiy\n    skills: []\nskills: []\n"
	if err := os.WriteFile(filepath.Join(root, "install", "bundles.yaml"), []byte(bundles), 0o644); err != nil {
		t.Fatal(err)
	}
}

func read(t *testing.T, p string) []byte {
	t.Helper()
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func saveLock(t *testing.T, stateDir string, files map[string]string) {
	t.Helper()
	lock := manifest.New("opencode", "test", files)
	if err := manifest.Save(filepath.Join(stateDir, "warp.lock"), lock); err != nil {
		t.Fatal(err)
	}
}

func TestCheck1UsesLockFileSet(t *testing.T) {
	root := t.TempDir()
	testKit(t, root)
	dest := filepath.Join(t.TempDir(), ".config", "opencode")
	state := filepath.Join(t.TempDir(), ".config", "aiy-warp")
	// Selective install: only aiy.md is on the host and the lock records only it.
	agentDir := filepath.Join(dest, "agents")
	if err := os.MkdirAll(agentDir, 0o755); err != nil {
		t.Fatal(err)
	}
	aiy := read(t, filepath.Join(root, "agents", "aiy.md"))
	if err := os.WriteFile(filepath.Join(agentDir, "aiy.md"), aiy, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(state, 0o755); err != nil {
		t.Fatal(err)
	}
	saveLock(t, state, map[string]string{"agents/aiy.md": manifest.KitHash(string(aiy), "agents/aiy.md")})

	res, err := Run(Options{RepoRoot: root, DestDir: dest, StateDir: state, ConfigPath: filepath.Join(root, "nonexistent")})
	if err != nil {
		t.Fatal(err)
	}
	if res.Exit != 0 {
		t.Fatalf("selective install must not report drift: exit %d, checks %+v", res.Exit, res.Checks)
	}
	c := res.Checks[0]
	if c.Status != Pass {
		t.Fatalf("check1 = %s, want PASS", c.Status)
	}
	if len(c.Warns) == 0 {
		t.Fatal("expected a warn about uninstalled kit files")
	}
}

func TestCheck1RepoUnreadableIsFail(t *testing.T) {
	root := t.TempDir()
	testKit(t, root)
	dest := filepath.Join(t.TempDir(), ".config", "opencode")
	state := filepath.Join(t.TempDir(), ".config", "aiy-warp")
	agentDir := filepath.Join(dest, "agents")
	if err := os.MkdirAll(agentDir, 0o755); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{}
	for _, name := range []string{"aiy", "lin"} {
		content := read(t, filepath.Join(root, "agents", name+".md"))
		if err := os.WriteFile(filepath.Join(agentDir, name+".md"), content, 0o644); err != nil {
			t.Fatal(err)
		}
		rel := "agents/" + name + ".md"
		files[rel] = manifest.KitHash(string(content), rel)
	}
	if err := os.MkdirAll(state, 0o755); err != nil {
		t.Fatal(err)
	}
	saveLock(t, state, files)
	// Make a repo file unreadable: replace it with a directory.
	if err := os.Remove(filepath.Join(root, "agents", "lin.md")); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(root, "agents", "lin.md"), 0o755); err != nil {
		t.Fatal(err)
	}

	res, err := Run(Options{RepoRoot: root, DestDir: dest, StateDir: state, ConfigPath: filepath.Join(root, "nonexistent")})
	if err != nil {
		t.Fatal(err)
	}
	if res.Exit != 3 {
		t.Fatalf("repo-unreadable file must fail doctor: exit %d, checks %+v", res.Exit, res.Checks)
	}
	c := res.Checks[0]
	if c.Status != Fail {
		t.Fatalf("check1 = %s, want FAIL", c.Status)
	}
}
