package main

import (
	"flag"
	"fmt"
	"io"

	"github.com/anco1031-sudo/aiy-warp/install/src/internal/config"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/doctor"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/errs"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/ui"
)

func cmdDoctor(u *ui.UI, rest []string) int {
	platform, rest := extractPlatform("doctor", rest)
	fs := flag.NewFlagSet("doctor", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	platformFlag := fs.String("platform", "opencode", "platform to check (P0: opencode)")
	jsonOut := fs.Bool("json", false, "machine-readable JSON report")
	cfgPath := fs.String("config", "", "path to warp.config")
	verbose := fs.Bool("verbose", false, "verbose output")
	if handled, code := parseCmdFlags(u, fs, rest); handled {
		return code
	}
	p := *platformFlag
	if platform != "" {
		p = platform
	}
	if p != "opencode" {
		fmt.Fprintf(u.Err, "aiy warp: unknown platform %q (P0 supports: opencode)\n", p)
		return errs.CodeUsage
	}

	repoRoot, cfg, _, code := setup(u, *cfgPath, *verbose)
	if code != errs.CodeOK {
		return code
	}

	res, err := doctor.Run(doctor.Options{
		RepoRoot:   repoRoot,
		DestDir:    cfg.OpencodeDest(),
		StateDir:   cfg.StateDirResolved(),
		ConfigPath: pick(*cfgPath, config.DefaultPath()),
	})
	if err != nil {
		return reportErr(u, err)
	}

	if *jsonOut {
		if err := u.JSON(res); err != nil {
			fmt.Fprintln(u.Err, "aiy warp:", err)
			return errs.CodeRuntime
		}
		return res.Exit
	}
	u.Doctor(res)
	if res.Exit == 0 {
		u.OK("Doctor: all checks pass.")
	} else {
		u.Warn("Doctor: drift detected — run 'aiy warp install opencode' to fix.")
	}
	return res.Exit
}
