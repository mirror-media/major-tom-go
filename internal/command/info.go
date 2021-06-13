package command

import (
	"context"
	"fmt"

	"github.com/mirror-media/major-tom-go/config"
	"github.com/mirror-media/major-tom-go/internal/k8sop"
	"github.com/pkg/errors"
)

// Info provides infomation of service. textParts in interpreted as [project, stage, service]
func Info(ctx context.Context, clusterConfigs config.K8S, textParts []string) (messages []string, err error) {
	switch len(textParts) {
	// Info does not provice information for projects and stages
	case 0, 1, 2:
		messages = []string{"call help"}
		err = errors.New("Please provide more details...")
	// service
	case 3:
		kubeConfigPath, errConfig := k8sop.SwitchKubeConfig(clusterConfigs, textParts[0], textParts[1])
		if errConfig != nil {
			messages = []string{"call help"}
			err = errConfig
			break
		}

		info, errGet := k8sop.GetDeploymentInfo(ctx, kubeConfigPath, textParts[2])
		if errGet != nil {
			err = errGet
			messages = []string{"call list"}
			break
		}
		messages = []string{fmt.Sprintf("%s\n\tImageTag: %s\n\tAvailable pods: %d\n\tReady pods: %d\n\tUpdated pods: %d", textParts[2], info.ImageTag, info.Available, info.Ready, info.Updated)}
	}

	return messages, err
}
