package command

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/mirror-media/major-tom-go/config"
	"github.com/mirror-media/major-tom-go/internal/k8sop"
	"github.com/pkg/errors"
)

// List provide infomation about cluster, stages, and services. It should also provide helm message if input is invalid
func List(ctx context.Context, clusterConfigs config.K8S, textParts []string) (messages []string, err error) {
	switch len(textParts) {
	case 0:
		var projects []string
		for key := range clusterConfigs {
			projects = append(projects, string(key))
		}
		sort.Strings(projects)
		messages = []string{
			"The following projects are available: " + strings.Join(projects, ", "),
		}

	case 1:
		// List stages
		project := textParts[0]
		stageConfigs, isExisting := clusterConfigs[config.Project(project)]
		if !isExisting {
			// TODO call help
			return []string{"call help"}, errors.Errorf("project(%s) doesn't exist", project)
		}
		stages := make([]string, 0, len(stageConfigs))
		for k := range stageConfigs {
			stages = append(stages, string(k))
		}
		sort.Strings(stages)
		messages = []string{
			fmt.Sprintf("The following stages are available for %s: %s", project, strings.Join(stages, ", ")),
		}
	case 2:
		kubeConfigPath, err := k8sop.SwitchKubeConfig(clusterConfigs, textParts[0], textParts[1])
		if err != nil {
			return nil, err
		}
		releases, err := k8sop.ListReleases(ctx, kubeConfigPath)
		if err != nil {
			return nil, err
		}

		messages = make([]string, len(releases))
		for i, release := range releases {
			messages[i] = fmt.Sprintf("%s: %s", release.Name, release.Status)
		}
	default:
		err = errors.Errorf("What is going on now?")
	}

	return messages, err
}
