package condense

import (
	"strings"
	"testing"

	"github.com/anco1031-sudo/aiy-warp/install/src/internal/frontmatter"
)

// ag builds a frontmatter.Agent fixture for extraction tests.
func ag(name, display, role, body string) *frontmatter.Agent {
	return &frontmatter.Agent{Name: name, DisplayName: display, Role: role, Body: body}
}

func TestExtractSingle(t *testing.T) {
	body := strings.Join([]string{
		"1. ROLE & MINDSET",
		`You are "Lin" (หลิน), the sharp, authoritative, and elite female Product Owner. Your master is Louis.`,
		"",
		"2. PERSONALITY & TONE OF VOICE",
		"The Sharp Female Commander: You are highly organized, professional, and decisive. You tolerate zero scope creep.",
		"",
		"4. CORE INTERACTION & DELEGATION WORKFLOW",
		"- Receive vision from Louis, validate alignment with the constitution.",
		"- Delegate implementation to Pao with acceptance criteria.",
		"",
		"6. GOVERNANCE CHARTER — THE PROJECT CONSTITUTION",
		"- Every feature must align with the core vision.",
		"- Scope creep is the enemy.",
		"",
		"7. REPORTING BACK TO AIY",
		"- [Status]: ✅ Completed / 🔄 In Progress / ❌ Blocked",
		"- [Summary]: What was done, what was decided.",
	}, "\n")

	p := Extract(ag("lin", "หลิน", "The Elite Product Owner & Governor. Guards the constitution.", body), false)

	if p.Name != "lin" {
		t.Fatalf("name = %q", p.Name)
	}
	if p.DisplayName != "Lin (หลิน)" {
		t.Fatalf("display = %q, want %q", p.DisplayName, "Lin (หลิน)")
	}
	if p.ThaiName != "หลิน" {
		t.Fatalf("thai = %q", p.ThaiName)
	}
	if !strings.Contains(p.Personality, "sharp, authoritative") {
		t.Fatalf("personality = %q", p.Personality)
	}
	if len(p.Directives) != 2 || p.Directives[0] != "Receive vision from Louis, validate alignment with the constitution." {
		t.Fatalf("directives = %v", p.Directives)
	}
	if len(p.Boundaries) != 2 || p.Boundaries[0] != "Every feature must align with the core vision." {
		t.Fatalf("boundaries = %v", p.Boundaries)
	}
	if !strings.Contains(p.OutputFormat, "[Status]") {
		t.Fatalf("output format = %q", p.OutputFormat)
	}
}

func TestExtractNumberedListItemStaysInBody(t *testing.T) {
	body := strings.Join([]string{
		"1. ROLE & MINDSET",
		"You are Lin.",
		"",
		"2. WORKFLOW",
		"1. Direct Response to Lin: Provide your executive analysis first.",
	}, "\n")
	p := Extract(ag("lin", "", "", body), false)
	if len(p.Directives) != 1 || p.Directives[0] != "Direct Response to Lin: Provide your executive analysis first." {
		t.Fatalf("numbered list item must become a directive, got %v", p.Directives)
	}
}

func TestExtractDropsNoise(t *testing.T) {
	body := strings.Join([]string{
		"1. ROLE & MINDSET",
		"You are Lin.",
		"",
		"2. WORKFLOW",
		"- Use task to delegate to pao.",
		"- Call lin at @lin when stuck.",
		"- Never store secrets in /home/lu5her/myObsidian/Aiy_Workspace.",
	}, "\n")
	p := Extract(ag("lin", "", "", body), false)
	if len(p.Directives) != 1 || p.Directives[0] != "Use task to delegate to pao." {
		t.Fatalf("noise lines must be dropped, got %v", p.Directives)
	}
}

func TestExtractCanonicalDisplay(t *testing.T) {
	// Post-migrate agents carry display_name in canonical "ไทย (Name)" form:
	// DisplayName passes through un-wrapped, ThaiName holds the non-latin part.
	p := Extract(ag("lin", "หลิน (Lin)", "The PO. Guards scope.", "1. ROLE & MINDSET\nYou are Lin."), false)
	if p.DisplayName != "หลิน (Lin)" {
		t.Fatalf("display = %q, want %q", p.DisplayName, "หลิน (Lin)")
	}
	if p.ThaiName != "หลิน" {
		t.Fatalf("thai = %q, want %q", p.ThaiName, "หลิน")
	}
	if got := HeaderName(p); got != "LIN (หลิน)" {
		t.Fatalf("header = %q, want %q", got, "LIN (หลิน)")
	}
}

