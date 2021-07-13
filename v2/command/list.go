package command

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/mirror-media/major-tom-go/v2/config"
	"github.com/mirror-media/major-tom-go/v2/k8sop"
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
	case 3:
		// kubeConfigPath, err := k8sop.SwitchKubeConfig(clusterConfigs, textParts[0], textParts[1])
		// if err != nil {
		// 	return nil, err
		// }
		// // var namespace string
		// if textParts[2] == "cronjobs" {
		// 	namespace = "cron"
		// } else {
		// 	namespace = "default"
		// }
		// TODO
		messages = append(messages, "history is not implemented yet")
		// versions, err := k8sop.ListHelmReleaseVersion(kubeConfigPath, textParts[2], namespace, 11)
		// if err != nil {
		// 	return nil, err
		// }
		// for _, version := range versions {
		// 	messages = append(messages, fmt.Sprintf("%s(%d)", textParts[2], version.Version))
		// 	messages = append(messages, fmt.Sprintf("\t%s: %v", "ImageTag", version.ImageTag))
		// 	messages = append(messages, fmt.Sprintf("\t%s: %v", "Pods", version.ReplicaCount))
		// 	messages = append(messages, fmt.Sprintf("\t%s", "AutoScaling"))
		// 	messages = append(messages, fmt.Sprintf("\t\t%s: %v", "Enabled", version.AutoScaling.Enabled))
		// 	messages = append(messages, fmt.Sprintf("\t\t%s: %v", "Max Pods", version.AutoScaling.MaxReplicas))
		// 	messages = append(messages, fmt.Sprintf("\t\t%s: %v", "Min Pods", version.AutoScaling.MinReplicas))
		// }
	default:
		err = errors.Errorf("What is going on now?")
	}

	return messages, err
}
