package command

import (
	"context"
	"fmt"
	"strings"

	"github.com/mirror-media/major-tom-go/v2/config"
	"github.com/pkg/errors"
)

// List provide infomation about cluster, stages, and services. It should also provide helm message if input is invalid
func List(ctx context.Context, clusterConfigs config.K8S, textParts []string) (message []string, err error) {
	switch len(textParts) {
	case 0:
		var projects []string
		for key := range clusterConfigs {
			projects = append(projects, string(key))
		}
		message = []string{
			"The following projects are available: " + strings.Join(projects, ", "),
		}

	case 1:
		// List stages
		project := textParts[0]
		stages, isExisting := clusters[project]
		if !isExisting {
			// TODO call help
			return []string{"call help"}, errors.Errorf("project(%s) doesn't exist", project)
		}
		stages := make([]string, 0, len(stageConfigs))
		for k := range stageConfigs {
			stages = append(stages, string(k))
		}
		message = []string{
			fmt.Sprintf("The following stages are available for %s: %s", project, strings.Join(stages, ", ")),
		}
	}
	return message, err
}
