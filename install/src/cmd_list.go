package main

import (
	"flag"
	"fmt"
	"io"

	"github.com/anco1031-sudo/aiy-warp/install/src/internal/errs"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/kit"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/ui"
)

func cmdList(u *ui.UI, rest []string) int {
	platform, rest := extractPlatform("list", rest)
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	platformFlag := fs.String("platform", "opencode", "platform to list (P0: opencode)")
	verbose := fs.Bool("verbose", false, "verbose output")
	if err := fs.Parse(rest); err != nil {
		fmt.Fprintln(u.Err, "aiy warp list:", err)
		return errs.CodeUsage
	}
	p := *platformFlag
	if platform != "" {
		p = platform
	}
	if p != "opencode" {
		fmt.Fprintf(u.Err, "aiy warp: unknown platform %q (P0 supports: opencode)\n", p)
		return errs.CodeUsage
	}

	repoRoot, _, bundles, code := setup(u, "", *verbose)
	if code != errs.CodeOK {
		return code
	}
	k, err := kit.Discover(repoRoot)
	if err != nil {
		return reportErr(u, err)
	}

	u.Info("Kit: %d agents, %d skills, %d templates, %d playbooks (platform: %s)",
		len(k.Agents), len(k.Skills), len(k.Templates), len(k.Playbooks), p)
	u.Info("")
	u.Info("AGENT  RANK      DEPT      REPORTS_TO")
	u.Info("-----  --------  --------  ----------")
	for _, name := range bundles.AgentNames() {
		u.Info(bundles.Describe(name))
	}
	u.Info("")
	u.Info("SKILLS:")
	for _, s := range bundles.SkillNames() {
		owner := ""
		for _, sk := range bundles.Skills {
			if sk.Name == s {
				owner = sk.OwnedBy
			}
		}
		u.Info("  %-22s owned_by: %s", s, owner)
	}
	return errs.CodeOK
}
