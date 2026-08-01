// Package manifest implements the warp.lock sha256 manifest, the 3-way
// install diff, and the host-local frontmatter merge policy (CLI_SPEC.md §6).
package manifest

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Lock is the installed-state manifest written beside warp.config.
type Lock struct {
	Platform  string            `yaml:"platform"`
	SourceRev string            `yaml:"source_rev"`
	Timestamp string            `yaml:"timestamp"`
	Files     map[string]string `yaml:"files"` // repo-relative path → drift hash
}

// Load reads a lock file; a missing file yields (nil, nil).
func Load(path string) (*Lock, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var l Lock
	if err := yaml.Unmarshal(b, &l); err != nil {
		return nil, err
	}
	return &l, nil
}

// Save writes the lock atomically (tmp + rename). The tmp file is fsynced
// before rename so a crash can never leave an empty/partial lock (F11).
func Save(path string, l *Lock) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := yaml.Marshal(l)
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	f, err := os.Open(tmp)
	if err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// New builds a lock for an install.
func New(platform, rev string, files map[string]string) *Lock {
	return &Lock{
		Platform:  platform,
		SourceRev: rev,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Files:     files,
	}
}

// Hashes computes the drift hash of each file content (relpath → hash).
func Hashes(contents map[string]string) map[string]string {
	out := make(map[string]string, len(contents))
	for p, c := range contents {
		out[p] = KitHash(c, p)
	}
	return out
}

// Sha256Hex returns the hex sha256 of content.
func Sha256Hex(content string) string {
	h := sha256.Sum256([]byte(content))
	return hex.EncodeToString(h[:])
}
