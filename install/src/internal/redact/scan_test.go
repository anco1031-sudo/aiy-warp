package redact

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestBlocksFakeSecrets(t *testing.T) {
	cases := map[string]string{
		"prefix":     "the api key is sk-test-abcdefghijklmnop123456789",
		"pem":        "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEA",
		"pgp-block":  "-----BEGIN PGP PRIVATE KEY BLOCK-----\nMIIEowIBAAKCAQEA",
		"jwt":        "eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U",
		"assignment": "token: supersecretvalue1234567890",
		"ghp":        "ghp_abcdefghijklmnopqrstuvwxyz123456",
		"ak":         "AKIAIOSFODNN7EXAMPLE",
	}
	for name, content := range cases {
		rep := Scan(map[string]string{"agents/x.md": content}, nil)
		if !rep.HasBlocks() {
			t.Errorf("%s: expected a hard block, got none", name)
		}
	}
}

func TestPlaceholderDoesNotBlock(t *testing.T) {
	content := "token: your-api-key-here\nsecret: changeme\npassword: example"
	rep := Scan(map[string]string{"x": content}, nil)
	if rep.HasBlocks() {
		t.Fatalf("placeholders must not block: %+v", rep.Blocks)
	}
}

func TestProseDoesNotBlock(t *testing.T) {
	content := "Bot Token (embedded in script)\nToken optimization saves ~76% of tokens\nreason: Token optimization"
	rep := Scan(map[string]string{"x": content}, nil)
	if rep.HasBlocks() {
		t.Fatalf("prose must not hard-block: %+v", rep.Blocks)
	}
	if len(rep.Warns) == 0 {
		t.Fatal("expected a warn for 'Bot Token' reference")
	}
}

func TestLongIDExportBlockNotHardBlock(t *testing.T) {
	rep := Scan(map[string]string{"x": "channel 1527698229347487904 is the todo list"}, nil)
	if len(rep.Exports) != 1 {
		t.Fatalf("expected 1 export finding, got %d", len(rep.Exports))
	}
	if rep.HasBlocks() {
		t.Fatal("numeric IDs must not hard-block install")
	}
}

func TestAllowIDsExempts(t *testing.T) {
	allow := map[string]bool{"1527698229347487904": true}
	rep := Scan(map[string]string{"x": "channel 1527698229347487904"}, allow)
	if len(rep.Exports) != 0 {
		t.Fatalf("allowlisted ID must be exempt, got %d", len(rep.Exports))
	}
}

func TestShortNumericIDIgnored(t *testing.T) {
	// Numbers shorter than 9 digits are not identifiers and must not trip the
	// export gate. 9-13 digit numbers are Telegram-class identifiers by design
	// (DESIGN-NOTE §2, fail-closed: phones included; --allow-identifiers is the
	// override).
	rep := Scan(map[string]string{"x": "date 20260801 counter 123456"}, nil)
	if len(rep.Exports) != 0 {
		t.Fatalf("short numbers must not trip the export gate: %+v", rep.Exports)
	}
}

func TestTelegramIDExportBlock(t *testing.T) {
	// DESIGN-NOTE §2: Telegram identifiers are export-gated (9-13 digits,
	// optional -100 prefix) — unlike short numbers they are host-routing IDs.
	cases := map[string]string{
		"plain chat":    "chat 852106923",
		"supergroup":    "supergroup -1001234567890",
		"in flags":      "--chat_id 852106923",
		"negative chat": "chat -123456789",
		"backtick":      "`852106923`",
		"adjacent word": "id=852106923",
	}
	for name, content := range cases {
		rep := Scan(map[string]string{"x": content}, nil)
		if len(rep.Exports) != 1 {
			t.Errorf("%s (%q): expected 1 export finding, got %d (%+v)", name, content, len(rep.Exports), rep.Exports)
		}
		if rep.HasBlocks() {
			t.Errorf("%s: Telegram IDs must not hard-block install", name)
		}
	}
	// Allowlist exemption must work for the raw ID (plain and -100 prefixed).
	for _, id := range []string{"852106923", "-1001234567890"} {
		rep := Scan(map[string]string{"x": "chat " + id}, map[string]bool{id: true})
		if len(rep.Exports) != 0 {
			t.Fatalf("allowlisted Telegram ID %s must be exempt: %+v", id, rep.Exports)
		}
	}
}

func TestLongIDNotPartiallyMatchedAsTelegram(t *testing.T) {
	// A 15+ digit number is a long-numeric-id finding, never a partial
	// telegram match — must not double-report or mis-classify.
	rep := Scan(map[string]string{"x": "snowflake 1527698229347487904"}, nil)
	if len(rep.Exports) != 1 || rep.Exports[0].Pattern != "long-numeric-id" {
		t.Fatalf("expected 1 long-numeric-id finding, got %+v", rep.Exports)
	}
}

func TestAssignmentInURLIsWarnNotBlock(t *testing.T) {
	content := "See https://example.com/auth?api_key=abcdefghijklmnop123456 for details"
	rep := Scan(map[string]string{"docs/guide.md": content}, nil)
	if rep.HasBlocks() {
		t.Fatalf("query-string assignment must not hard-block: %+v", rep.Blocks)
	}
	found := false
	for _, w := range rep.Warns {
		if w.Pattern == "secret-assignment-in-url" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected a secret-assignment-in-url warn, got %+v", rep.Warns)
	}
}

// TestKitContentProducesNoBlocks scans the real kit — the P0 exit criterion
// depends on install/doctor never false-positive on today's payload.
func TestKitContentProducesNoBlocks(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(file), "..", "..", "..", "..")
	glob := []string{
		"agents/*.md",
		"skills/*/SKILL.md",
		"templates/*.md",
		"playbooks/*.md",
	}
	files := map[string]string{}
	for _, pat := range glob {
		matches, err := filepath.Glob(filepath.Join(root, filepath.FromSlash(pat)))
		if err != nil {
			t.Fatal(err)
		}
		for _, m := range matches {
			b, err := os.ReadFile(m)
			if err != nil {
				t.Fatal(err)
			}
			rel, _ := filepath.Rel(root, m)
			files[filepath.ToSlash(rel)] = string(b)
		}
	}
	rep := Scan(files, nil)
	if rep.HasBlocks() {
		t.Fatalf("kit must contain zero credential values, got %d: %+v", len(rep.Blocks), rep.Blocks)
	}
	// The kit does carry identifiers + host paths today — expected warnings only.
	t.Logf("kit: %d export-gate findings, %d warnings (expected; P1 parameterizes)", len(rep.Exports), len(rep.Warns))
}
