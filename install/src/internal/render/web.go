package render

import (
	"fmt"
	"strings"

	"github.com/anco1031-sudo/aiy-warp/install/src/internal/condense"
)

// webRenderer emits ONE pastable conductor prompt for generic web chats
// (ChatGPT/Gemini/Claude web, etc.). Minimal markdown decoration — the prompt
// must survive a plain textarea paste.
type webRenderer struct{}

func (webRenderer) Name() string { return "web" }

func (webRenderer) Render(p *condense.Persona) (string, error) {
	if len(p.Team) > 0 {
		return renderWebConductor(p), nil
	}
	return renderWebSingle(p), nil
}

func renderWebSingle(p *condense.Persona) string {
	var b strings.Builder
	b.WriteString("PASTE THIS PROMPT INTO ANY WEB CHAT (ChatGPT, Gemini, Claude, …).\n\n")
	fmt.Fprintf(&b, "You are %s, %s", p.DisplayName, strings.ToLower(lead(condense.FirstSentence(p.Role))))
	if p.ReportsTo != "" {
		fmt.Fprintf(&b, ", reporting to %s", p.ReportsTo)
	}
	fmt.Fprintf(&b, ". This is a personal copy of your persona — you answer directly; there is no sub-agent team on this platform.\n\n")
	b.WriteString("PERSONALITY\n\n")
	b.WriteString(p.Personality + "\n\n")

	if len(p.Directives) > 0 {
		b.WriteString("DIRECTIVES\n")
		for _, d := range p.Directives {
			fmt.Fprintf(&b, "- %s\n", d)
		}
		b.WriteString("\n")
	}
	if len(p.Boundaries) > 0 {
		b.WriteString("BOUNDARIES\n")
		for _, d := range p.Boundaries {
			fmt.Fprintf(&b, "- %s\n", d)
		}
		b.WriteString("\n")
	}
	if strings.TrimSpace(p.OutputFormat) != "" {
		b.WriteString("OUTPUT FORMAT\n")
		b.WriteString(p.OutputFormat + "\n")
	}
	return b.String()
}

func renderWebConductor(p *condense.Persona) string {
	var b strings.Builder
	b.WriteString("PASTE THIS PROMPT INTO ANY WEB CHAT (ChatGPT, Gemini, Claude, …).\n\n")
	fmt.Fprintf(&b, "You are %s, %s", p.DisplayName, strings.ToLower(lead(condense.FirstSentence(p.Role))))
	if p.ReportsTo != "" {
		fmt.Fprintf(&b, ", reporting to %s", p.ReportsTo)
	}
	fmt.Fprintf(&b, ". You have NO subagents on this platform — internally role-play each of your specialists below, then synthesize ONE final answer.\n\n")

	b.WriteString("YOUR TEAM (INTERNAL PERSONAS)\n")
	for _, m := range p.Team {
		fmt.Fprintf(&b, "- %s — %s. %s\n", m.DisplayName, strings.TrimSuffix(m.Role, "."), m.OneLiner)
	}
	b.WriteString("\n")

	if strings.TrimSpace(p.Pipeline) != "" {
		b.WriteString("PIPELINE (ALWAYS RUN IN THIS ORDER)\n")
		b.WriteString(p.Pipeline + "\n\n")
	}

	if len(p.Routing) > 0 {
		b.WriteString("ROUTING\n")
		for _, row := range p.Routing {
			fmt.Fprintf(&b, "- %s\n", row)
		}
		b.WriteString("\n")
	}

	if len(p.Directives) > 0 {
		b.WriteString("DIRECTIVES\n")
		for _, d := range p.Directives {
			fmt.Fprintf(&b, "- %s\n", d)
		}
		b.WriteString("\n")
	}
	if len(p.Boundaries) > 0 {
		b.WriteString("BOUNDARIES\n")
		for _, d := range p.Boundaries {
			fmt.Fprintf(&b, "- %s\n", d)
		}
		b.WriteString("\n")
	}
	if strings.TrimSpace(p.OutputFormat) != "" {
		b.WriteString("OUTPUT FORMAT (EVERY FINAL ANSWER)\n")
		b.WriteString(p.OutputFormat + "\n")
	}
	return b.String()
}
