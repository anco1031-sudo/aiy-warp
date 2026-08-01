package manifest

import (
	"reflect"

	"github.com/anco1031-sudo/aiy-warp/install/src/internal/frontmatter"
	"gopkg.in/yaml.v3"
)

// MergeAgentContent performs the field-level merge for agent files
// (CLI_SPEC.md §6.2): repo content wins, except the host-local `model` and
// `permission` frontmatter values, which are preserved. If either file has no
// parseable frontmatter, the repo content is returned unchanged.
func MergeAgentContent(repoContent, hostContent string) (string, error) {
	merged, _, err := MergeAgentContentV(repoContent, hostContent)
	return merged, err
}

// MergeAgentContentV is MergeAgentContent plus warnings: when the host file's
// frontmatter cannot be parsed, host-local tuning (model/permission) is
// silently discarded by the repo-wins fallback — the warning makes that visible
// (F12).
func MergeAgentContentV(repoContent, hostContent string) (string, []string, error) {
	var warns []string
	repoFM, repoBody, err := frontmatter.SplitFrontmatter(repoContent)
	if err != nil || repoFM == "" {
		return repoContent, warns, nil
	}
	hostFM, _, err := frontmatter.SplitFrontmatter(hostContent)
	if err != nil || hostFM == "" {
		warns = append(warns, "host file frontmatter unparseable — repo wins, host model/permission tuning discarded")
		return repoContent, warns, nil
	}
	var repoMap, hostMap map[string]any
	if err := yaml.Unmarshal([]byte(repoFM), &repoMap); err != nil {
		return repoContent, warns, nil
	}
	if err := yaml.Unmarshal([]byte(hostFM), &hostMap); err != nil {
		warns = append(warns, "host file frontmatter YAML invalid — repo wins, host model/permission tuning discarded")
		return repoContent, warns, nil
	}

	changed := false
	for _, key := range HostLocalFields {
		hv, ok := hostMap[key]
		if !ok {
			continue
		}
		if rv, exists := repoMap[key]; !exists || !reflect.DeepEqual(rv, hv) {
			repoMap[key] = hv
			changed = true
		}
	}
	if !changed {
		return repoContent, warns, nil
	}

	mergedFM, err := yaml.Marshal(repoMap)
	if err != nil {
		return repoContent, warns, nil
	}
	return "---\n" + string(mergedFM) + "---\n" + repoBody, warns, nil
}
