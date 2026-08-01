package manifest

import (
	"reflect"
	"testing"

	"github.com/anco1031-sudo/aiy-warp/install/src/internal/frontmatter"
	"gopkg.in/yaml.v3"
)

func TestDiffStates(t *testing.T) {
	repo := map[string]string{"agents/a.md": "h1", "agents/b.md": "h2", "agents/c.md": "h3"}
	host := map[string]string{"agents/a.md": "h1", "agents/b.md": "DIFF", "agents/d.md": "h4"}

	entries := Diff(repo, host)
	got := map[string]State{}
	for _, e := range entries {
		got[e.Rel] = e.State
	}
	want := map[string]State{"agents/a.md": StateSame, "agents/b.md": StateDrifted, "agents/c.md": StateNew}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("diff states = %v, want %v", got, want)
	}

	orphans := Orphans(repo, host)
	if len(orphans) != 1 || orphans[0] != "agents/d.md" {
		t.Fatalf("orphans = %v", orphans)
	}
}

func TestKitHashIgnoresHostLocalFields(t *testing.T) {
	repo := "---\ndescription: \"LIN (หลิน) — PO\"\nmode: subagent\nmodel: opencode/old\ncolor: \"#E53935\"\n---\nbody"
	tuned := "---\ndescription: \"LIN (หลิน) — PO\"\nmode: subagent\nmodel: opencode/host-custom\npermission:\n  edit: allow\ncolor: \"#E53935\"\n---\nbody"

	if KitHash(repo, "agents/lin.md") != KitHash(tuned, "agents/lin.md") {
		t.Fatal("host-local model/permission tuning must not change the drift hash")
	}
}

func TestKitHashCatchesKitDrift(t *testing.T) {
	repo := "---\ndescription: \"LIN (หลิน) — PO\"\n---\nbody"
	edited := "---\ndescription: \"LIN (หลิน) — PO\"\n---\nbody changed"

	if KitHash(repo, "agents/lin.md") == KitHash(edited, "agents/lin.md") {
		t.Fatal("kit body changes must change the drift hash")
	}
}

func TestMergePreservesHostModelPermission(t *testing.T) {
	repo := "---\ndescription: \"LIN (หลิน) — PO\"\nmode: subagent\nmodel: opencode/old\ncolor: \"#E53935\"\n---\nrepo body"
	host := "---\ndescription: \"LIN (หลิน) — PO\"\nmode: subagent\nmodel: opencode/host-tuned\npermission:\n  edit: allow\n  bash: read\n---\nhost body"

	merged, err := MergeAgentContent(repo, host)
	if err != nil {
		t.Fatal(err)
	}
	a, err := frontmatter.ParseAgent("agents/lin.md", merged)
	if err != nil {
		t.Fatal(err)
	}
	if a.ModelHint != "opencode/host-tuned" {
		t.Fatalf("model = %q, want host-tuned preserved", a.ModelHint)
	}
	if a.Description == "" || a.Color != "#E53935" {
		t.Fatalf("kit fields lost: %+v", a)
	}
	if a.Body != "repo body" {
		t.Fatalf("body = %q, want repo body", a.Body)
	}
	if !reflect.DeepEqual(hostPerm(t, merged), map[string]any{"edit": "allow", "bash": "read"}) {
		t.Fatal("permission block lost in merge")
	}
}

func hostPerm(t *testing.T, content string) map[string]any {
	t.Helper()
	fm, _, err := frontmatter.SplitFrontmatter(content)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := yaml.Unmarshal([]byte(fm), &m); err != nil {
		t.Fatal(err)
	}
	p, _ := m["permission"].(map[string]any)
	return p
}
