// Package redact implements the no-secrets gate (CLI_SPEC.md §5.4).
//
// Per DESIGN-NOTE.md §2, the gate detects credential VALUES, not
// documentation words: hard-block on high-entropy credentials (exit 5),
// export-only block on long numeric identifiers (overridable via
// --allow-identifiers), and warnings for prose references and host paths.
package redact

import (
	"regexp"
	"sort"
	"strings"
	"unicode"
)

// Kind classifies a finding's severity.
type Kind string

const (
	// KindBlock is a credential value — abort install/export, exit 5.
	KindBlock Kind = "block"
	// KindExport blocks export unless the identifier is allowlisted.
	KindExport Kind = "export"
	// KindWarn is informational only — never blocks.
	KindWarn Kind = "warn"
)

// Finding is one suspicious match within a file.
type Finding struct {
	File    string `json:"file"`
	Line    int    `json:"line"`
	Kind    Kind   `json:"kind"`
	Pattern string `json:"pattern"`
	Match   string `json:"match"`
}

// Report aggregates findings by severity.
type Report struct {
	Blocks  []Finding `json:"blocks"`  // credential values — hard fail
	Exports []Finding `json:"exports"` // identifiers — block export unless allowlisted
	Warns   []Finding `json:"warns"`   // prose references / host paths
}

// HasBlocks reports whether any credential value was found.
func (r Report) HasBlocks() bool { return len(r.Blocks) > 0 }

// BlockFiles returns the unique files containing credential values, sorted.
func (r Report) BlockFiles() []string { return uniqueFiles(r.Blocks) }

// ExportBlockFiles returns the unique files with non-allowlisted identifiers.
func (r Report) ExportBlockFiles() []string { return uniqueFiles(r.Exports) }

func uniqueFiles(fs []Finding) []string {
	seen := map[string]bool{}
	var out []string
	for _, f := range fs {
		if !seen[f.File] {
			seen[f.File] = true
			out = append(out, f.File)
		}
	}
	sort.Strings(out)
	return out
}

var (
	// Hard-block: PEM/OpenSSH private key blocks. The optional " BLOCK" suffix
	// covers armored PGP keys ("-----BEGIN PGP PRIVATE KEY BLOCK-----").
	pemRE = regexp.MustCompile(`-----BEGIN [A-Z0-9 ]*(PRIVATE|EC|RSA|OPENSSH|DSA|PGP)[ A-Z0-9]*KEY(?: BLOCK)?-----`)

	// Hard-block: well-known credential prefixes.
	prefixRE = regexp.MustCompile(`(?i)\b(sk-[A-Za-z0-9_-]{8,}|sk-ant-[A-Za-z0-9_-]{16,}|ghp_[A-Za-z0-9]{20,}|gho_[A-Za-z0-9]{20,}|xox[baprs]-[A-Za-z0-9-]{10,}|glpat-[A-Za-z0-9_-]{16,}|rk_live_[A-Za-z0-9]{16,}|ya29\.[A-Za-z0-9_-]{20,}|AIza[0-9A-Za-z_-]{30,}|AKIA[0-9A-Z]{16})\b`)

	// Hard-block: JWT-shaped blobs.
	jwtRE = regexp.MustCompile(`\beyJ[A-Za-z0-9_-]{12,}\.[A-Za-z0-9_-]{12,}\.[A-Za-z0-9_-]{12,}\b`)

	// Hard-block: secret assignments with an opaque value (>=16 chars).
	assignRE = regexp.MustCompile(`(?i)\b(api[_-]?key|apikey|client[_-]?secret|access[_-]?token|auth[_-]?token|secret|password|token)\b\s*[:=]\s*['"]?([A-Za-z0-9_\-/+=.]{16,})['"]?`)

	// Export-only: long numeric identifiers (Discord snowflakes, etc.).
	longIDRE = regexp.MustCompile(`\b\d{15,}\b`)

	// Export-only: Telegram chat/supergroup identifiers — 9-13 digits with an
	// optional -100 prefix (DESIGN-NOTE §2 contract). Boundary chars are
	// consumed so a longer number is never partially matched; the ID itself is
	// captured in group 1 for allowlist comparison.
	telegramRE = regexp.MustCompile(`(?:^|[^0-9])((?:-100)?\d{9,13})(?:[^0-9]|$)`)

	// Warn: prose references to credential concepts.
	refRE = regexp.MustCompile(`(?i)\b(api[_-]?key|apikey|client[_-]?secret|password|(?:bot|api|access|auth|page|user|discord|telegram|facebook)[- ]?token)\b`)

	// Warn: host paths leaking into the kit — cross-platform (F10: ~/, /Users/,
	// C:\, /home/).
	hostPathRE = regexp.MustCompile(`(?:/home/[A-Za-z0-9_./-]+|~/[A-Za-z0-9_./-]+|/Users/[A-Za-z0-9_./-]+|[A-Za-z]:\\[^\s]+)`)
)

