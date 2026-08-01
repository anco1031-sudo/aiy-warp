package redact

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestHostPathWarnsCrossPlatform(t *testing.T) {
	lines := []string{
		"ref: /home/lu5her/01-Projects/aiy-warp",
		"ref: ~/myObsidian/Aiy_Workspace",
		"ref: /Users/joe/secret",
		`ref: C:\Users\joe\secret`,
	}
	for _, line := range lines {
		rep := Scan(map[string]string{"x.md": line}, nil)
		if rep.HasBlocks() {
			t.Fatalf("host path must warn, not block: %q (blocks=%+v)", line, rep.Blocks)
		}
		var found bool
		for _, w := range rep.Warns {
			if w.Pattern == "host-path" {
				found = true
			}
		}
		if !found {
			t.Errorf("host path %q must produce a host-path warning (F10)", line)
		}
	}
}

func TestClipRuneSafe(t *testing.T) {
	long := strings.Repeat("กขคง", 20) // 80 runes, multibyte
	got := clip(long)
	if !utf8.ValidString(got) {
		t.Fatal("clip must never split a multibyte rune (F11)")
	}
	r := []rune(got)
	if len(r) > 41 {
		t.Fatalf("clip produced %d runes, want ≤ 41", len(r))
	}
	if !strings.HasSuffix(got, "…") {
		t.Errorf("truncated clip must end with ellipsis, got %q", got)
	}
	if clip("abc") != "abc" {
		t.Error("short clip must pass through unchanged")
	}
}
