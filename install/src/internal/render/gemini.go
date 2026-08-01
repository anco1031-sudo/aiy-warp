package render

import (
	"fmt"
	"strings"
	"time"

	"github.com/anco1031-sudo/aiy-warp/install/src/internal/condense"
)

// geminiRenderer emits Gem "Instructions" text. Gem instructions are markdown;
// skills can be attached as files where the host allows.
type geminiRenderer struct{}

func (geminiRenderer) Name() string { return "gemini" }

func (geminiRenderer) Render(p *condense.Persona) (string, error) {
	if len(p.Team) > 0 {
		return renderGeminiConductor(p), nil
	}
	return renderGeminiSingle(p), nil
}

func renderGeminiSingle(p *condense.Persona) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# ✨ %s — %s (Gem instructions)\n\n", condense.HeaderName(p), condense.FirstSentence(p.Role))
	fmt.Fprintf(&b, "You are **%s**, %s", p.DisplayName, strings.ToLower(lead(condense.FirstSentence(p.Role))))
	if p.ReportsTo != "" {
		fmt.Fprintf(&b, ", reporting to %s", p.ReportsTo)
	}
	fmt.Fprintf(&b, ".\n")
	fmt.Fprintf(&b, "Exported by `aiy warp export gemini --agent %s` · %s · Spec §3.1 (condensed).\n\n",
		p.Name, time.Now().Format("2006-01-02"))

	b.WriteString("## Core identity\n\n")
	b.WriteString(p.Personality + "\n\n")

	if len(p.Directives) > 0 {
		b.WriteString("## Directives\n\n")
		for _, d := range p.Directives {
			fmt.Fprintf(&b, "- %s\n", d)
		}
		b.WriteString("\n")
	}
	if len(p.Boundaries) > 0 {
		b.WriteString("## Boundaries\n\n")
		for _, d := range p.Boundaries {
			fmt.Fprintf(&b, "- %s\n", d)
		}
		b.WriteString("\n")
	}
	if strings.TrimSpace(p.OutputFormat) != "" {
		b.WriteString("## Output format\n\n")
		b.WriteString(p.OutputFormat + "\n\n")
	}
	b.WriteString("> Attach owned skills (SKILL.md) as files where Gemini allows; otherwise paste their workflow sections.\n")
	return b.String()
}

func renderGeminiConductor(p *condense.Persona) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# ✨ %s TEAM — %s (Gem conductor, collapsed from %d agents)\n\n",
		strings.ToUpper(p.Name), condense.FirstSentence(p.Role), len(p.Team)+1)
	fmt.Fprintf(&b, "You are **%s**, %s", p.DisplayName, strings.ToLower(lead(condense.FirstSentence(p.Role))))
	if p.ReportsTo != "" {
		fmt.Fprintf(&b, ", reporting to %s", p.ReportsTo)
	}
	fmt.Fprintf(&b, ".\n")
	fmt.Fprintf(&b, "You have **no subagents on this platform** — internally role-play each specialist, then synthesize ONE final answer.\n\n")

	b.WriteString("## Your team (internal personas)\n\n")
	for _, m := range p.Team {
		fmt.Fprintf(&b, "- **%s** — %s. %s\n", m.DisplayName, strings.TrimSuffix(m.Role, "."), m.OneLiner)
	}
	b.WriteString("\n")

	if strings.TrimSpace(p.Pipeline) != "" {
		b.WriteString("## Pipeline (always run in this order)\n\n```\n")
		b.WriteString(p.Pipeline + "\n```\n\n")
	}

	if len(p.Routing) > 0 {
		b.WriteString("## Routing table\n\n| If the request concerns… | Run lens |\n|---|---|\n")
		for _, row := range p.Routing {
			what, target, _ := strings.Cut(row, " → ")
			fmt.Fprintf(&b, "| %s | → %s |\n", what, target)
		}
		b.WriteString("\n")
	}

	if len(p.Directives) > 0 {
		b.WriteString("## Directives\n\n")
		for _, d := range p.Directives {
			fmt.Fprintf(&b, "- %s\n", d)
		}
		b.WriteString("\n")
	}
	if len(p.Boundaries) > 0 {
		b.WriteString("## Boundaries\n\n")
		for _, d := range p.Boundaries {
			fmt.Fprintf(&b, "- %s\n", d)
		}
		b.WriteString("\n")
	}
	if strings.TrimSpace(p.OutputFormat) != "" {
		b.WriteString("## Output format (every final answer)\n\n```\n")
		b.WriteString(p.OutputFormat + "\n```\n")
	}
	return b.String()
}
