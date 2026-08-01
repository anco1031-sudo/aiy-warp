// Package config loads the optional host-local warp.config (CLI_SPEC.md §5.3).
// All fields are optional in P0 — opencode destinations default under $HOME.
package config

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config is the host-local configuration. A missing file yields zero values;
// callers fall back to $HOME-derived defaults.
type Config struct {
	Version int `yaml:"version"`
	Paths   struct {
		Home         string `yaml:"home"`
		Workspace    string `yaml:"workspace"`
		ObsidianRoot string `yaml:"obsidian_root"`
	} `yaml:"paths"`
	Platforms struct {
		Opencode struct {
			ConfigDir string `yaml:"config_dir"`
		} `yaml:"opencode"`
	} `yaml:"platforms"`
	StateDir         string   `yaml:"state_dir"`
	AllowIdentifiers []string `yaml:"allow_identifiers"`
}

// DefaultPath returns the default config location under $HOME.
func DefaultPath() string {
	h, err := os.UserHomeDir()
	if err != nil {
		return ".config/aiy-warp/warp.config"
	}
	return filepath.Join(h, ".config", "aiy-warp", "warp.config")
}

// Load reads the config if present; a missing file yields an empty Config.
func Load(path string) (*Config, error) {
	cfg := &Config{}
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}
	if err := yaml.Unmarshal(b, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// OpencodeDest resolves the opencode config directory (default ~/.config/opencode).
func (c *Config) OpencodeDest() string {
	if c.Platforms.Opencode.ConfigDir != "" {
		return expand(c.Platforms.Opencode.ConfigDir)
	}
	h, err := os.UserHomeDir()
	if err != nil {
		return ".config/opencode"
	}
	return filepath.Join(h, ".config", "opencode")
}

// StateDirResolved resolves where warp.lock lives (default ~/.config/aiy-warp).
func (c *Config) StateDirResolved() string {
	if c.StateDir != "" {
		return expand(c.StateDir)
	}
	h, err := os.UserHomeDir()
	if err != nil {
		return ".config/aiy-warp"
	}
	return filepath.Join(h, ".config", "aiy-warp")
}

// AllowIDSet returns the declared public identifiers as a lookup set.
func (c *Config) AllowIDSet() map[string]bool {
	out := map[string]bool{}
	for _, id := range c.AllowIdentifiers {
		if t := strings.TrimSpace(id); t != "" {
			out[t] = true
		}
	}
	return out
}

// expand resolves a leading "~" against $HOME.
func expand(p string) string {
	if p == "~" || strings.HasPrefix(p, "~/") {
		if h, err := os.UserHomeDir(); err == nil {
			return filepath.Join(h, strings.TrimPrefix(p, "~/"))
		}
	}
	return p
}
