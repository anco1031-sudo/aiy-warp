// Package condense implements the condensed persona engine (CLI_SPEC.md §3).
// It extracts a platform-portable persona from agent markdown: keeps
// personality/directives/boundaries/reporting, drops opencode-only syntax and
// host paths, and caps the result (≤ ~1500 words, spec §3.1). Teams collapse
// into a single conductor persona (§3.2) via Collapse.
package condense

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/anco1031-sudo/aiy-warp/install/src/internal/frontmatter"
)

// WordCap is the total rendered-persona budget (spec §3.1).
const WordCap = 1500

// personalityWords is the budget for the personality block.
const personalityWords = 150

// Member is one executor in a collapsed team.
type Member struct {
	Name        string
	DisplayName string
	Role        string
	OneLiner    string
}

// Persona is the condensed, platform-portable representation of an agent or a
// collapsed team. Renderers (internal/render) format it per platform.
type Persona struct {
	Name         string
	DisplayName  string // full form, e.g. "Kwan (ขวัญ)"
	ThaiName     string // the non-latin display part, e.g. "ขวัญ" ("" when absent)
	Role         string
	Color        string
	Dept         string
	Rank         string
	ReportsTo    string
	Personality  string
	Directives   []string
	Boundaries   []string
	Team         []Member // collapse only
	Pipeline     string   // collapse only
	Routing      []string // collapse only (routing-table rows)
	OutputFormat string
	Sources      []string // source agent relpaths (collapse only)
}

var (
	sectionRE = regexp.MustCompile(`^([0-9]+)\.\s+(.+?)\s*$`)
	bulletRE  = regexp.MustCompile(`^\s*[-*]\s+(.+?)\s*$`)
	numItemRE = regexp.MustCompile(`^[0-9]+\.\s+(.+?)\s*$`)
	parenRE   = regexp.MustCompile(`\([^)]*\)`)
)

// section is one numbered body section.
type section struct {
	title string
	body  string
}

// splitSections splits a persona body into numbered sections. A numbered line
// is a section header when its title is predominantly UPPERCASE (after
// stripping parentheticals) and short — numbered list items ("1. Direct
// Response to Lin: …") are sentence-case and long, so they stay in the body.
func splitSections(body string) []section {
	var secs []section
	var cur *section
	flush := func() {
		if cur != nil {
			cur.body = strings.TrimRight(cur.body, "\n")
			secs = append(secs, *cur)
		}
	}
	for _, ln := range strings.Split(body, "\n") {
		if m := sectionRE.FindStringSubmatch(ln); m != nil && isHeader(m[2]) {
			flush()
			cur = &section{title: m[2]}
			continue
		}
		if cur != nil {
			cur.body += ln + "\n"
		}
	}
	flush()
	return secs
}

// isHeader reports whether a numbered line title is a section header (not a
// list item). Headers are short and mostly uppercase; parentheticals such as
// "(Delegation via task tool)" are ignored for the check.
func isHeader(title string) bool {
	if len(title) > 80 {
		return false
	}
	t := parenRE.ReplaceAllString(title, "")
	var upper, lower int
	for _, r := range t {
		if unicode.IsUpper(r) {
			upper++
		} else if unicode.IsLower(r) {
			lower++
		}
	}
	if upper+lower == 0 {
		return false
	}
	return float64(upper)/float64(upper+lower) >= 0.6
}

// keepSection classifies a section title as persona-relevant.
func keepSection(title string) bool {
	t := strings.ToUpper(title)
	for _, bad := range []string{"WORKSPACE", "KNOWLEDGE", "SHARED RULES", "PERMISSION"} {
		if strings.Contains(t, bad) {
			return false
		}
	}
	if strings.Contains(t, "SUBAGENT_TYPE") || strings.Contains(t, "WHEN USING TASK") {
		return false
	}
	for _, good := range []string{"ROLE", "MINDSET", "PERSONALITY", "TONE", "VOICE",
		"RELATIONSHIP", "WORKFLOW", "INTERACTION", "PROTOCOL", "DIRECTIVE",
		"RESPONSIBILIT", "REPORTING", "BOUNDAR", "LANE", "CHARTER", "TEAM MATRIX", "THE TEAM"} {
		if strings.Contains(t, good) {
			return true
		}
	}
	return false
}

