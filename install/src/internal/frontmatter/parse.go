// Package frontmatter parses agent identity files in either the legacy
// opencode-native format or the canonical warp v1 format (CLI_SPEC.md §2).
package frontmatter

import (
	"fmt"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Agent is a parsed agent identity file.
type Agent struct {
	Name        string // routing id = filename base
	RelPath     string // repo-relative path, e.g. agents/aiy.md
	Canonical   bool   // true if warp_version frontmatter is present
	WarpVersion int
	DisplayName string
	Role        string
	Color       string
	Department  string
	Rank        string
	ReportsTo   string
	ModelHint   string
	Description string
	Mode        string
	Skills      []string
	Body        string // persona content after the frontmatter block
}

// SplitFrontmatter returns the frontmatter YAML and body of a `---`-delimited
// markdown file. Files without a leading `---` yield empty frontmatter and the
// full content as body. The FIRST line that is exactly `---` (trailing spaces
// or tabs allowed; LF and CRLF both supported) terminates the block; a later
// `---` is body content, e.g. a horizontal rule.
//
// Malformed delimiters are errors, never silent truncation: a line starting
// with `---` that is not an exact closing marker (e.g. `----` or `--- x`)
// returns an error instead of dropping the body.
func SplitFrontmatter(content string) (fm, body string, err error) {
	var eol string
	switch {
	case strings.HasPrefix(content, "---\r\n"):
		eol = "\r\n"
	case strings.HasPrefix(content, "---\n"):
		eol = "\n"
	default:
		return "", content, nil
	}
	start := len("---") + len(eol)

	i := start
	for i < len(content) {
		lineEnd := len(content)
		if nl := strings.Index(content[i:], eol); nl >= 0 {
			lineEnd = i + nl
		}
		line := strings.TrimRight(content[i:lineEnd], " \t\r")
		if line == "---" {
			fm = content[start:i]
			if lineEnd == len(content) {
				return fm, "", nil
			}
			return fm, content[lineEnd+len(eol):], nil
		}
		if strings.HasPrefix(line, "---") {
			return "", "", fmt.Errorf("malformed frontmatter: line %q looks like a closing '---' but is not", line)
		}
		if lineEnd == len(content) {
			break
		}
		i = lineEnd + len(eol)
	}
	return "", "", fmt.Errorf("unterminated frontmatter: missing closing '---'")
}

// ParseAgent parses an agent file at repo-relative path relPath.
func ParseAgent(relPath, content string) (*Agent, error) {
	name := strings.TrimSuffix(filepath.Base(relPath), filepath.Ext(relPath))
	fm, body, err := SplitFrontmatter(content)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", relPath, err)
	}
	a := &Agent{Name: name, RelPath: relPath, Body: body}
	if fm == "" {
		return a, nil // body-only file — valid for non-agent payload
	}
	var m map[string]any
	if err := yaml.Unmarshal([]byte(fm), &m); err != nil {
		return nil, fmt.Errorf("%s: frontmatter YAML: %w", relPath, err)
	}
	if _, ok := m["warp_version"]; ok {
		a.Canonical = true
		parseCanonical(a, m)
	} else {
		parseLegacy(a, m)
	}
	return a, nil
}

func parseLegacy(a *Agent, m map[string]any) {
	a.Description = str(m["description"])
	a.Mode = str(m["mode"])
	a.ModelHint = str(m["model"])
	a.Color = str(m["color"])
	if a.Description != "" {
		if i := strings.Index(a.Description, "("); i >= 0 {
			if j := strings.Index(a.Description[i:], ")"); j > 0 {
				a.DisplayName = strings.TrimSpace(a.Description[i+1 : i+j])
			}
		}
		if k := strings.Index(a.Description, "—"); k >= 0 {
			a.Role = strings.TrimSpace(a.Description[k+3:])
		}
	}
}

func parseCanonical(a *Agent, m map[string]any) {
	a.WarpVersion = intVal(m["warp_version"])
	a.DisplayName = str(m["display_name"])
	a.Role = str(m["role"])
	a.Color = str(m["color"])
	a.Department = str(m["department"])
	a.Rank = str(m["rank"])
	a.ReportsTo = str(m["reports_to"])
	a.ModelHint = str(m["model_hint"])
	a.Description = str(m["description"])
	a.Mode = str(m["mode"])
	a.Skills = strSlice(m["skills"])
}

func str(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func intVal(v any) int {
	if i, ok := v.(int); ok {
		return i
	}
	return 0
}

func strSlice(v any) []string {
	switch t := v.(type) {
	case []any:
		var out []string
		for _, e := range t {
			if s, ok := e.(string); ok {
				out = append(out, s)
			}
		}
		return out
	case []string:
		return t
	}
	return nil
}
