package workspace

import (
	"os"
	"path/filepath"
	"testing"
)

func mkTemplate(t *testing.T) string {
	t.Helper()
	tpl := t.TempDir()
	for _, d := range []string{"00-Inbox", "01-Projects"} {
		if err := os.MkdirAll(filepath.Join(tpl, d), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(tpl, d, "README.md"), []byte("# "+d), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return tpl
}

func TestRunCreates(t *testing.T) {
	tpl := mkTemplate(t)
	home := t.TempDir()
	res, err := Run(Options{Home: home, TemplateDir: tpl})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Dirs) != 2 || res.Dirs[0] != "00-Inbox" || res.Dirs[1] != "01-Projects" {
		t.Fatalf("dirs = %v", res.Dirs)
	}
	if len(res.Copied) != 2 {
		t.Fatalf("copied = %v", res.Copied)
	}
	if res.Noop {
		t.Fatal("first run must not be a no-op")
	}
	for _, d := range res.Dirs {
		if _, err := os.Stat(filepath.Join(home, d)); err != nil {
			t.Errorf("dir %s not created: %v", d, err)
		}
	}
	readme, err := os.ReadFile(filepath.Join(home, "00-Inbox", "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(readme) != "# 00-Inbox" {
		t.Fatalf("README content = %q", readme)
	}
}

func TestRunNoop(t *testing.T) {
	tpl := mkTemplate(t)
	home := t.TempDir()
	if _, err := Run(Options{Home: home, TemplateDir: tpl}); err != nil {
		t.Fatal(err)
	}
	res, err := Run(Options{Home: home, TemplateDir: tpl})
	if err != nil {
		t.Fatal(err)
	}
	if !res.Noop {
		t.Fatal("second run must be a no-op")
	}
	if len(res.Existing) == 0 {
		t.Fatal("existing components must be reported")
	}
}

func TestRunDryRun(t *testing.T) {
	tpl := mkTemplate(t)
	home := t.TempDir()
	res, err := Run(Options{Home: home, TemplateDir: tpl, DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Dirs) != 2 {
		t.Fatalf("dry-run dirs = %v", res.Dirs)
	}
	for _, d := range res.Dirs {
		if _, err := os.Stat(filepath.Join(home, d)); !os.IsNotExist(err) {
			t.Errorf("dry-run must not create %s", d)
		}
	}
}

func TestSummary(t *testing.T) {
	if got := Summary([]string{"00-Inbox", "01-Projects"}); got != "00-Inbox, 01-Projects" {
		t.Fatalf("Summary = %q", got)
	}
	if got := Summary(nil); got != "—" {
		t.Fatalf("Summary(nil) = %q", got)
	}
}
