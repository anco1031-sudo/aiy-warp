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
	// Telegram chat IDs (<15 digits) are P1 parameterization debt, not gate items.
	rep := Scan(map[string]string{"x": "chat 852106923"}, nil)
	if len(rep.Exports) != 0 {
		t.Fatalf("9-digit ID must not trip the export gate in P0: %+v", rep.Exports)
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
