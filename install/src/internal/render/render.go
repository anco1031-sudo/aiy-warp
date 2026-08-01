// Package render formats a condensed Persona for each web-chat platform
// (CLI_SPEC.md §3 renderer matrix): ChatGPT Custom GPT instructions, Gemini Gem
// instructions, and a pastable web-chat conductor prompt.
package render

import (
	"fmt"
	"sort"
	"strings"

	"github.com/anco1031-sudo/aiy-warp/install/src/internal/condense"
)

// Renderer formats a condensed persona for one platform.
type Renderer interface {
	// Name is the platform key (chatgpt | gemini | web).
	Name() string
	// Render emits the platform-native persona text. A persona with a non-empty
	// Team is a collapsed team conductor (§3.2); otherwise it is a single agent.
	Render(p *condense.Persona) (string, error)
}

// Registry maps platform names to renderers.
var Registry = map[string]Renderer{
	"chatgpt": chatgptRenderer{},
	"gemini":  geminiRenderer{},
	"web":     webRenderer{},
}

// Names returns the registered platform names, sorted.
func Names() []string {
	out := make([]string, 0, len(Registry))
	for n := range Registry {
		out = append(out, n)
	}
	sort.Strings(out)
	return out
}

// Get returns the renderer for a platform, or an error for unknown platforms.
func Get(platform string) (Renderer, error) {
	r, ok := Registry[platform]
	if !ok {
		return nil, fmt.Errorf("unknown render platform %q (have: %s)", platform, strings.Join(Names(), ", "))
	}
	return r, nil
}

// sourceList renders the source-agent annotation, e.g.
// "`agents/kwan.md` + `agents/fon.md`" — backticked relpaths joined by " + ".
func sourceList(sources []string) string {
	var out []string
	for _, s := range sources {
		out = append(out, "`"+s+"`")
	}
	return strings.Join(out, " + ")
}
