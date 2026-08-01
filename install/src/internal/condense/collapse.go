package condense

import (
	"regexp"
	"sort"
	"strings"

	"github.com/anco1031-sudo/aiy-warp/install/src/internal/frontmatter"
)

// Collapse builds a single conductor Persona from a department head + its
// executors (CLI_SPEC.md §3.2): the head's directives/boundaries form the
// frame; each executor becomes a one-line internal persona; the department's
// pipeline skill becomes the workflow section.
func Collapse(head *frontmatter.Agent, members []*frontmatter.Agent, pipelineSkill string) *Persona {
	p := Extract(head, true)
	p.Team = make([]Member, 0, len(members))
	for _, m := range members {
		p.Team = append(p.Team, Member{
			Name:        m.Name,
			DisplayName: FullName(m.Name, m.DisplayName),
			Role:        CompactRole(m.Role),
			OneLiner:    memberOneLiner(m),
		})
	}
	p.Pipeline = pipeline(pipelineSkill, p.Team)
	p.Team = orderByPipeline(p.Team, p.Pipeline)
	if p.Routing == nil {
		p.Routing = fallbackRouting(p.Team)
	}
	// Ensure the output-format block exists; fall back to the head's reporting
	// bullet list.
	if strings.TrimSpace(p.OutputFormat) == "" {
		p.OutputFormat = reportingFallback(p)
	}
	return p
}

// orderByPipeline sorts team members by first mention in the pipeline text so
// the conductor lists analysts in pipeline order (e.g. Fon → June → Bee → Nam).
// Matching is case-insensitive: pipeline text names "Fon"/"June" while member
// routing ids are lowercase.
func orderByPipeline(team []Member, pipeline string) []Member {
	if pipeline == "" {
		return team
	}
	lower := strings.ToLower(pipeline)
	pos := func(name string) int {
		if i := strings.Index(lower, name); i >= 0 {
			return i
		}
		return 1 << 30
	}
	out := append([]Member{}, team...)
	sort.SliceStable(out, func(i, j int) bool { return pos(out[i].Name) < pos(out[j].Name) })
	return out
}

var sentenceRE = regexp.MustCompile(`([^.!?\n]+[.!?])\s*`)

// memberOneLiner distills an executor's personality into one sentence. A
// personality/voice/tone section is preferred; role & mindset is only a
// fallback (the "You are X, …" framing is dropped either way).
func memberOneLiner(m *frontmatter.Agent) string {
	secs := splitSections(m.Body)
	var pick *section
	for _, s := range secs {
		t := strings.ToUpper(s.title)
		if strings.Contains(t, "PERSONALITY") || strings.Contains(t, "VOICE") ||
			strings.Contains(t, "TONE") {
			pick = &s
			break
		}
	}
	if pick == nil {
		for _, s := range secs {
			t := strings.ToUpper(s.title)
			if strings.Contains(t, "MINDSET") {
				pick = &s
				break
			}
		}
	}
	if pick != nil {
		clean := cleanBody(pick.body)
		if p := paragraphs(clean, 1); len(p) > 0 {
			// Strip a "The X: " label prefix, e.g. "The Intuitive Analyst: ".
			one := p[0]
			if i := strings.Index(one, ": "); i > 0 && i < 60 {
				one = one[i+2:]
			}
			// Drop a "You are …" framing sentence.
			one = strings.TrimSpace(strings.TrimPrefix(one, "You are "))
			one = strings.TrimSpace(strings.TrimPrefix(one, "You're "))
			if s := sentenceRE.FindString(one); s != "" {
				return titleCase(strings.TrimSpace(s))
			}
			return titleCase(strings.TrimSpace(capWords(one, 24)))
		}
	}
	if m.Role != "" {
		return CompactRole(m.Role) + "."
	}
	return m.Name
}

// pipeline extracts the department workflow from the pipeline skill, or
// synthesizes one from the team roster.
func pipeline(skill string, team []Member) string {
	if skill != "" {
		if p := skillPipeline(skill); p != "" {
			return p
		}
	}
	var names []string
	for _, m := range team {
		names = append(names, m.Name)
	}
	if len(names) == 0 {
		return ""
	}
	return strings.Join(names, " → ") + " → " + "head synthesis"
}

var pipelineFlowRE = regexp.MustCompile(`(?s)## The Pipeline Flow(.*?)(?:\n## |\z)`)

// skillPipeline pulls the pipeline-flow block (code fence) from a skill file.
func skillPipeline(skill string) string {
	if m := pipelineFlowRE.FindStringSubmatch(skill); m != nil {
		block := m[1]
		var out []string
		inFence := false
		for _, ln := range strings.Split(block, "\n") {
			trim := strings.TrimSpace(ln)
			if strings.HasPrefix(trim, "```") {
				inFence = !inFence
				continue
			}
			if inFence {
				out = append(out, ln)
			}
		}
		if len(out) > 0 {
			return strings.TrimSpace(strings.Join(out, "\n"))
		}
	}
	// Fallback: the "> Orchestrate …" blockquote line.
	for _, ln := range strings.Split(skill, "\n") {
		if s := strings.TrimSpace(ln); strings.HasPrefix(s, "> Orchestrate") {
			return strings.TrimPrefix(s, "> ")
		}
	}
	return ""
}

// fallbackRouting produces a routing row per team member when the head's
// team-matrix text carries no explicit "delegate X to Y" clauses.
func fallbackRouting(team []Member) []string {
	var out []string
	for _, m := range team {
		domain := strings.TrimSuffix(m.Role, ".")
		if domain == "" {
			domain = m.Name
		}
		out = append(out, domain+" → "+titleCase(m.Name))
	}
	return out
}

// reportingFallback builds an output-format block from the head's reporting
// bullets, or a generic verdict template.
func reportingFallback(p *Persona) string {
	if len(p.Directives) > 0 {
		var out []string
		for _, d := range p.Directives {
			out = append(out, "- "+d)
		}
		return strings.Join(out, "\n")
	}
	return "[Verdict]: …\n[Rationale]: …\n[Risk]: LOW / MED / HIGH"
}
