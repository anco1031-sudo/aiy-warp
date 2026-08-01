package manifest

import (
	"strings"
	"testing"
)

func TestMergeWarnsOnUnparseableHostFM(t *testing.T) {
	repo := "---\nname: kwan\nmodel: opencode/repo-model\n---\nbody-repo"
	host := "no frontmatter at all"
	got, warns, err := MergeAgentContentV(repo, host)
	if err != nil {
		t.Fatal(err)
	}
	if got != repo {
		t.Fatal("repo content must win when host FM is unparseable")
	}
	if len(warns) != 1 || !strings.Contains(warns[0], "unparseable") {
		t.Fatalf("warns = %v, want unparseable-host warning (F12)", warns)
	}
}

func TestMergeWarnsOnInvalidHostYAML(t *testing.T) {
	repo := "---\nname: kwan\nmodel: opencode/repo-model\n---\nbody-repo"
	host := "---\nmodel: [1,2\n---\nbody-host"
	got, warns, err := MergeAgentContentV(repo, host)
	if err != nil {
		t.Fatal(err)
	}
	if got != repo {
		t.Fatal("repo content must win on invalid host YAML")
	}
	if len(warns) != 1 || !strings.Contains(warns[0], "YAML invalid") {
		t.Fatalf("warns = %v, want invalid-YAML warning (F12)", warns)
	}
}

func TestMergePreservesHostLocal(t *testing.T) {
	repo := "---\nname: kwan\nrole: Head Trader\nmodel: opencode/repo-model\n---\nbody-repo"
	host := "---\nname: kwan\nmodel: opencode/host-model\npermission:\n  edit: false\n---\nbody-host"
	got, warns, err := MergeAgentContentV(repo, host)
	if err != nil {
		t.Fatal(err)
	}
	if len(warns) != 0 {
		t.Fatalf("warns = %v, want none", warns)
	}
	if !strings.Contains(got, "opencode/host-model") {
		t.Errorf("host model must win:\n%s", got)
	}
	if !strings.Contains(got, "edit: false") {
		t.Errorf("host permission must be preserved:\n%s", got)
	}
	if !strings.Contains(got, "body-repo") {
		t.Errorf("repo body must win:\n%s", got)
	}
}

func TestMergeNoChangeReturnsRepo(t *testing.T) {
	repo := "---\nname: kwan\nmodel: opencode/x\n---\nbody-repo"
	host := "---\nname: kwan\nmodel: opencode/x\n---\nbody-host"
	got, _, err := MergeAgentContentV(repo, host)
	if err != nil {
		t.Fatal(err)
	}
	if got != repo {
		t.Fatalf("identical host-local values must return repo content unchanged")
	}
}
