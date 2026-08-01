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
	repoFM, repoBody, err := frontmatter.SplitFrontmatter(repoContent)
	if err != nil || repoFM == "" {
		return repoContent, nil
	}
	hostFM, _, err := frontmatter.SplitFrontmatter(hostContent)
	if err != nil || hostFM == "" {
		return repoContent, nil
	}
	var repoMap, hostMap map[string]any
	if err := yaml.Unmarshal([]byte(repoFM), &repoMap); err != nil {
		return repoContent, nil
	}
	if err := yaml.Unmarshal([]byte(hostFM), &hostMap); err != nil {
		return repoContent, nil
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
		return repoContent, nil
	}

	mergedFM, err := yaml.Marshal(repoMap)
	if err != nil {
		return repoContent, nil
	}
	return "---\n" + string(mergedFM) + "---\n" + repoBody, nil
}