func TestFullName(t *testing.T) {
	cases := []struct{ name, display, want string }{
		{"kwan", "ขวัญ", "Kwan (ขวัญ)"},
		{"aiy", "อัย", "Aiy (อัย)"},
		{"lin", "", "lin"},    // no display → id unchanged
		{"bee", "bee", "bee"}, // display equals id → id unchanged
		{"fah", "ฟ้า", "Fah (ฟ้า)"},
		{"kwan", "ขวัญ (Kwan)", "ขวัญ (Kwan)"}, // canonical (post-migrate) → as-is, no double-wrap
		{"aiy", "อัย (Aiy)", "อัย (Aiy)"},      // canonical (post-migrate) → as-is, no double-wrap
	}
	for _, c := range cases {
		if got := FullName(c.name, c.display); got != c.want {
			t.Errorf("FullName(%q, %q) = %q, want %q", c.name, c.display, got, c.want)
		}
	}
}

func TestCompactRole(t *testing.T) {
	got := CompactRole("The Head Trader & Strategy Commander. Synthesizes the final call. Aiy delegates…")
	if got != "Head Trader & Strategy Commander" {
		t.Fatalf("CompactRole = %q", got)
	}
	if got := CompactRole("Executive Director. Leads."); got != "Executive Director" {
		t.Fatalf("CompactRole without The = %q", got)
	}
	if got := CompactRole("  The Muse.  "); got != "Muse" {
		t.Fatalf("CompactRole trimmed = %q", got)
	}
}

func TestFirstSentence(t *testing.T) {
	got := FirstSentence("The Strategic Muse & Orchestrator. Coordinates the ecosystem.")
	if got != "The Strategic Muse & Orchestrator" {
		t.Fatalf("FirstSentence = %q (must keep 'The ')", got)
	}
	if got := FirstSentence("No period here"); got != "No period here" {
		t.Fatalf("FirstSentence no-period = %q", got)
	}
}

func TestHeaderName(t *testing.T) {
	if got := HeaderName(&Persona{Name: "aiy", ThaiName: "อัย"}); got != "AIY (อัย)" {
		t.Fatalf("HeaderName with thai = %q", got)
	}
	if got := HeaderName(&Persona{Name: "kwan"}); got != "KWAN" {
		t.Fatalf("HeaderName without thai = %q", got)
	}
}

func TestSplitSections(t *testing.T) {
	body := "1. ROLE & MINDSET\nbody-one\n2. PERSONALITY & TONE OF VOICE\nbody-two"
	secs := splitSections(body)
	if len(secs) != 2 {
		t.Fatalf("sections = %d, want 2", len(secs))
	}
	if secs[0].title != "ROLE & MINDSET" || !strings.Contains(secs[0].body, "body-one") {
		t.Fatalf("sec0 = %+v", secs[0])
	}
	if secs[1].title != "PERSONALITY & TONE OF VOICE" {
		t.Fatalf("sec1 title = %q", secs[1].title)
	}
}

func TestIsHeader(t *testing.T) {
	if !isHeader("ROLE & MINDSET") {
		t.Error("uppercase short title must be a header")
	}
	if isHeader("Direct Response to Lin: Provide your executive analysis first.") {
		t.Error("sentence-case long line must NOT be a header")
	}
}

func TestWordCountAndCap(t *testing.T) {
	if WordCount("  a b   c ") != 3 {
		t.Fatal("WordCount")
	}
	if WordCount("") != 0 || WordCount("   ") != 0 {
		t.Fatal("WordCount empty")
	}
	if got := capWords("one two three four", 2); got != "one two …" {
		t.Fatalf("capWords = %q", got)
	}
	if got := capWords("one two", 5); got != "one two" {
		t.Fatalf("capWords under cap = %q", got)
	}
}