// noiseRE matches opencode-only syntax and host paths stripped from condensed
// output (spec §3.1: drop permission/model/workspace paths/bash examples).
var noiseRE = regexp.MustCompile(`(?i)(/home/|/Users/|~/|myObsidian|02-Areas|01-Projects|` +
	`Aiy_Workspace|subagent_type|@mention|permission:|model:|mode:|color:|\btask tool\b|` +
	`(?:^|[[:space:]])@[a-z][a-z-]*(?:[[:space:]]|$))`)

// isNoiseLine reports whether a body line should be dropped.
func isNoiseLine(ln string) bool {
	s := strings.TrimSpace(ln)
	if s == "" {
		return false
	}
	// Drop host paths, opencode syntax, and workspace references.
	if noiseRE.MatchString(s) {
		return true
	}
	// Drop bash command examples.
	if strings.HasPrefix(s, "cd ") || strings.HasPrefix(s, "./") ||
		strings.HasPrefix(s, "python3 ") || strings.HasPrefix(s, "kill ") ||
		strings.HasPrefix(s, "pgrep ") || strings.HasPrefix(s, "source ") {
		return true
	}
	return false
}

// stripFences removes fenced code blocks whose content is opencode task-tool
// templates (they carry subagent_type), leaving other fences intact.
func stripFences(s string) string {
	var out []string
	inFence := false
	var block []string
	for _, ln := range strings.Split(s, "\n") {
		if strings.HasPrefix(strings.TrimSpace(ln), "```") {
			if inFence {
				if !strings.Contains(strings.Join(block, "\n"), "subagent_type") {
					out = append(out, block...)
				}
				block = nil
				inFence = false
				continue
			}
			inFence = true
			continue
		}
		if inFence {
			block = append(block, ln)
		} else {
			out = append(out, ln)
		}
	}
	return strings.Join(out, "\n")
}

// cleanBody drops noise lines and opencode task-tool fences from a section body.
func cleanBody(body string) string {
	body = stripFences(body)
	var out []string
	for _, ln := range strings.Split(body, "\n") {
		if isNoiseLine(ln) {
			continue
		}
		out = append(out, ln)
	}
	return strings.TrimSpace(strings.Join(out, "\n"))
}

// bullets extracts dash bullets and numbered list items from cleaned text.
func bullets(text string, cap int) []string {
	var out []string
	for _, ln := range strings.Split(text, "\n") {
		s := strings.TrimSpace(ln)
		if s == "" {
			continue
		}
		if m := bulletRE.FindStringSubmatch(s); m != nil {
			out = append(out, strings.TrimSpace(m[1]))
		} else if m := numItemRE.FindStringSubmatch(s); m != nil {
			out = append(out, strings.TrimSpace(m[1]))
		}
		if cap > 0 && len(out) >= cap {
			break
		}
	}
	return out
}

// paragraphs returns the first n non-empty paragraph strings.
func paragraphs(text string, n int) []string {
	var out []string
	for _, p := range strings.Split(text, "\n\n") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
			if len(out) >= n {
				break
			}
		}
	}
	return out
}

// WordCount counts whitespace-separated words.
func WordCount(s string) int {
	if strings.TrimSpace(s) == "" {
		return 0
	}
	return len(strings.Fields(s))
}

// capWords truncates s to at most n words without splitting a word.
func capWords(s string, n int) string {
	f := strings.Fields(s)
	if len(f) <= n {
		return s
	}
	return strings.Join(f[:n], " ") + " …"
}

// FullName renders the display name in its canonical "Name (ไทย)" form, e.g.
// "Kwan (ขวัญ)". When the agent has no distinct display name the routing id is
// returned unchanged. A display that already carries the routing id (e.g. the
// migrate-derived canonical "ขวัญ (Kwan)") is returned as-is so it is never
// double-wrapped.
func FullName(name, display string) string {
	if display == "" || strings.EqualFold(display, name) {
		return name
	}
	if strings.Contains(strings.ToLower(display), strings.ToLower(name)) {
		return display
	}
	return titleCase(name) + " (" + display + ")"
}

