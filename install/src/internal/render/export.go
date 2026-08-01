package render

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/anco1031-sudo/aiy-warp/install/src/internal/condense"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/errs"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/frontmatter"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/kit"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/redact"
)

// ExportOptions configures `aiy warp export chatgpt|gemini|web`.
type ExportOptions struct {
	RepoRoot   string
	Out        string // output dir (ignored when Stdout)
	Stdout     bool
	Platform   string
	Bundles    *kit.Bundles
	Selector   kit.Selector
	AllowIDs   map[string]bool
	NoCollapse bool // explicit --no-collapse (web hosts have no subagents → usage error)
}

// ExportResult reports what was emitted.
type ExportResult struct {
	Files []string // written relative paths (or rendered names for stdout)
	Words int      // word count of the rendered persona (single/condensed)
}

// RunExport condenses + renders the selected agents for a web-chat platform.
// Teams auto-collapse into one conductor. The redaction gate scans the rendered
// payload (exit 5 on credential values / non-allowlisted identifiers).
func RunExport(o ExportOptions) (*ExportResult, error) {
	r, err := Get(o.Platform)
	if err != nil {
		return nil, errs.Usage(err.Error())
	}
	k, err := kit.Discover(o.RepoRoot)
	if err != nil {
		return nil, errs.Wrapf(err, "discover kit")
	}
	files, err := k.Resolve(o.Bundles, o.Selector)
	if err != nil {
		return nil, errs.Usage(err.Error())
	}

	var agents []string
	var skills []string
	for _, f := range files {
		if strings.HasPrefix(f, "agents/") {
			agents = append(agents, f)
		} else if strings.HasPrefix(f, "skills/") {
			skills = append(skills, f)
		}
	}
	if len(agents) == 0 {
		return nil, errs.Usage("selector resolves to no agent files")
	}

	parsed := make([]*frontmatter.Agent, 0, len(agents))
	byName := map[string]*frontmatter.Agent{}
	for _, rel := range agents {
		a, err := parseAgent(o.RepoRoot, rel)
		if err != nil {
			return nil, err
		}
		parsed = append(parsed, a)
		byName[a.Name] = a
	}

	var text string
	var outRel string
	var words int
	if o.Selector.Team != "" {
		if o.NoCollapse {
			return nil, errs.Usage("--no-collapse is not supported on web-chat platforms (no subagents); teams collapse into one conductor")
		}
		dept := o.Bundles.ResolveTeam(o.Selector.Team)
		head := teamHead(o.Bundles, dept, parsed)
		members := teamMembers(o.Bundles, dept, parsed, head)
		pipelineSkill := skillContent(o, k, dept)
		p := condense.Collapse(head, members, pipelineSkill)
		p.Dept = o.Selector.Team // echo the selector the user typed (e.g. "kwan")
		p.Sources = agents
		text, err = r.Render(p)
		words = condense.WordCount(text)
		outRel = dept + "-conductor.md"
	} else {
		a := byName[o.Selector.Agent]
		if a == nil {
			a = parsed[0]
		}
		p := condense.Extract(a, false)
		text, err = r.Render(p)
		words = condense.WordCount(text)
		outRel = a.Name + ".md"
	}
	if err != nil {
		return nil, errs.Wrapf(err, "render %s", o.Platform)
	}

	// Redaction gate on the rendered payload.
	rep := redact.Scan(map[string]string{outRel: text}, o.AllowIDs)
	if rep.HasBlocks() {
		return nil, errs.Secret(fmt.Sprintf("credential value in rendered persona — refusing to export:\n  %s",
			strings.Join(rep.BlockFiles(), "\n  ")))
	}
	if len(rep.Exports) > 0 {
		ids := map[string]bool{}
		for _, f := range rep.Exports {
			ids[f.Match] = true
		}
		keys := make([]string, 0, len(ids))
		for id := range ids {
			keys = append(keys, id)
		}
		sort.Strings(keys)
		return nil, errs.Secret(fmt.Sprintf(
			"sensitive identifiers in rendered persona — refusing to export:\n  %s\nDeclare with --allow-identifiers %s",
			strings.Join(rep.ExportBlockFiles(), "\n  "), strings.Join(keys, ",")))
	}

	if o.Stdout {
		fmt.Print(text)
		if !strings.HasSuffix(text, "\n") {
			fmt.Println()
		}
		return &ExportResult{Files: []string{outRel}, Words: words}, nil
	}
	dest := filepath.Join(o.Out, outRel)
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return nil, errs.Wrapf(err, "mkdir %s", filepath.Dir(dest))
	}
	if err := os.WriteFile(dest, []byte(text), 0o644); err != nil {
		return nil, errs.Wrapf(err, "write %s", dest)
	}
	return &ExportResult{Files: []string{outRel}, Words: words}, nil
}

func parseAgent(repoRoot, rel string) (*frontmatter.Agent, error) {
	b, err := os.ReadFile(filepath.Join(repoRoot, filepath.FromSlash(rel)))
	if err != nil {
		return nil, errs.Wrapf(err, "read %s", rel)
	}
	a, err := frontmatter.ParseAgent(rel, string(b))
	if err != nil {
		return nil, errs.Wrapf(err, "parse %s", rel)
	}
	return a, nil
}

func teamHead(b *kit.Bundles, dept string, parsed []*frontmatter.Agent) *frontmatter.Agent {
	byName := map[string]*frontmatter.Agent{}
	for _, a := range parsed {
		byName[a.Name] = a
	}
	for _, n := range b.Team(dept) {
		if a := byName[n]; a != nil && b.Find(n) != nil && b.Find(n).Rank == "head" {
			return a
		}
	}
	return parsed[0]
}

func teamMembers(b *kit.Bundles, dept string, parsed []*frontmatter.Agent, head *frontmatter.Agent) []*frontmatter.Agent {
	byName := map[string]*frontmatter.Agent{}
	for _, a := range parsed {
		byName[a.Name] = a
	}
	var out []*frontmatter.Agent
	for _, n := range b.Team(dept) {
		if n == head.Name {
			continue
		}
		if a := byName[n]; a != nil {
			out = append(out, a)
		}
	}
	return out
}

// skillContent returns the pipeline skill's SKILL.md content for a department
// (e.g. trading → skills/trading-pipeline/SKILL.md), or "" when absent.
func skillContent(o ExportOptions, k *kit.Kit, dept string) string {
	var skillName string
	for _, s := range o.Bundles.Skills {
		for _, d := range s.Departments {
			if d == dept {
				skillName = s.Name
			}
		}
	}
	if skillName == "" {
		return ""
	}
	for _, rel := range k.Skills {
		if filepath.Base(filepath.Dir(filepath.FromSlash(rel))) == skillName {
			if b, err := os.ReadFile(filepath.Join(o.RepoRoot, filepath.FromSlash(rel))); err == nil {
				return string(b)
			}
		}
	}
	return ""
}
