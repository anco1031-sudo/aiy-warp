package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/anco1031-sudo/aiy-warp/install/src/internal/config"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/errs"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/kit"
	"github.com/anco1031-sudo/aiy-warp/install/src/internal/ui"
)

// setup loads shared context: repo root, config, bundles.
func setup(u *ui.UI, cfgPath string, verbose bool) (string, *config.Config, *kit.Bundles, int) {
	u.Verbose = verbose
	repoRoot, err := locateRepo()
	if err != nil {
		fmt.Fprintln(u.Err, "aiy warp:", err)
		return "", nil, nil, errs.CodeRuntime
	}
	cfg, err := config.Load(pick(cfgPath, config.DefaultPath()))
	if err != nil {
		fmt.Fprintln(u.Err, "aiy warp: config:", err)
		return "", nil, nil, errs.CodeRuntime
	}
	bundles, err := kit.LoadBundles(filepath.Join(repoRoot, "install", "bundles.yaml"))
	if err != nil {
		fmt.Fprintln(u.Err, "aiy warp: bundles.yaml:", err)
		return "", nil, nil, errs.CodeRuntime
	}
	return repoRoot, cfg, bundles, errs.CodeOK
}

// locateRepo finds the warp kit root from CWD, walking upward.
func locateRepo() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir := cwd
	for i := 0; i < 6; i++ {
		if _, err := os.Stat(filepath.Join(dir, "install", "bundles.yaml")); err == nil {
			if _, err := os.Stat(filepath.Join(dir, "agents")); err == nil {
				return dir, nil
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("warp kit not found (looked for install/bundles.yaml + agents/ from %s upward); run from the aiy-warp repo root", cwd)
}

// valueFlags lists per-command flags that consume a following value; used by
// extractPlatform to pull the positional <platform> out of argv while keeping
// flag/value pairs intact (stdlib flag stops at the first positional).
var valueFlags = map[string]map[string]bool{
	"install": {"agent": true, "team": true, "config": true, "allow-identifiers": true},
	"export":  {"agent": true, "team": true, "out": true, "config": true, "allow-identifiers": true},
	"doctor":  {"platform": true, "config": true},
	"list":    {"platform": true, "config": true},
}

// parseCmdFlags parses fs against rest. On -h/--help it prints the command's
// flags to u.Out and reports handled=true with exit code 0; on a parse error it
// reports handled=true with exit code 2; on success handled=false.
func parseCmdFlags(u *ui.UI, fs *flag.FlagSet, rest []string) (bool, int) {
	if err := fs.Parse(rest); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			fmt.Fprintf(u.Out, "aiy warp %s\nUsage: aiy warp %s [flags]\n\nFlags:\n", fs.Name(), fs.Name())
			fs.SetOutput(u.Out)
			fs.PrintDefaults()
			return true, errs.CodeOK
		}
		fmt.Fprintln(u.Err, "aiy warp "+fs.Name()+":", err)
		return true, errs.CodeUsage
	}
	return false, 0
}

// extractPlatform removes the first bare (non-flag) token from argv and
// returns it as platform, keeping all flags and their values intact.
func extractPlatform(cmd string, argv []string) (string, []string) {
	vf := valueFlags[cmd]
	if vf == nil {
		vf = map[string]bool{}
	}
	var rest []string
	platform := ""
	for i := 0; i < len(argv); i++ {
		a := argv[i]
		if strings.HasPrefix(a, "-") {
			rest = append(rest, a)
			name := strings.TrimLeft(a, "-")
			if eq := strings.Index(name, "="); eq >= 0 {
				continue
			}
			if vf[name] && i+1 < len(argv) {
				rest = append(rest, argv[i+1])
				i++
			}
			continue
		}
		if platform == "" {
			platform = a
		} else {
			rest = append(rest, a)
		}
	}
	return platform, rest
}
