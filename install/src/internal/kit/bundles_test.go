package kit

import "testing"

func TestResolveTeam(t *testing.T) {
	b := &Bundles{
		Departments: []string{"trading"},
		Agents: []Agent{
			{Name: "kwan", Rank: "head", Department: "trading"},
			{Name: "bee", Rank: "executor", Department: "trading"},
			{Name: "aiy", Rank: "head", Department: "core"},
		},
	}
	cases := map[string]string{
		"trading": "trading", // declared department → passthrough
		"kwan":    "trading", // head agent name → its department
		"bee":     "bee",     // executor name → unchanged (caller reports unknown)
		"bogus":   "bogus",   // unknown → unchanged
	}
	for in, want := range cases {
		if got := b.ResolveTeam(in); got != want {
			t.Errorf("ResolveTeam(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestTeamAndFind(t *testing.T) {
	b := &Bundles{
		Departments: []string{"core", "trading"},
		Agents: []Agent{
			{Name: "kwan", Rank: "head", Department: "trading"},
			{Name: "bee", Rank: "executor", Department: "trading"},
			{Name: "aiy", Rank: "head", Department: "core"},
		},
	}
	if got := b.Team("trading"); len(got) != 2 || got[0] != "bee" || got[1] != "kwan" {
		t.Fatalf("Team(trading) = %v", got)
	}
	if b.Find("kwan") == nil || b.Find("nope") != nil {
		t.Fatal("Find")
	}
	if b.HasDept("trading") != true || b.HasDept("nope") != false {
		t.Fatal("HasDept")
	}
}
