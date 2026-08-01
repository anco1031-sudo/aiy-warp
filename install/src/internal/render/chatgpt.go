package render

import (
	"fmt"
	"strings"
	"time"

	"github.com/anco1031-sudo/aiy-warp/install/src/internal/condense"
)

// chatgptRenderer emits Custom GPT "Instructions" text. Single agents follow
// the shape of export/examples/chatgpt-aiy-persona.md; collapsed teams follow
// export/examples/chatgpt-kwan-conductor.md (CLI_SPEC.md §3.2).
type chatgptRenderer struct{}

func (chatgptRenderer) Name() string { return "chatgpt" }

func (chatgptRenderer) Render(p *condense.Persona) (string, error) {
	if len(p.Team) > 0 {
		return renderChatGPTConductor(p), nil
	}
	return renderChatGPTSingle(p), nil
}

func renderChatGPTSingle(p *condense.Persona) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# 🌀 %s — %s (personal copy)\n\n", condense.HeaderName(p), condense.FirstSentence(p.Role))
	fmt.Fprintf(&b, "> Exported by `aiy warp export chatgpt --agent %s`\n", p.Name)
	fmt.Fprintf(&b, "> Generated %s · Source: `agents/%s.md` · Spec: `install/CLI_SPEC.md` §3.1 (condensed)\n\n",
		time.Now().Format("2006-01-02"), p.Name)
	fmt.Fprintf(&b, "You are **%s**, %s. ", p.DisplayName, strings.ToLower(lead(condense.FirstSentence(p.Role))))
	if p.ReportsTo != "" {
		fmt.Fprintf(&b, "You report to %s. ", p.ReportsTo)
	}
	fmt.Fprintf(&b, "This is a personal copy: a single-agent vessel of your persona, without the sub-agent team.\n\n")

	b.WriteString("## Personality\n\n")
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
	b.WriteString("---\n\n*Condensed persona v1 — single-agent vessel per WARP.md §3.1. Soul preserved; delegation layer removed by platform limits.*\n")
	return b.String()
}

func renderChatGPTConductor(p *condense.Persona) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# 🎯 %s TEAM — %s (collapsed from %d agents)\n\n",
		strings.ToUpper(p.Name), condense.FirstSentence(p.Role), len(p.Team)+1)
	fmt.Fprintf(&b, "> Exported by `aiy warp export chatgpt --team %s --collapse`\n", p.Dept)
	fmt.Fprintf(&b, "> Generated %s · Source: %s · Spec: `install/CLI_SPEC.md` §3.2\n\n",
		time.Now().Format("2006-01-02"), sourceList(p.Sources))
	fmt.Fprintf(&b, "You are **%s**, %s", p.DisplayName, strings.ToLower(lead(condense.FirstSentence(p.Role))))
	if p.ReportsTo != "" {
		fmt.Fprintf(&b, ", reporting to %s", p.ReportsTo)
	}
	fmt.Fprintf(&b, ".\n\n")
	b.WriteString("You have **NO subagents on this platform** — you internally role-play each specialist, then synthesize ONE final answer.\n\n")

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

// lead lowercases the first rune of s.
func lead(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = []rune(strings.ToLower(string(r[0])))[0]
	return string(r)
}