// placeholderValues are documentation placeholders that must never hard-block.
var placeholderValues = map[string]bool{
	"changeme": true, "your-api-key": true, "your-api-key-here": true,
	"your_token_here": true, "your-token-here": true, "password123": true,
	"example": true, "dummy": true, "redacted": true, "secret": true,
	"token": true, "true": true, "false": true, "null": true, "none": true,
}

// looksOpaque reports whether v is high-entropy enough to be a real credential:
// at least two character classes and not a single repeated character.
func looksOpaque(v string) bool {
	if placeholderValues[strings.ToLower(v)] {
		return false
	}
	var lower, upper, digit, other bool
	for _, r := range v {
		switch {
		case unicode.IsLower(r):
			lower = true
		case unicode.IsUpper(r):
			upper = true
		case unicode.IsDigit(r):
			digit = true
		default:
			other = true
		}
	}
	classes := 0
	for _, b := range []bool{lower, upper, digit, other} {
		if b {
			classes++
		}
	}
	allSame := true
	for i := 1; i < len(v); i++ {
		if v[i] != v[0] {
			allSame = false
			break
		}
	}
	return classes >= 2 && !allSame
}

// assignValue returns the captured assignment value if it looks like a real
// credential, otherwise "" and the match start (or -1 when not applicable).
func assignValue(line string) (string, int) {
	m := assignRE.FindStringSubmatchIndex(line)
	if len(m) != 6 {
		return "", -1
	}
	val := line[m[4]:m[5]]
	if looksOpaque(val) {
		return val, m[0]
	}
	return "", -1
}

// isQueryContext reports whether an assignment at byte offset start sits in a
// URL query string (e.g. "https://example.com/auth?api_key=..."). Such values
// are documentation references, not credentials.
func isQueryContext(line string, start int) bool {
	prefix := strings.TrimRight(line[:start], " \t")
	return strings.Contains(prefix, "://") ||
		strings.HasSuffix(prefix, "?") || strings.HasSuffix(prefix, "&")
}

// Scan inspects file contents (relpath → content) and returns a Report.
// allow is the set of explicitly-declared public identifiers exempted from
// the export gate.
func Scan(files map[string]string, allow map[string]bool) Report {
	var rep Report
	paths := make([]string, 0, len(files))
	for p := range files {
		paths = append(paths, p)
	}
	sort.Strings(paths)

	for _, path := range paths {
		lines := strings.Split(files[path], "\n")
		for i, line := range lines {
			n := i + 1
			if m := pemRE.FindString(line); m != "" {
				rep.Blocks = append(rep.Blocks, Finding{File: path, Line: n, Kind: KindBlock, Pattern: "pem", Match: clip(m)})
			}
			if m := prefixRE.FindString(line); m != "" {
				rep.Blocks = append(rep.Blocks, Finding{File: path, Line: n, Kind: KindBlock, Pattern: "credential-prefix", Match: clip(m)})
			}
			if m := jwtRE.FindString(line); m != "" {
				rep.Blocks = append(rep.Blocks, Finding{File: path, Line: n, Kind: KindBlock, Pattern: "jwt", Match: clip(m)})
			}
			if val, start := assignValue(line); val != "" {
				f := Finding{File: path, Line: n, Kind: KindBlock, Pattern: "secret-assignment", Match: clip(val)}
				if isQueryContext(line, start) {
					f.Kind = KindWarn
					f.Pattern = "secret-assignment-in-url"
					rep.Warns = append(rep.Warns, f)
				} else {
					rep.Blocks = append(rep.Blocks, f)
				}
			}
			if m := longIDRE.FindString(line); m != "" && !allow[m] {
				rep.Exports = append(rep.Exports, Finding{File: path, Line: n, Kind: KindExport, Pattern: "long-numeric-id", Match: m})
			}
			if sm := telegramRE.FindStringSubmatch(line); sm != nil && !allow[sm[1]] {
				rep.Exports = append(rep.Exports, Finding{File: path, Line: n, Kind: KindExport, Pattern: "telegram-id", Match: sm[1]})
			}
			if m := refRE.FindString(line); m != "" {
				rep.Warns = append(rep.Warns, Finding{File: path, Line: n, Kind: KindWarn, Pattern: "credential-reference", Match: clip(m)})
			}
			if m := hostPathRE.FindString(line); m != "" {
				rep.Warns = append(rep.Warns, Finding{File: path, Line: n, Kind: KindWarn, Pattern: "host-path", Match: clip(m)})
			}
		}
	}
	return rep
}

// clip truncates s to 40 runes, never splitting a multibyte character (F11).
func clip(s string) string {
	r := []rune(s)
	if len(r) > 40 {
		return string(r[:40]) + "…"
	}
	return s
}
