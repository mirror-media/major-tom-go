package command

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/go-git/go-git/v5/plumbing"
	gootkitconfig "github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"github.com/mirror-media/major-tom-go/v2/config"
	"github.com/mirror-media/major-tom-go/v2/gitop"
	mjcontext "github.com/mirror-media/major-tom-go/v2/internal/context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type yamlCFG struct {
	path      string
	valueType string
}

type Deployment struct {
	ctx      context.Context
	codebase *config.Codebase
	service  *config.Service
	project  string
	stage    string
	imageTag string
	caller   string
}

var deployChannel = make(chan Deployment, 64)

// Deploy certain configuration to a service. textParts in interpreted as [project, stage, service, ...cfg:arg]
func Deploy(ctx context.Context, k8sRepo config.KubernetesConfigsRepo, textParts []string, caller string) (messages []string, err error) {
	if !DeployWorker.isRunning {
		return nil, errors.New("deploy worker is not running")
	}

	if len(textParts) < 1 {
		return nil, errors.New("call help")
	}

	// Compare and retrieve the repo before we engage the deployment, so we can pass the repo to deploy worker for clearer intention
	codebases := k8sRepo.Configs
	repoNameInCMD, textParts := pop(textParts, 0)
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

	// deploy requires env only and it only supports image-tag only

	textParts, stage, err := popValue(textParts, "env", "=")
	if err != nil {
		return nil, errors.Wrap(err, "getting env for deployment encountered an error")
	} else if stage == "prod" {
		return nil, errors.New("deploy command doesn't support prod env")
	}

	textParts, image, err := popValue(textParts, "image-tag", "=")
	if err != nil {
		return nil, errors.Wrap(err, "getting image-tag for deployment encountered an error")
	}

	if len(textParts) != 0 {
		return nil, errors.New(strings.Join(textParts, ", ") + " are not supported")
	}

	timeout := 5 * time.Minute
	ch := make(chan response)
	newCtx := context.WithValue(ctx, mjcontext.ResponseChannel, ch)
	newCtx, cancelFn := context.WithTimeout(newCtx, timeout)
	defer cancelFn()
	deployChannel <- Deployment{
		ctx:      newCtx,
		codebase: codebase,
		stage:    stage,
		imageTag: image,
		caller:   caller,
	}

	select {
	case commandResponse := <-ch:
		return commandResponse.Messages, commandResponse.Error
	case <-newCtx.Done():
		return nil, errors.Errorf("\"%s\" command has timeouted(%f)", strings.Join(textParts, " "), timeout.Minutes())
	}
}

type deployWorker struct {
	once      sync.Once
	isRunning bool
	k8sRepo   *gitop.Repository
}

var DeployWorker deployWorker

func (w *deployWorker) Set(gitConfigs config.GitConfig) {
	var err error
	w.k8sRepo, err = gitop.GetK8SConfigsRepository(gitConfigs)
	if err != nil {
		logrus.Fatal(err)
	}

	go w.once.Do(func() {
		w.isRunning = true
		logrus.Info("the deploy worker is running now....")
		for {
			deployment := <-deployChannel
			deploy(deployment.ctx, w.k8sRepo, *deployment.codebase, deployment.stage, deployment.project, deployment.imageTag, deployment.caller)
		}
	})
}

func hardReset(repository *gitop.Repository, commit plumbing.Hash) (hardResetFN func() error) {
	return func() error { return repository.HardResetToCommit(commit) }
}

// deploy certain configuration to a service. textParts in interpreted as [repo, env=value..., ...cfg=value]
func deploy(ctx context.Context, k8sRepo *gitop.Repository, codebase config.Codebase, stage, project, imageTag, caller string) {
	var messages = []string{}
	var err error
	ch := ctx.Value(mjcontext.ResponseChannel).(chan response)

	repo := k8sRepo
	hash, errHash := repo.GetHeadHash()
	if errHash != nil {
		err = errors.Wrap(errHash, fmt.Sprintf("getting head hash of repo(%s) has error", "kubernetes-configs"))
		logrus.Error(err)
	}
	if err != nil {
		ch <- response{
			Messages: messages,
			Error:    err,
		}
		return
	}
	hardResetFn := hardReset(repo, hash)

	path, err := codebase.GetImageKustomizationPath(stage, project)
	logrus.Debug(path)
	if err != nil {
		ch <- response{
			Messages: messages,
			Error:    err,
		}
		return
	}
	kustomization, err := repo.GetFile(path)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("cannot get file(%s)", path))
		ch <- response{
			Messages: messages,
			Error:    err,
		}
		return
	}

	f := kustomization

	// operation starts here. worktree needs to be cleaned if disaster happens
	b, err := io.ReadAll(f)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("reading %s has error", path))
		logrus.Warn(err)
		ch <- response{
			Messages: messages,
			Error:    err,
		}
		_ = hardResetFn()
		return
	}
	valueConfig := gootkitconfig.New(path)
	valueConfig.AddDriver(yaml.Driver)
	err = valueConfig.LoadStrings(gootkitconfig.Yaml, string(b))
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("loading YAML from %s has error", path))
		logrus.Warn(err)
		ch <- response{
			Messages: messages,
			Error:    err,
		}
		_ = hardResetFn()
		return
	}

	images0 := valueConfig.Get(".images.0", true).(map[interface{}]interface{})
	images0["newTag"] = imageTag
	err = valueConfig.Set(".images.0", images0, true)
	if err != nil {
		ch <- response{
			Messages: messages,
			Error:    errors.Wrap(err, fmt.Sprintf("fail to set newTag in %s", path)),
		}
		_ = hardResetFn()
		return
	}

	var pendingProject string

	if project != "" {
		pendingProject = "/" + project
	}

	messages = append(messages, fmt.Sprintf("deploy(%s/%s%s): deployed by %s", codebase.Repo, stage, pendingProject, caller), "", fmt.Sprintf("Set %s(%s) to %v", "image-tag", "images.0.newTag", imageTag))

	err = f.Truncate(0)
	f.Close()
	if err != nil {
		ch <- response{
			Messages: messages,
			Error:    errors.Wrap(err, fmt.Sprintf("clean file before writing to %s has error", path)),
		}
		_ = hardResetFn()
		return
	}

	f, _ = repo.GetFile(path)
	valueConfig.DumpTo(f, gootkitconfig.Yaml)
	if err != nil {
		ch <- response{
			Messages: messages,
			Error:    errors.Wrap(err, fmt.Sprintf("writing to %s has error", path)),
		}
		_ = hardResetFn()
		return
	}

	// command operation finished
	// now git operations starts

	repo.AddFile(path)
	if err != nil {
		ch <- response{
			Messages: messages,
			Error:    errors.Wrap(err, fmt.Sprintf("adding %s to staging area has error", path)),
		}
		_ = hardResetFn()
		return
	}

	err = repo.Commit(path, caller, strings.Join(messages, "\n"))
	if err != nil {
		ch <- response{
			Messages: append([]string{"this operation failed"}, messages...),
			Error:    errors.Wrap(err, fmt.Sprintf("commits for project(%s) has error", project)),
		}
		_ = hardResetFn()
		return
	}
	err = repo.Push()
	if err != nil {
		ch <- response{
			Messages: append([]string{"this operation failed"}, messages...),
			Error:    errors.Wrap(err, fmt.Sprintf("push commits for project(%s) has error", project)),
		}
		_ = hardResetFn()
		return
	}

	ch <- response{
		Messages: messages,
		Error:    err,
	}
}