// titleCase capitalizes the first rune of s.
func titleCase(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// thaiPart returns the non-latin display part of a display name, stripping a
// canonical "ไทย (Name)" parenthetical when present.
func thaiPart(display string) string {
	if i := strings.Index(display, " ("); i > 0 {
		return display[:i]
	}
	return display
}

// CompactRole shortens a role sentence to its first clause and drops a leading
// "The ", e.g. "The Head Trader & Strategy Commander. Aiy delegates…" →
// "Head Trader & Strategy Commander". Used for member lines and routing triggers.
func CompactRole(role string) string {
	role = strings.TrimSpace(role)
	if i := strings.Index(role, "."); i > 0 {
		role = role[:i]
	}
	return strings.TrimSpace(strings.TrimPrefix(role, "The "))
}

// FirstSentence returns the first sentence of s (up to the first period),
// trimmed. Keeps a leading "The " — used for renderer title lines.
func FirstSentence(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.Index(s, "."); i > 0 {
		s = s[:i]
	}
	return strings.TrimSpace(s)
}

// HeaderName renders the title form "AIY (อัย)" — uppercased routing id plus
// the thai display in parens — used by renderer title lines. Falls back to the
// uppercased id when no thai name exists.
func HeaderName(p *Persona) string {
	if p.ThaiName != "" {
		return strings.ToUpper(p.Name) + " (" + p.ThaiName + ")"
	}
	return strings.ToUpper(p.Name)
}

// Extract builds a Persona from a parsed agent. collapsed=false yields a
// single-agent persona; collapsed=true also carries the team roster (members,
// routing) for the conductor renderers.
func Extract(a *frontmatter.Agent, collapsed bool) *Persona {
	p := &Persona{
		Name:        a.Name,
		DisplayName: FullName(a.Name, a.DisplayName),
		ThaiName:    thaiPart(a.DisplayName),
		Role:        a.Role,
		Color:       a.Color,
		Dept:        a.Department,
		Rank:        a.Rank,
		ReportsTo:   a.ReportsTo,
	}
	secs := splitSections(a.Body)
	for _, s := range secs {
		if !keepSection(s.title) {
			continue
		}
		t := strings.ToUpper(s.title)
		clean := cleanBody(s.body)
		switch {
		case strings.Contains(t, "REPORTING"):
			p.OutputFormat = strings.TrimSpace(clean)
		case strings.Contains(t, "BOUNDAR") || strings.Contains(t, "LANE") || strings.Contains(t, "CHARTER"):
			p.Boundaries = append(p.Boundaries, bullets(clean, 4)...)
		case strings.Contains(t, "ROLE") || strings.Contains(t, "MINDSET") ||
			strings.Contains(t, "PERSONALITY") || strings.Contains(t, "TONE") ||
			strings.Contains(t, "VOICE") || strings.Contains(t, "RELATIONSHIP"):
			if p.Personality == "" {
				p.Personality = paragraphs(clean, 2)[0]
			}
		case strings.Contains(t, "TEAM MATRIX") || strings.Contains(t, "THE TEAM"):
			// roster handled by collapse; routing rows extracted below
			if collapsed {
				p.Routing = append(p.Routing, routingRows(clean)...)
			}
		default: // WORKFLOW / INTERACTION / PROTOCOL / DIRECTIVE / RESPONSIBILIT
			p.Directives = append(p.Directives, bullets(clean, 6)...)
		}
	}
	p.Personality = capWords(cleanPersonality(p.Personality), personalityWords)
	return p
}

// cleanPersonality normalizes the personality block (collapse internal
// newlines into spaces for a tight, pastable paragraph).
func cleanPersonality(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// routingRows extracts "…→ member" routing rows from a team-matrix section,
// e.g. "- Fon (ฝน) - The News & Sentiment Analyst: … Delegate news scans,
// event analysis, and sentiment tracking to Fon." → "news scans, event
// analysis, and sentiment tracking → Fon".
func routingRows(text string) []string {
	var out []string
	for _, ln := range strings.Split(text, "\n") {
		s := strings.TrimSpace(ln)
		m := bulletRE.FindStringSubmatch(s)
		if m == nil {
			continue
		}
		body := m[1]
		low := strings.ToLower(body)
		for _, kw := range []string{"delegate ", "assign ", "send ", "route "} {
			if i := strings.Index(low, kw); i >= 0 {
				rest := body[i+len(kw):]
				if j := strings.Index(rest, " to "); j > 0 {
					target := strings.TrimSpace(strings.TrimSuffix(rest[j+4:], "."))
					what := strings.TrimSpace(rest[:j])
					if what != "" && target != "" {
						out = append(out, what+" → "+target)
					}
				}
				break
			}
		}
	}
	return out
}
