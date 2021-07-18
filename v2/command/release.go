package command

import (
	"context"
	"strings"
	"time"

	"github.com/mirror-media/major-tom-go/v2/config"
	mjcontext "github.com/mirror-media/major-tom-go/v2/internal/context"
	"github.com/pkg/errors"
)

// Release a new image tag to a repo in a project. texts is interpreted as [project=value, imag-tag=value]
func Release(ctx context.Context, k8sRepo config.KubernetesConfigsRepo, texts []string, message, caller string) (messages []string, err error) {
	if !DeployWorker.isRunning {
		return nil, errors.New("deploy worker is not running")
	}

	if len(texts) < 1 {
		return nil, errors.New("call help")
	}

	// Compare and retrieve the repo before we engage the deployment, so we can pass the repo to deploy worker for clearer intention
	codebases := k8sRepo.Configs
	repoNameInCMD, texts := pop(texts, 0)
	var codebase *config.Codebase
	for _, c := range codebases {
		if repoNameInCMD == c.Repo {
			codebase = &c
			break
		}
	}

	if codebase == nil {
		return nil, errors.New("invalid repo name")
	}

	// deploy requires project only and it retquires image-tag and project

	texts, project, err := popValue(texts, "project", "=")
	if err != nil {
		return nil, errors.Wrap(err, "getting project for deployment encountered an error")
	}

	texts, image, err := popValue(texts, "image-tag", "=")
	if err != nil {
		return nil, errors.Wrap(err, "getting image-tag for deployment encountered an error")
	}

	if len(texts) != 0 {
		return nil, errors.New("Major Tom does not support: " + strings.Join(texts, ", "))
	}

	timeout := 5 * time.Minute
	ch := make(chan response)
	newCtx := context.WithValue(ctx, mjcontext.ResponseChannel, ch)
	newCtx, cancelFn := context.WithTimeout(newCtx, timeout)
	defer cancelFn()
	deployChannel <- Deployment{
		ctx:      newCtx,
		codebase: codebase,
		stage:    "prod",
		project:  project,
		imageTag: image,
		caller:   caller,
		message:  message,
	}

	select {
	case commandResponse := <-ch:
		return commandResponse.Messages, commandResponse.Error
	case <-newCtx.Done():
		return nil, errors.Errorf("\"%s\" command has timeouted(%f)", strings.Join(texts, " "), timeout.Minutes())
	}
}