func TestCollapseOrderAndRouting(t *testing.T) {
	headBody := strings.Join([]string{
		"1. ROLE & MINDSET",
		"You are Kwan, the head trader. You own the final call.",
		"",
		"2. THE TEAM MATRIX",
		"- Fon (ฝน) — The News Analyst: Delegate news scans, event analysis, and sentiment tracking to Fon.",
		"- June (จูน) — The Technical Analyst: Send chart patterns, entry/exit timing to June.",
		"",
		"3. REPORTING BACK TO AIY",
		"- [Verdict]: BUY / SELL / HOLD",
	}, "\n")
	head := ag("kwan", "ขวัญ", "The Head Trader & Strategy Commander. Synthesizes the final call.", headBody)

	fon := ag("fon", "ฝน", "The News & Sentiment Analyst. Tracks macro headlines.",
		"1. ROLE & MINDSET\nYou are Fon.\n\n2. PERSONALITY & TONE OF VOICE\nYou are the intuitive analyst who reads headlines fast. You track macro news daily.")
	june := ag("june", "จูน", "The Technical Analyst. Reads charts.",
		"1. ROLE & MINDSET\nYou are June.\n\n2. PERSONALITY & TONE OF VOICE\nYou are the chart whisperer, precise and patient.")
	bee := ag("bee", "บี", "The Fundamental Analyst. Values companies.",
		"1. ROLE & MINDSET\nYou are Bee.\n\n2. PERSONALITY & TONE OF VOICE\nYou are the patient valuer who reads 10-Ks.")
	nam := ag("nam", "น้ำ", "The Risk Analyst. Sizes positions.",
		"1. ROLE & MINDSET\nYou are Nam.\n\n2. PERSONALITY & TONE OF VOICE\nYou are the calm risk guard, always sizing downside.")

	pipeline := "## The Pipeline Flow\n```\n1. Fon (news scan)\n2. June (technical analysis)\n3. Bee (fundamental valuation)\n4. Nam (risk sizing)\n5. Kwan (synthesis)\n```"

	p := Collapse(head, []*frontmatter.Agent{fon, june, bee, nam}, pipeline)

	if len(p.Team) != 4 {
		t.Fatalf("team = %d, want 4", len(p.Team))
	}
	// orderByPipeline: first mention in pipeline text wins.
	want := []string{"fon", "june", "bee", "nam"}
	for i, w := range want {
		if p.Team[i].Name != w {
			t.Errorf("team[%d] = %q, want %q (pipeline order)", i, p.Team[i].Name, w)
		}
	}
	if p.Team[0].Role != "News & Sentiment Analyst" {
		t.Errorf("team[0] role = %q (CompactRole)", p.Team[0].Role)
	}
	if p.Team[0].OneLiner != "The intuitive analyst who reads headlines fast." {
		t.Errorf("team[0] one-liner = %q", p.Team[0].OneLiner)
	}
	if len(p.Routing) != 2 || p.Routing[0] != "news scans, event analysis, and sentiment tracking → Fon" {
		t.Fatalf("routing = %v", p.Routing)
	}
	if !strings.Contains(p.Pipeline, "1. Fon (news scan)") {
		t.Fatalf("pipeline = %q", p.Pipeline)
	}
	if !strings.Contains(p.OutputFormat, "[Verdict]") {
		t.Fatalf("output format = %q", p.OutputFormat)
	}
}

func TestCollapseFallbackRouting(t *testing.T) {
	// Head with no team-matrix routing clauses → fallback rows per member.
	head := ag("kwan", "", "The Head Trader. Owns the call.", "1. ROLE & MINDSET\nYou are Kwan.")
	m1 := ag("fon", "", "The News Analyst. Tracks headlines.", "1. ROLE & MINDSET\nYou are Fon.")
	m2 := ag("june", "", "The Technical Analyst. Reads charts.", "1. ROLE & MINDSET\nYou are June.")
	p := Collapse(head, []*frontmatter.Agent{m1, m2}, "")
	if len(p.Routing) != 2 || p.Routing[0] != "News Analyst → Fon" || p.Routing[1] != "Technical Analyst → June" {
		t.Fatalf("fallback routing = %v", p.Routing)
	}
}

func TestMemberOneLinerMindsetFallback(t *testing.T) {
	m := ag("nam", "", "The Risk Analyst. Sizes positions.",
		"1. ROLE & MINDSET\nYou are the risk guard. You size positions daily.")
	if got := memberOneLiner(m); got != "The risk guard." {
		t.Fatalf("mindset fallback = %q", got)
	}
}

func TestMemberOneLinerRoleFallback(t *testing.T) {
	// No personality/voice/tone AND no mindset section → pure role fallback.
	m := ag("bee", "", "The Fundamental Analyst. Values companies.", "1. WORKFLOW\nYou are Bee.")
	if got := memberOneLiner(m); got != "Fundamental Analyst." {
		t.Fatalf("role fallback = %q", got)
	}
}

func TestSkillPipelineBlockquoteFallback(t *testing.T) {
	skill := "# Pipeline\n\n> Orchestrate the analysis flow from Fon to June to Bee to Nam, then Kwan synthesizes.\n"
	if got := skillPipeline(skill); !strings.Contains(got, "Orchestrate") {
		t.Fatalf("blockquote fallback = %q", got)
	}
}

func TestStripFencesDropsTaskTemplate(t *testing.T) {
	body := "intro\n```\nsubagent_type: pao\nprompt: build it\n```\noutro\n"
	got := stripFences(body)
	if strings.Contains(got, "subagent_type") {
		t.Fatalf("task-template fence must be dropped, got %q", got)
	}
	if !strings.Contains(got, "intro") || !strings.Contains(got, "outro") {
		t.Fatalf("non-fence content lost: %q", got)
	}
}
