package render

import (
	"strings"
	"testing"

	"github.com/anco1031-sudo/aiy-warp/install/src/internal/condense"
)

func TestRegistryGet(t *testing.T) {
	for _, n := range []string{"chatgpt", "gemini", "web"} {
		if _, err := Get(n); err != nil {
			t.Errorf("Get(%q) = %v", n, err)
		}
	}
	if _, err := Get("slack"); err == nil || !strings.Contains(err.Error(), "unknown render platform") {
		t.Fatalf("Get(slack) err = %v, want unknown-render error", err)
	}
}

func TestNamesSorted(t *testing.T) {
	got := Names()
	want := []string{"chatgpt", "gemini", "web"}
	if len(got) != len(want) {
		t.Fatalf("Names = %v", got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Names = %v, want %v", got, want)
		}
	}
}

func TestSourceList(t *testing.T) {
	if got := sourceList([]string{"agents/kwan.md", "agents/fon.md"}); got != "`agents/kwan.md` + `agents/fon.md`" {
		t.Fatalf("sourceList = %q", got)
	}
	if got := sourceList(nil); got != "" {
		t.Fatalf("sourceList(nil) = %q", got)
	}
}

func TestChatGPTSingle(t *testing.T) {
	p := &condense.Persona{
		Name:        "aiy",
		DisplayName: "Aiy (อัย)",
		ThaiName:    "อัย",
		Role:        "The Strategic Muse & Orchestrator. Coordinates the ecosystem.",
		Personality: "You are the sharp orchestrator of Louis's ecosystem, decisive and elegant.",
	}
	r, err := Get("chatgpt")
	if err != nil {
		t.Fatal(err)
	}
	out, err := r.Render(p)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "# 🌀 AIY (อัย) — The Strategic Muse & Orchestrator (personal copy)") {
		t.Errorf("header missing: %q", out[:80])
	}
	if !strings.Contains(out, "You are **Aiy (อัย)**, the strategic muse & orchestrator.") {
		t.Errorf("opening line missing")
	}
	if !strings.Contains(out, "personal copy") {
		t.Errorf("personal-copy marker missing")
	}
	// FIX: role period must terminate the title line, no glued "orchestrator.This".
	if strings.Contains(out, "orchestrator.This") || strings.Contains(out, "Commander.Synthesizes") {
		t.Errorf("role period defect: %q", out[:200])
	}
}

func TestChatGPTConductor(t *testing.T) {
	p := &condense.Persona{
		Name:        "kwan",
		DisplayName: "Kwan (ขวัญ)",
		ThaiName:    "ขวัญ",
		Role:        "The Head Trader & Strategy Commander. Synthesizes the final call.",
		Dept:        "kwan",
		Team: []condense.Member{
			{Name: "fon", DisplayName: "Fon (ฝน)", Role: "News & Sentiment Analyst", OneLiner: "The intuitive analyst who reads headlines fast."},
			{Name: "june", DisplayName: "June (จูน)", Role: "Technical Analyst", OneLiner: "The chart whisperer."},
			{Name: "bee", DisplayName: "Bee (บี)", Role: "Fundamental Analyst", OneLiner: "The patient valuer."},
			{Name: "nam", DisplayName: "Nam (น้ำ)", Role: "Risk Analyst", OneLiner: "The calm risk guard."},
		},
		Pipeline:     "1. Fon (news scan)\n2. June (technical analysis)\n3. Bee (fundamental valuation)\n4. Nam (risk sizing)\n5. Kwan (synthesis)",
		Routing:      []string{"News & Sentiment Analyst → Fon", "Technical Analyst → June"},
		Sources:      []string{"agents/kwan.md", "agents/fon.md", "agents/june.md", "agents/bee.md", "agents/nam.md"},
		Directives:   []string{"Run every lens before answering."},
		Boundaries:   []string{"Never trade on a single lens."},
		OutputFormat: "- [Verdict]: BUY / SELL / HOLD\n- [Rationale]: thesis",
	}
	r, err := Get("chatgpt")
	if err != nil {
		t.Fatal(err)
	}
	out, err := r.Render(p)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "# 🎯 KWAN TEAM — The Head Trader & Strategy Commander (collapsed from 5 agents)") {
		t.Errorf("conductor header missing")
	}
	if !strings.Contains(out, "Exported by `aiy warp export chatgpt --team kwan --collapse`") {
		t.Errorf("export line missing")
	}
	if !strings.Contains(out, "Source: `agents/kwan.md` + `agents/fon.md` + `agents/june.md` + `agents/bee.md` + `agents/nam.md`") {
		t.Errorf("source list missing")
	}
	if !strings.Contains(out, "- **Fon (ฝน)** — News & Sentiment Analyst. The intuitive analyst who reads headlines fast.") {
		t.Errorf("member line missing")
	}
	if !strings.Contains(out, "## Routing table") || !strings.Contains(out, "| News & Sentiment Analyst | → Fon |") {
		t.Errorf("routing table missing")
	}
	if !strings.Contains(out, "## Pipeline (always run in this order)") {
		t.Errorf("pipeline block missing")
	}
}

func TestRenderWordCap(t *testing.T) {
	// A large conductor must stay under the 1500-word cap (CLI_SPEC §3.1).
	team := make([]condense.Member, 0, 5)
	routing := make([]string, 0, 5)
	for _, name := range []string{"fon", "june", "bee", "nam"} {
		team = append(team, condense.Member{Name: name, DisplayName: name, Role: "Analyst " + name, OneLiner: "One liner for " + name + "."})
		routing = append(routing, "domain "+name+" → "+name)
	}
	p := &condense.Persona{
		Name: "kwan", DisplayName: "Kwan", Role: "The Head Trader. Owns the call.",
		Dept: "kwan", Team: team, Pipeline: "pipeline flow", Routing: routing,
	}
	r, _ := Get("chatgpt")
	out, err := r.Render(p)
	if err != nil {
		t.Fatal(err)
	}
	if w := condense.WordCount(out); w > condense.WordCap {
		t.Fatalf("rendered %d words, cap %d", w, condense.WordCap)
	}
}
