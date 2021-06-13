// Package k8sop is responsible of the implementation involving helm and Kubernetes
package k8sop

import (
	"github.com/mirror-media/major-tom-go/config"
	"github.com/pkg/errors"
)

func SwitchKubeConfig(clusterConfigs config.K8S, project, stage string) (kubeConfigPath string, err error) {

	s, isExisting := clusterConfigs[config.Project(project)]
	if !isExisting {
		return "", errors.Errorf("project(%s) doesn't exist", project)
	}
	config, isExisting := s[config.Stage(stage)]
	if !isExisting {
		return "", errors.Errorf("stage(%s) doesn't exist for project(%s)", stage, project)
	}

	return string(config), nil
}
