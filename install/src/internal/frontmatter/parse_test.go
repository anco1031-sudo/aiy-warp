package frontmatter

import (
	"strings"
	"testing"
)

func TestSplitAndLegacyParse(t *testing.T) {
	content := "---\ndescription: \"LIN (หลิน) — The Elite Product Owner\"\nmode: subagent\nmodel: opencode/deepseek-v4-flash-free\ncolor: \"#E53935\"\n---\n\n# Body"
	fm, body, err := SplitFrontmatter(content)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(fm, "mode: subagent") {
		t.Fatalf("unexpected frontmatter: %q", fm)
	}
	if body != "\n# Body" {
		t.Fatalf("unexpected body: %q", body)
	}

	a, err := ParseAgent("agents/lin.md", content)
	if err != nil {
		t.Fatal(err)
	}
	if a.Name != "lin" {
		t.Fatalf("name = %q, want lin", a.Name)
	}
	if a.Canonical {
		t.Fatal("legacy file must not be canonical")
	}
	if a.Mode != "subagent" || a.ModelHint != "opencode/deepseek-v4-flash-free" || a.Color != "#E53935" {
		t.Fatalf("legacy fields wrong: %+v", a)
	}
	if a.DisplayName != "หลิน" {
		t.Fatalf("display name = %q", a.DisplayName)
	}
	if !strings.HasPrefix(a.Role, "The Elite Product Owner") {
		t.Fatalf("role = %q", a.Role)
	}
}

func TestCanonicalParse(t *testing.T) {
	content := "---\nwarp_version: 1\nname: lin\ndisplay_name: \"หลิน (Lin)\"\nrole: PO\ndepartment: tech\nrank: head\nreports_to: aiy\nmodel_hint: opencode/x\nskills: [obsidian, aiy-messaging]\n---\nbody"
	a, err := ParseAgent("agents/lin.md", content)
	if err != nil {
		t.Fatal(err)
	}
	if !a.Canonical || a.WarpVersion != 1 {
		t.Fatalf("canonical flags wrong: %+v", a)
	}
	if a.Department != "tech" || a.Rank != "head" || a.ReportsTo != "aiy" {
		t.Fatalf("org fields wrong: %+v", a)
	}
	if len(a.Skills) != 2 || a.Skills[0] != "obsidian" {
		t.Fatalf("skills wrong: %v", a.Skills)
	}
}

func TestNoFrontmatterIsBodyOnly(t *testing.T) {
	content := "just a body with no delimiters"
	fm, body, err := SplitFrontmatter(content)
	if err != nil {
		t.Fatal(err)
	}
	if fm != "" || body != content {
		t.Fatalf("fm=%q body=%q", fm, body)
	}
}

func TestUnterminatedFrontmatter(t *testing.T) {
	if _, _, err := SplitFrontmatter("---\nname: x"); err == nil {
		t.Fatal("expected error for unterminated frontmatter")
	}
}

func TestLaterDelimiterIsBody(t *testing.T) {
	content := "---\nname: obsidian\ndescription: x\n---\n# Title\n\n---\n## Section\nrest"
	fm, body, err := SplitFrontmatter(content)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(fm, "name: obsidian") {
		t.Fatalf("fm = %q", fm)
	}
	if !strings.HasPrefix(body, "# Title") || !strings.Contains(body, "## Section") {
		t.Fatalf("body must keep later delimiters: %q", body)
	}
}
