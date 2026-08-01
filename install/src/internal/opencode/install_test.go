package opencode

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/anco1031-sudo/aiy-warp/install/src/internal/kit"
)

func testKit(t *testing.T, root string) *kit.Bundles {
	t.Helper()
	for _, d := range []string{"agents", "install", "skills", "templates", "playbooks"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	agent := "---\ndescription: \"AIY (อัย) — test\"\nmode: primary\n---\nbody\n"
	if err := os.WriteFile(filepath.Join(root, "agents", "aiy.md"), []byte(agent), 0o644); err != nil {
		t.Fatal(err)
	}
	bundles := "warp_version: 1\ndepartments: [core]\nagents:\n  - name: aiy\n    rank: primary\n    department: core\n    reports_to: null\n    skills: []\nskills: []\n"
	if err := os.WriteFile(filepath.Join(root, "install", "bundles.yaml"), []byte(bundles), 0o644); err != nil {
		t.Fatal(err)
	}
	b, err := kit.LoadBundles(filepath.Join(root, "install", "bundles.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestInstallDoesNotWriteThroughSymlink(t *testing.T) {
	root := t.TempDir()
	b := testKit(t, root)
	dest := filepath.Join(t.TempDir(), ".config", "opencode")
	state := filepath.Join(t.TempDir(), ".config", "aiy-warp")
	target := filepath.Join(t.TempDir(), "precious.txt")
	if err := os.WriteFile(target, []byte("PRECIOUS DATA"), 0o644); err != nil {
		t.Fatal(err)
	}
	agentDir := filepath.Join(dest, "agents")
	if err := os.MkdirAll(agentDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, filepath.Join(agentDir, "aiy.md")); err != nil {
		t.Fatal(err)
	}

	res, err := RunInstall(InstallOptions{
		RepoRoot: root, DestDir: dest, StateDir: state, Bundles: b,
		Selector: kit.Selector{}, Force: true,
	})
	if err != nil {
		t.Fatalf("install failed: %v", err)
	}

	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "PRECIOUS DATA" {
		t.Fatalf("symlink target was written through: %q", got)
	}
	fi, err := os.Lstat(filepath.Join(agentDir, "aiy.md"))
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode()&os.ModeSymlink != 0 {
		t.Fatal("dest file must be a regular file after install, not a symlink")
	}
	if len(res.Notes) == 0 {
		t.Fatal("expected a symlink-replaced note")
	}
}

func TestDryRunReportsNoopWhenInSync(t *testing.T) {
	root := t.TempDir()
	b := testKit(t, root)
	dest := filepath.Join(t.TempDir(), ".config", "opencode")
	state := filepath.Join(t.TempDir(), ".config", "aiy-warp")
	opts := InstallOptions{
		RepoRoot: root, DestDir: dest, StateDir: state, Bundles: b, Selector: kit.Selector{},
	}
	if _, err := RunInstall(opts); err != nil {
		t.Fatal(err)
	}
	dry := opts
	dry.DryRun = true
	res, err := RunInstall(dry)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Noop {
		t.Fatal("dry-run on already-in-sync state must report Noop")
	}
}
